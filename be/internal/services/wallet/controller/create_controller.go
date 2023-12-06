package controller

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"github.com/chientranthien/goldenpay/internal/proto"
	"github.com/chientranthien/goldenpay/internal/services/wallet/biz"
)

type CreateController struct {
	biz *biz.WalletBiz
}

func NewCreateController(biz *biz.WalletBiz) *CreateController {
	return &CreateController{biz: biz}
}

func (c CreateController) Create(req *proto.CreateWalletReq) (*proto.CreateWalletResp, error) {
	if err := c.validate(req); err != nil {
		return nil, err
	}

	return c.biz.Create(req)
}

func (c CreateController) validate(req *proto.CreateWalletReq) error {
	if req.UserId <= 0 {
		return status.Errorf(codes.InvalidArgument, "invalid user_id")
	}

	_, err := c.biz.GetByUserID(req.UserId)
	if err != nil && err != gorm.ErrRecordNotFound {
		return status.Errorf(codes.Internal, "failed to call DB")
	}

	if err == nil {
		return status.Errorf(codes.AlreadyExists, "already exists")
	}

	return nil
}
