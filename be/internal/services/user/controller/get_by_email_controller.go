package controller

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/chientranthien/goldenpay/internal/common"
	"github.com/chientranthien/goldenpay/internal/proto"
	"github.com/chientranthien/goldenpay/internal/services/user/biz"
)

type GetByEmailController struct {
	biz *biz.UserBiz
}

func NewGetByEmailController(biz *biz.UserBiz) *GetByEmailController {
	return &GetByEmailController{biz: biz}
}

func (c GetByEmailController) GetByEmail(req *proto.GetByEmailReq) (*proto.GetByEmailResp, error) {
	if err := c.validate(req); err != nil {
		return nil, err
	}

	return c.biz.GetByEmail(req)
}

func (c GetByEmailController) validate(req *proto.GetByEmailReq) error {
	if common.ValidateEmail(req.Email) != nil {
		return status.New(codes.InvalidArgument, "invalid email").Err()
	}

	return nil
}
