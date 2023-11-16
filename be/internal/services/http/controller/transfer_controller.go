package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/chientranthien/goldenpay/internal/common"
	"github.com/chientranthien/goldenpay/internal/proto"
)

type TransferReq struct {
	ToEmail string `json:"to_email,omitempty"`
	Amount  int64  `json:"amount,omitempty"`
}

type TransferResp struct {
	Code *common.Code `json:"code"`
}
type TransferController struct {
	uClient proto.UserServiceClient
	wClient proto.WalletServiceClient
}

func NewTransferController(wClient proto.WalletServiceClient) *TransferController {
	return &TransferController{wClient: wClient}
}

func (c TransferController) Transfer(ctx gin.Context) {
	req := &TransferReq{}
	if ctx.BindJSON(req) != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	if code := c.validate(req); !code.IsSuccess() {
		resp := &TransferResp{
			Code: code,
		}
		ctx.JSON(http.StatusOK, resp)
		return
	}

	token, _ := ctx.Cookie(TokenCookie)
	if token == "" {
		ctx.JSON(http.StatusOK, &TransferResp{Code: common.CodeUnauthenticated})
		return
	}

	_, err := c.uClient.Authz(common.Ctx(), &proto.AuthzReq{Token: token})
	if err != nil {
		ctx.JSON(http.StatusOK, &TransferResp{Code: common.GetCodeFromErr(err)})
		return
	}

	c.wClient.Transfer(common.Ctx(),&proto.TransferReq{
		ToEmail:   0,
		Amount:   req.Amount,
	})

}

func (c TransferController) validate(req *TransferReq) *common.Code {
	if common.ValidateEmail(req.ToEmail) != nil {
		return common.CodeInvalidArgument
	}

	if req.Amount <= 0 {
		return common.CodeInvalidArgument
	}
	return common.CodeSuccess
}
