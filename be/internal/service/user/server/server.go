package server

import (
	"context"

	"github.com/chientranthien/goldenpay/internal/common"
	commonproto "github.com/chientranthien/goldenpay/internal/common/proto"
	"github.com/chientranthien/goldenpay/internal/proto"
	"github.com/chientranthien/goldenpay/internal/service/user/controller"
)

type Server struct {
	proto.UnimplementedUserServiceServer
	conf                              common.ServiceConfig
	signupController                  *controller.SignupController
	loginController                   *controller.LoginController
	getController                     *controller.GetController
	getBatchController                *controller.GetBatchController
	authzController                   *controller.AuthzController
	getByEmailController              *controller.GetByEmailController
	createContactIfNotExistController *controller.CreateContactIfNotExistController
	getContactsController             *controller.GetContactsController
}

func NewServer(
	conf common.ServiceConfig,
	signupController *controller.SignupController,
	loginController *controller.LoginController,
	getController *controller.GetController,
	getBatchController *controller.GetBatchController,
	authzController *controller.AuthzController,
	getByEmailController *controller.GetByEmailController,
	createContactIfNotExistController *controller.CreateContactIfNotExistController,
	getContactsController *controller.GetContactsController,
) *Server {
	return &Server{
		conf:                              conf,
		signupController:                  signupController,
		loginController:                   loginController,
		getController:                     getController,
		getBatchController:                getBatchController,
		authzController:                   authzController,
		getByEmailController:              getByEmailController,
		createContactIfNotExistController: createContactIfNotExistController,
		getContactsController:             getContactsController,
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
	return s.authzController.Do(ctx, req)
}

func (s Server) GetByEmail(ctx context.Context, req *proto.GetByEmailReq) (*proto.GetByEmailResp, error) {
	return s.getByEmailController.Do(req)
}

func (s Server) CreateContactIfNotExist(ctx context.Context, req *proto.CreateContactIfNotExistReq) (*proto.CreateContactIfNotExistResp, error) {
	return s.createContactIfNotExistController.Do(req)
}

func (s Server) GetContacts(ctx context.Context, req *proto.GetContactsReq) (*proto.GetContactsResp, error) {
	return s.getContactsController.Do(ctx, req)
}

func (s Server) Serve() {
	server, err := commonproto.NewServer(s.conf.Addr)
	if err != nil {
		common.L().Fatalw("createServerErr", "err", err)
	}
	proto.RegisterUserServiceServer(
		server,
		s,
	)

	common.FatalIfErr(server.ListenAndServe())
}
