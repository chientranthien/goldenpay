package controller

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/chientranthien/goldenpay/internal/proto"
	"github.com/chientranthien/goldenpay/internal/services/wallet/biz"
)

type GetController struct {
	biz *biz.WalletBiz
}

func NewGetController(biz *biz.WalletBiz) *GetController {
	return &GetController{biz: biz}
}

func (c GetController) Do(req *proto.GetWalletReq) (*proto.GetWalletResp, error) {
	if err := c.validate(req); err != nil {
		return nil, err
	}

	return c.biz.Get(req)
}

func (c GetController) validate(req *proto.GetWalletReq) error {
	if req.UserId <= 0 {
		return status.Errorf(codes.InvalidArgument, "invalid user_id")
	}

	return nil
}
