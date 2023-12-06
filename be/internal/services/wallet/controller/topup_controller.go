package controller

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/chientranthien/goldenpay/internal/proto"
	"github.com/chientranthien/goldenpay/internal/services/wallet/biz"
)

type TopupController struct {
	biz *biz.WalletBiz
}

func NewTopupController(biz *biz.WalletBiz) *TopupController {
	return &TopupController{biz: biz}
}


func (c TopupController) Do(req *proto.TopupReq) (*proto.TopupResp, error) {
	if err := c.validate(req); err != nil {
		return nil, err
	}


	return c.biz.Topup(req)
}


func (c TopupController) validate(req *proto.TopupReq) error {
	if req.UserId <= 0 {
		return status.New(codes.InvalidArgument, "invalid user_id").Err()
	}

	if req.Amount <= 0 {
		return status.New(codes.InvalidArgument, "invalid amount").Err()
	}

	return nil
}
