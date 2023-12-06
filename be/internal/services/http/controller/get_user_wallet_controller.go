package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/chientranthien/goldenpay/internal/common"
	"github.com/chientranthien/goldenpay/internal/proto"
)

type (
	GetUserWalletController struct {
		uClient proto.UserServiceClient
		wClient proto.WalletServiceClient
	}

	GetUserWalletReq struct {
	}

	GetUserWalletResp struct {
		Code *common.Code       `json:"code"`
		Data *GetUserWalletData `json:"data"`
	}

	GetUserWalletData struct {
		Balance int64 `json:"balance"`
	}
)

func NewGetUserWalletController(
	uClient proto.UserServiceClient,
	wClient proto.WalletServiceClient,
) *GetUserWalletController {
	return &GetUserWalletController{uClient: uClient, wClient: wClient}
}

func (c GetUserWalletController) Do(ctx *gin.Context) {
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

	getWalletResp, err := c.wClient.Get(reqCtx, &proto.GetWalletReq{UserId: authzResp.Metadata.UserId})

	resp := &GetUserWalletResp{
		Code: common.GetCodeFromErr(err),
	}
	if getWalletResp != nil {
		resp.Data = &GetUserWalletData{Balance: getWalletResp.Balance}
	}
	ctx.JSON(http.StatusOK, resp)

	if err != nil {
		common.L().Errorw("getUserWalletErr", "userId", authzResp.Metadata.UserId, "err", err)
		return
	}
}
