package controller

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/chientranthien/goldenpay/internal/proto"
	"github.com/chientranthien/goldenpay/internal/service/user/biz"
)

type LoginController struct {
	biz *biz.UserBiz
}

func NewLoginController(biz *biz.UserBiz) *LoginController {
	return &LoginController{biz: biz}
}

func (c LoginController) Login(_ context.Context, req *proto.LoginReq) (*proto.LoginResp, error) {
	if err := c.validate(req); err != nil {
		return nil, err
	}


	return c.biz.Login(req)
}

func (c LoginController) validate(req *proto.LoginReq) error {
	if len(req.Email) == 0 || len(req.Password) == 0 {
		return status.New(codes.InvalidArgument, "invalid email or password").Err()
	}

	return nil
}
