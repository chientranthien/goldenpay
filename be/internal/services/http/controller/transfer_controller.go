package controller

import (
	"log"
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
	Code        *common.Code `json:"code"`
	Transaction Transaction  `json:"transaction"`
}
type Transaction struct {
	id uint64 `json:"id"`
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

	authzResp, err := c.uClient.Authz(common.Ctx(), &proto.AuthzReq{Token: token})
	if err != nil {
		ctx.JSON(http.StatusOK, &TransferResp{Code: common.GetCodeFromErr(err)})
		return
	}

	transferResp, err := c.wClient.Transfer(common.Ctx(), &proto.TransferReq{
		FromUser: authzResp.Metadata.UserId,
		ToEmail:  req.ToEmail,
		Amount:   req.Amount,
	})
	ctx.JSON(
		http.StatusOK,
		&TransferResp{
			Code:        common.GetCodeFromErr(err),
			Transaction: Transaction{id: transferResp.TransactionId},
		},
	)

	if err != nil {
		log.Printf("failed to transfer, err=%v", err)
		return
	}
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
