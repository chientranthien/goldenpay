package controller

import (
	"context"

	"github.com/chientranthien/goldenpay/internal/common"
	httpcommon "github.com/chientranthien/goldenpay/internal/common/http"
	"github.com/chientranthien/goldenpay/internal/proto"
)

type (
	TopupBody struct {
		Amount int64 `json:"amount"`
	}

	TopupData struct {
		id uint64 `json:"id"`
	}

	TopupController struct {
		wClient proto.WalletServiceClient

		ctx  context.Context
		req  httpcommon.Req
		body *TopupBody
	}
)

func NewTopupController(
	wClient proto.WalletServiceClient,
) *TopupController {
	return &TopupController{wClient: wClient}
}

func (c TopupController) Do() (common.AnyPtr, common.Code) {
	topupResp, err := c.wClient.Topup(c.ctx, &proto.TopupReq{
		UserId: c.req.Metadata.UserId,
		Amount: c.body.Amount,
	})
	code := common.GetCodeFromErr(err)
	if !code.Success() {
		common.L().Errorw("TopupErr", "body", c.body, "err", err)
		return nil, code
	}

	return &TopupData{id: topupResp.TopupId}, common.CodeSuccess
}

func (c *TopupController) Take(ctx context.Context, req httpcommon.Req) common.Code {
	if req.Metadata.UserId <= 0 {
		return common.CodeInvalidMetadata
	}

	if b, ok := req.Body.(*TopupBody); ok {
		c.body = b
		c.ctx = ctx
		c.req = req
	} else {
		return common.CodeBody
	}

	if c.body.Amount <= 0 {
		return common.CodeInvalidArgument
	}

	return common.CodeSuccess
}
