package controller

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/chientranthien/goldenpay/internal/proto"
	"github.com/chientranthien/goldenpay/internal/services/user/biz"
)

type LoginController struct {
	biz *biz.UserBiz
}

func NewLoginController(biz *biz.UserBiz) *LoginController {
	return &LoginController{biz: biz}
}

func (c LoginController) Login(ctx context.Context, req *proto.LoginReq) (*proto.LoginResp, error) {
	if err := c.validate(req); err != nil {
		return nil, err
	}

	user, _ := c.biz.GetByUsername(ctx, req.Username)

	if user.HashedPassword != c.biz.HashPassword(req.Password) {
		return nil, status.New(codes.InvalidArgument, "incorrect password").Err()
	}

	token, err := c.biz.GenerateToken(user)
	if err != nil {
		return nil, status.New(codes.Internal, "unable to generate token").Err()
	}

	return &proto.LoginResp{Token: token}, nil
}

func (c LoginController) validate(req *proto.LoginReq) error {
	if len(req.Username) == 0 || len(req.Password) == 0 {
		return status.New(codes.InvalidArgument, "invalid username or password").Err()
	}

	return nil
}