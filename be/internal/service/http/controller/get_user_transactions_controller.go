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

	User struct {
		Id    uint64 `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	Transaction struct {
		Id     uint64 `json:"id"`
		From   User   `json:"from"`
		To     User   `json:"to"`
		Amount int64  `json:"amount"`
		Status uint64 `json:"status"`
		Ctime  uint64 `json:"ctime"`
	}

	GetUserTransactionsData struct {
		Transactions   []*Transaction `json:"transactions"`
		NextPagination *Pagination    `json:"next_pagination"`
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

	var transactions []*Transaction
	if len(getResp.GetTransactions()) > 0 {
		var ids []uint64
		for _, t := range getResp.Transactions {
			ids = append(ids, t.FromUser, t.ToUser)
			transactions = append(transactions, &Transaction{
				Id:     t.Id,
				Amount: t.Amount,
				Status: t.Status,
				Ctime:  t.Ctime,
			})
		}
		batchResp, err := c.uClient.GetBatch(reqCtx, &proto.GetBatchReq{Ids: ids})
		if err != nil {
			ctx.JSON(http.StatusOK, &TransferResp{Code: common.GetCodeFromErr(err)})
			return
		}

		for i, user := range batchResp.Users {
			if i%2 == 0 {
				transactions[i/2].From = User{
					Id:    user.Id,
					Name:  user.Name,
					Email: user.Email,
				}
			} else {
				transactions[i/2].To = User{
					Id:    user.Id,
					Name:  user.Name,
					Email: user.Email,
				}
			}
		}

	}

	resp := &GetUserTransactionsResp{
		Code: common.GetCodeFromErr(err),
		Data: &GetUserTransactionsData{
			Transactions:   transactions,
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
