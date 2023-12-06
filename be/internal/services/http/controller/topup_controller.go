package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/chientranthien/goldenpay/internal/common"
	"github.com/chientranthien/goldenpay/internal/proto"
)

type (
	TopupReq struct {
		Amount int64 `json:"amount"`
	}

	TopupResp struct {
		Code  *common.Code `json:"code"`
		Topup Topup        `json:"topup"`
	}

	Topup struct {
		id uint64 `json:"id"`
	}

	TopupController struct {
		uClient proto.UserServiceClient
		wClient proto.WalletServiceClient
	}
)

func NewTopupController(
	uClient proto.UserServiceClient,
	wClient proto.WalletServiceClient,
) *TopupController {
	return &TopupController{uClient: uClient, wClient: wClient}
}

func (c TopupController) Do(ctx *gin.Context) {
	req := &TopupReq{}
	if ctx.BindJSON(req) != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	if code := c.validate(req); !code.IsSuccess() {
		resp := &TopupResp{
			Code: code,
		}
		ctx.JSON(http.StatusOK, resp)
		return
	}

	token, _ := ctx.Cookie(TokenCookie)
	if token == "" {
		ctx.JSON(http.StatusOK, &TopupResp{Code: common.CodeUnauthenticated})
		return
	}

	reqCtx := common.Ctx()
	authzResp, err := c.uClient.Authz(reqCtx, &proto.AuthzReq{Token: token})
	if err != nil {
		ctx.JSON(http.StatusOK, &TopupResp{Code: common.GetCodeFromErr(err)})
		return
	}

	topupResp, err := c.wClient.Topup(reqCtx, &proto.TopupReq{
		UserId: authzResp.Metadata.UserId,
		Amount:   req.Amount,
	})
	resp := &TopupResp{
		Code: common.GetCodeFromErr(err),
	}
	if topupResp != nil {
		resp.Topup = Topup{id: topupResp.TopupId}
	}
	ctx.JSON(http.StatusOK, resp)

	if err != nil {
		common.L().Errorw("TopupErr", "req", req, "err", err)
		return
	}
}

func (c TopupController) validate(req *TopupReq) *common.Code {
	if req.Amount <= 0 {
		return common.CodeInvalidArgument
	}
	return common.CodeSuccess
}
