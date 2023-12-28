package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"

	"github.com/chientranthien/goldenpay/internal/common"
	"github.com/chientranthien/goldenpay/internal/proto"
)

type (
	TransferReq struct {
		ToEmail string `json:"to_email"`
		Amount  int64  `json:"amount"`
	}

	TransferResp struct {
		Code        *common.Code        `json:"code"`
		Transaction TransferTransaction `json:"transaction"`
	}

	TransferTransaction struct {
		id uint64 `json:"id"`
	}

	TransferController struct {
		uClient proto.UserServiceClient
		wClient proto.WalletServiceClient
	}
)

func NewTransferController(
	uClient proto.UserServiceClient,
	wClient proto.WalletServiceClient,
) *TransferController {
	return &TransferController{uClient: uClient, wClient: wClient}
}

func (c TransferController) Do(ctx *gin.Context) {
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

	reqCtx := common.Ctx()
	authzResp, err := c.uClient.Authz(reqCtx, &proto.AuthzReq{Token: token})
	if err != nil {
		ctx.JSON(http.StatusOK, &TransferResp{Code: &common.Code{Id: int32(codes.InvalidArgument), Msg: "can't send to yourself"}})
		return
	}

	if authzResp.Metadata.GetEmail() == req.ToEmail {
		ctx.JSON(http.StatusOK, &TransferResp{Code: common.GetCodeFromErr(err)})
		return

	}

	getResp, err := c.uClient.GetByEmail(reqCtx, &proto.GetByEmailReq{Email: req.ToEmail})
	if err != nil {
		ctx.JSON(http.StatusOK, &TransferResp{Code: common.GetCodeFromErr(err)})
		return
	}
	toUser := getResp.User
	if toUser.Status != common.UserStatusActive {
		ctx.JSON(
			http.StatusOK,
			&TransferResp{Code: common.NewCode(int32(codes.InvalidArgument), "user's status not active")},
		)
		return
	}

	transferResp, err := c.wClient.Transfer(reqCtx, &proto.TransferReq{
		FromUser: authzResp.Metadata.UserId,
		ToUser:   toUser.Id,
		Amount:   req.Amount,
	})
	resp := &TransferResp{
		Code: common.GetCodeFromErr(err),
	}
	if transferResp != nil {
		resp.Transaction = TransferTransaction{id: transferResp.TransactionId}
	}
	ctx.JSON(http.StatusOK, resp)

	if err != nil {
		common.L().Errorw("transferErr", "req", req, "err", err)
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
