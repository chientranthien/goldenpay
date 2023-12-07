package server

import (
	"context"
	"net"

	"google.golang.org/grpc"

	"github.com/chientranthien/goldenpay/internal/common"
	"github.com/chientranthien/goldenpay/internal/proto"
	"github.com/chientranthien/goldenpay/internal/service/user/controller"
)

type Server struct {
	proto.UnimplementedUserServiceServer
	conf                 common.ServiceConfig
	signupController     *controller.SignupController
	loginController      *controller.LoginController
	getController        *controller.GetController
	getBatchController   *controller.GetBatchController
	authzController      *controller.AuthzController
	getByEmailController *controller.GetByEmailController
}

func NewServer(
	conf common.ServiceConfig,
	signupController *controller.SignupController,
	loginController *controller.LoginController,
	getController *controller.GetController,
	getBatchController *controller.GetBatchController,
	authzController *controller.AuthzController,
	getByEmailController *controller.GetByEmailController,
) *Server {
	return &Server{
		conf: conf,
		signupController: signupController,
		loginController: loginController,
		getController: getController,
		getBatchController: getBatchController,
		authzController: authzController,
		getByEmailController: getByEmailController,
	}
}

func (s Server) Signup(ctx context.Context, req *proto.SignupReq) (*proto.SignupResp, error) {
	return s.signupController.Signup(ctx, req)
}

func (s Server) Login(ctx context.Context, req *proto.LoginReq) (*proto.LoginResp, error) {
	return s.loginController.Login(ctx, req)
}

func (s Server) Get(ctx context.Context, req *proto.GetReq) (*proto.GetResp, error) {
	return s.getController.Do(ctx, req)
}

func (s Server) GetBatch(ctx context.Context, req *proto.GetBatchReq) (*proto.GetBatchResp, error) {
	return s.getBatchController.Do(ctx, req)
}

func (s Server) Authz(ctx context.Context, req *proto.AuthzReq) (*proto.AuthzResp, error) {
	return s.authzController.Authz(ctx, req)
}

func (s Server) GetByEmail(ctx context.Context, req *proto.GetByEmailReq) (*proto.GetByEmailResp, error) {
	return s.getByEmailController.GetByEmail(req)
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
