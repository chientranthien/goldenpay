package controller

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"github.com/chientranthien/goldenpay/internal/proto"
	"github.com/chientranthien/goldenpay/internal/service/user/biz"
)

type SignupController struct {
	biz *biz.UserBiz
}

func NewSignupController(biz *biz.UserBiz) *SignupController {
	return &SignupController{biz: biz}
}

func (c *SignupController) Signup(_ context.Context, req *proto.SignupReq) (*proto.SignupResp, error) {
	if err := c.validate(req); err != nil {
		return nil, err
	}
	resp, err := c.biz.Signup(req)
	if err != nil {
		return nil, status.New(codes.Internal, "failed to signup").Err()
	}

	return resp, nil
}

func (c SignupController) validate(req *proto.SignupReq) error {
	if len(req.Email) == 0 || len(req.Password) == 0 {
		return status.New(codes.InvalidArgument, "invalid email or password").Err()
	}

	_, err := c.biz.GetByEmail(&proto.GetByEmailReq{Email: req.Email})
	if err != nil && err != gorm.ErrRecordNotFound {
		return status.New(codes.Internal, "failed to get user").Err()
	}

	if err == nil {
		return status.New(codes.InvalidArgument, "existed user").Err()
	}

	return nil
}
