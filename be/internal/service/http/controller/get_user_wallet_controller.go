package controller

import (
	"context"

	"github.com/chientranthien/goldenpay/internal/common"
	httpcommon "github.com/chientranthien/goldenpay/internal/common/http"
	"github.com/chientranthien/goldenpay/internal/proto"
)

type (
	GetUserWalletController struct {
		wClient proto.WalletServiceClient
		ctx     context.Context
		req     httpcommon.Req
	}

	GetUserWalletBody struct {
	}

	GetUserWalletData struct {
		Balance int64 `json:"balance"`
	}
)

func NewGetUserWalletController(
	wClient proto.WalletServiceClient,
) *GetUserWalletController {
	return &GetUserWalletController{wClient: wClient}
}

func (c *GetUserWalletController) Do() (common.AnyPtr, common.Code) {
	getWalletResp, err := c.wClient.Get(c.ctx, &proto.GetWalletReq{UserId: c.req.Metadata.UserId})

	code := common.GetCodeFromErr(err)
	if !code.Success() {
		common.L().Errorw("getUserWalletErr", "userId", c.req.Metadata.UserId, "err", err)
		return nil, code
	}

	return &GetUserWalletData{Balance: getWalletResp.Balance}, common.CodeSuccess
}

func (c *GetUserWalletController) Take(ctx context.Context, req httpcommon.Req) common.Code {
	if req.Metadata.UserId <= 0 {
		return common.CodeInvalidMetadata
	}

	c.ctx = ctx
	c.req = req

	return common.CodeSuccess
}
