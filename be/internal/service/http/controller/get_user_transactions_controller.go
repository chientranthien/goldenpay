package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/chientranthien/goldenpay/internal/common"
	"github.com/chientranthien/goldenpay/internal/proto"
)

type (
	Pagination struct {
		Val     int64  `json:"val"`
		Limit   uint32 `json:"limit"`
		HasMore bool   `json:"has_more"`
	}

	GetUserTransactionsReq struct {
		Pagination *Pagination `json:"pagination"`
	}

	GetUserTransactionsData struct {
		Transactions   []*proto.Transaction `json:"transactions"`
		NextPagination *Pagination          `json:"next_pagination"`
	}

	GetUserTransactionsResp struct {
		Code *common.Code             `json:"code"`
		Data *GetUserTransactionsData `json:"data"`
	}

	GetUserTransactionsController struct {
		uClient proto.UserServiceClient
		wClient proto.WalletServiceClient
	}
)

func NewGetUserTransactionsController(uClient proto.UserServiceClient, wClient proto.WalletServiceClient) *GetUserTransactionsController {
	return &GetUserTransactionsController{uClient: uClient, wClient: wClient}
}

func NewPagination(val int64, limit uint32, hasMore bool) *Pagination {
	return &Pagination{Val: val, Limit: limit, HasMore: hasMore}
}
func FromPagination(p *proto.Pagination) *Pagination {
	if p == nil {
		return nil
	}

	return &Pagination{Val: p.Val, Limit: p.Limit, HasMore: p.HasMore}
}

func (c GetUserTransactionsController) Do(ctx *gin.Context) {
	req := &GetUserTransactionsReq{}
	if ctx.BindJSON(req) != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	token, _ := ctx.Cookie(TokenCookie)
	if token == "" {
		ctx.JSON(http.StatusOK, &TransferResp{Code: common.CodeUnauthenticated})
		return
	}

	reqCtx := common.Ctx()
	authzResp, err := c.uClient.Authz(reqCtx, &proto.AuthzReq{Token: token})
	if err != nil {
		ctx.JSON(http.StatusOK, &TransferResp{Code: common.GetCodeFromErr(err)})
		return
	}

	getResp, err := c.wClient.GetUserTransactions(reqCtx, &proto.GetUserTransactionsReq{
		Cond: &proto.GetUserTransactionsCond{
			User: &proto.GetUserTransactionsCond_UserCond{Eq: authzResp.Metadata.UserId},
		},
		Pagination: req.Pagination.ToServicePagination(),
	})

	resp := &GetUserTransactionsResp{
		Code: common.GetCodeFromErr(err),
		Data: &GetUserTransactionsData{
			Transactions:   getResp.Transactions,
			NextPagination: FromPagination(getResp.NextPagination),
		},
	}

	ctx.JSON(http.StatusOK, resp)
	if err != nil {
		common.L().Errorw("getUserTransactionsErr", "req", req, "err", err)
		return
	}
}

func (c GetUserTransactionsController) validate(req *GetUserTransactionsReq) *common.Code {
	return common.CodeSuccess
}

func (p Pagination) ToServicePagination() *proto.Pagination {
	return &proto.Pagination{
		Val:     p.Val,
		Limit:   p.Limit,
		HasMore: p.HasMore,
	}
}
