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
	*controller.AuthzController
	*controller.GetByEmailController
}

func NewService(
	signupController *controller.SignupController,
	loginController *controller.LoginController,
	getController *controller.GetController,
	authzController *controller.AuthzController,
	getByEmailController *controller.GetByEmailController,
) *Service {
	return &Service{
		SignupController: signupController,
		LoginController: loginController,
		GetController: getController,
		AuthzController: authzController,
		GetByEmailController: getByEmailController,
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

func (u Service) Authz(ctx context.Context, req *proto.AuthzReq) (*proto.AuthzResp, error) {
	return u.AuthzController.Authz(ctx, req)
}

func (u Service) GetByEmail(ctx context.Context, req *proto.GetByEmailReq) (*proto.GetByEmailResp, error) {
	return u.GetByEmailController.GetByEmail(req)
}

