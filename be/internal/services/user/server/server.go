package server

import (
	"context"
	"net"

	"google.golang.org/grpc"

	"github.com/chientranthien/goldenpay/internal/common"
	"github.com/chientranthien/goldenpay/internal/proto"
	"github.com/chientranthien/goldenpay/internal/services/user/controller"
)

type Server struct {
	proto.UnimplementedUserServiceServer
	conf common.ServiceConfig
	*controller.SignupController
	*controller.LoginController
	*controller.GetController
	*controller.AuthzController
	*controller.GetByEmailController
}

func NewServer(
	conf common.ServiceConfig,
	signupController *controller.SignupController,
	loginController *controller.LoginController,
	getController *controller.GetController,
	authzController *controller.AuthzController,
	getByEmailController *controller.GetByEmailController,
) *Server {
	return &Server{
		conf:                 conf,
		SignupController:     signupController,
		LoginController:      loginController,
		GetController:        getController,
		AuthzController:      authzController,
		GetByEmailController: getByEmailController,
	}
}

func (s Server) Signup(ctx context.Context, req *proto.SignupReq) (*proto.SignupResp, error) {
	return s.SignupController.Signup(ctx, req)
}

func (s Server) Login(ctx context.Context, req *proto.LoginReq) (*proto.LoginResp, error) {
	return s.LoginController.Login(ctx, req)
}

func (s Server) Get(ctx context.Context, req *proto.GetUserReq) (*proto.GetUserResp, error) {
	return s.GetController.Get(ctx, req)
}

func (s Server) Authz(ctx context.Context, req *proto.AuthzReq) (*proto.AuthzResp, error) {
	return s.AuthzController.Authz(ctx, req)
}

func (s Server) GetByEmail(ctx context.Context, req *proto.GetByEmailReq) (*proto.GetByEmailResp, error) {
	return s.GetByEmailController.GetByEmail(req)
}

func (s Server) Serve() {
	server := grpc.NewServer(
		common.ServerLoggingInterceptor,
	)
	proto.RegisterUserServiceServer(
		server,
		s,
	)

	lis, err := net.Listen("tcp", s.conf.Addr)
	if err != nil {
		common.L().Fatalw("netListenErr", "config", s.conf, "err", err)
	}
	err = server.Serve(lis)
	if err != nil {
		common.L().Fatalw("serveErr", "err", err)
	}
}
