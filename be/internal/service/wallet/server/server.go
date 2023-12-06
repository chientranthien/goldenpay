package server

import (
	"context"
	"net"

	"google.golang.org/grpc"

	"github.com/chientranthien/goldenpay/internal/common"
	"github.com/chientranthien/goldenpay/internal/proto"
	"github.com/chientranthien/goldenpay/internal/service/wallet/controller"
)

type Server struct {
	proto.UnimplementedWalletServiceServer
	conf                common.ServiceConfig
	transferController  *controller.TransferController
	topupController     *controller.TopupController
	getController       *controller.GetController
	createController    *controller.CreateController
	getUserTransactions *controller.GetUserTransactionsController
}

func NewServer(
	conf common.ServiceConfig,
	transferController *controller.TransferController,
	topupController *controller.TopupController,
	getController *controller.GetController,
	createController *controller.CreateController,
	getUserTransactions *controller.GetUserTransactionsController,
) *Server {
	return &Server{
		conf:                conf,
		transferController:  transferController,
		topupController:     topupController,
		getController:       getController,
		createController:    createController,
		getUserTransactions: getUserTransactions,
	}
}

func (s Server) Get(ctx context.Context, req *proto.GetWalletReq) (*proto.GetWalletResp, error) {
	return s.getController.Do(req)
}

func (s Server) Transfer(ctx context.Context, req *proto.TransferReq) (*proto.TransferResp, error) {
	return s.transferController.Do(req)
}

func (s Server) Topup(ctx context.Context, req *proto.TopupReq) (*proto.TopupResp, error) {
	return s.topupController.Do(req)
}

func (s Server) Create(ctx context.Context, req *proto.CreateWalletReq) (*proto.CreateWalletResp, error) {
	return s.createController.Create(req)
}

func (s Server) GetUserTransactions(ctx context.Context, req *proto.GetUserTransactionsReq) (*proto.GetUserTransactionsResp, error) {
	return s.getUserTransactions.Do(req)
}

func (s Server) Serve() {
	server := grpc.NewServer(
		common.ServerLoggingInterceptor,
	)
	proto.RegisterWalletServiceServer(
		server,
		s,
	)

	lis, err := net.Listen("tcp", s.conf.Addr)
	if err != nil {
		common.L().Fatalw("netListenErr", "config", s.conf, "err", err.Error())
	} else {
		common.L().Infow("listening", "add", s.conf.Addr)
	}
	err = server.Serve(lis)
	if err != nil {
		common.L().Fatalw("serveErr", "err", err)
	}

}
