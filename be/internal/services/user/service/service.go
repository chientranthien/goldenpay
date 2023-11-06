package service

import (
	"context"

	"github.com/chientranthien/goldenpay/internal/proto"
	"github.com/chientranthien/goldenpay/internal/services/user/controller"
)

type Service struct {
	proto.UnimplementedUserServiceServer
	*controller.SignupController
	*controller.LoginController
	*controller.GetController
}

func NewService(
	signupController *controller.SignupController,
	loginController *controller.LoginController,
	getController *controller.GetController,
) *Service {
	return &Service{
		SignupController: signupController,
		LoginController:  loginController,
		GetController:    getController,
	}
}

func (u Service) Signup(ctx context.Context, req *proto.SignupReq) (*proto.SignupResp, error) {
	return u.SignupController.Signup(ctx, req)
}

func (u Service) Login(ctx context.Context, req *proto.LoginReq) (*proto.LoginResp, error) {
	return u.LoginController.Login(ctx, req)
}

func (u Service) Get(ctx context.Context, req *proto.GetUserReq) (*proto.GetUserResp, error) {
	return u.GetController.Get(ctx, req)
}
