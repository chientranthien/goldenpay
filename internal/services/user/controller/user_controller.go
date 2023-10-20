package controller

import (
	"context"

	"github.com/chientranthien/goldenpay/internal/proto"
)

type UserController struct {
	proto.UnimplementedUserServiceServer
}

func (u UserController) Signup(ctx context.Context, req *proto.SignupReq) (*proto.SignupResp, error) {
	panic("implement me")
}

func (u UserController) Login(ctx context.Context, req *proto.LoginReq) (*proto.LoginResp, error) {
	panic("implement me")
}

func (u UserController) Get(ctx context.Context, req *proto.GetUserReq) (*proto.GetUserResp, error) {
	panic("implement me")
}
