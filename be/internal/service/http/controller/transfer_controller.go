package controller

import (
	"context"

	"google.golang.org/grpc/codes"

	"github.com/chientranthien/goldenpay/internal/common"
	httpcommon "github.com/chientranthien/goldenpay/internal/common/http"
	"github.com/chientranthien/goldenpay/internal/proto"
)

type (
	TransferBody struct {
		ToEmail string `json:"to_email"`
		Amount  int64  `json:"amount"`
	}

	TransferData struct {
		id uint64 `json:"id"`
	}

	TransferController struct {
		uClient proto.UserServiceClient
		wClient proto.WalletServiceClient

		ctx    context.Context
		req    httpcommon.Req
		body   *TransferBody
		toUser *proto.User
	}
)

func NewTransferController(
	uClient proto.UserServiceClient,
	wClient proto.WalletServiceClient,
) *TransferController {
	return &TransferController{uClient: uClient, wClient: wClient}
}

func (c TransferController) Do() (common.AnyPtr, common.Code) {
	transferResp, err := c.wClient.Transfer(c.ctx, &proto.TransferReq{
		FromUser: c.req.Metadata.UserId,
		ToUser:   c.toUser.GetId(),
		Amount:   c.body.Amount,
	})
	code := common.GetCodeFromErr(err)
	if !code.Success() {
		common.L().Errorw("transferErr", "body", c.body, "err", err)
		return nil, code
	}

	return TransferData{id: transferResp.GetTransactionId()}, common.CodeSuccess
}

func (c *TransferController) Take(ctx context.Context, req httpcommon.Req) common.Code {
	if req.Metadata.UserId <= 0 {
		return common.CodeInvalidMetadata
	}

	if b, ok := req.Body.(*TransferBody); ok {
		c.body = b
		c.ctx = ctx
		c.req = req
	} else {
		return common.CodeBody
	}

	if c.req.Metadata.Email == c.body.ToEmail {
		return common.CodeInvalidArgument

	}

	if common.ValidateEmail(c.body.ToEmail) != nil {
		return common.CodeInvalidArgument
	}

	if c.body.Amount <= 0 {
		return common.CodeInvalidArgument
	}

	getResp, err := c.uClient.GetByEmail(c.ctx, &proto.GetByEmailReq{Email: c.body.ToEmail})
	code := common.GetCodeFromErr(err)
	if !code.Success() {
		return code
	}

	toUser := getResp.GetUser()
	if toUser.GetStatus() != common.UserStatusActive {
		return common.NewCode(int32(codes.InvalidArgument), "user's status not active")
	}
	c.toUser = toUser

	return common.CodeSuccess
}
