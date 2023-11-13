package controller

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/chientranthien/goldenpay/internal/proto"
	"github.com/chientranthien/goldenpay/internal/services/user/biz"
)

type SignupController struct {
	biz *biz.UserBiz
}

func NewSignupController(biz *biz.UserBiz) *SignupController {
	return &SignupController{biz: biz}
}

func (c *SignupController) Signup(ctx context.Context, req *proto.SignupReq) (*proto.SignupResp, error) {
	if err := c.validate(ctx, req); err != nil {
		return nil, err
	}
	resp, err := c.biz.Signup(req)
	if err != nil {
		return nil, status.New(codes.Internal, "failed to signup").Err()
	}

	return resp, nil
}

func (c SignupController) validate(ctx context.Context, req *proto.SignupReq) error {
	if len(req.Email) == 0 || len(req.Password) == 0 {
		return status.New(codes.InvalidArgument, "invalid email or password").Err()
	}

	user, err := c.biz.GetByEmail(ctx, req.Email)
	if err != nil {
		return status.New(codes.Internal, "failed to get user").Err()
	}

	if user != nil && user.Id != 0 {
		return status.New(codes.InvalidArgument, "existed user").Err()
	}

	return nil
}
