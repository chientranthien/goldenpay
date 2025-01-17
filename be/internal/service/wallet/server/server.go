package server

import (
	"context"

	"github.com/chientranthien/goldenpay/internal/common"
	commonproto "github.com/chientranthien/goldenpay/internal/common/proto"
	"github.com/chientranthien/goldenpay/internal/proto"
	"github.com/chientranthien/goldenpay/internal/service/wallet/controller"
)

type Server struct {
	proto.UnimplementedWalletServiceServer
	conf                      common.ServiceConfig
	transferController        *controller.TransferController
	processTransferController *controller.ProcessTransferController
	topupController           *controller.TopupController
	getController             *controller.GetController
	createController          *controller.CreateController
	getUserTransactions       *controller.GetUserTransactionsController
}

func NewServer(
	conf common.ServiceConfig,
	transferController *controller.TransferController,
	processTransferController *controller.ProcessTransferController,
	topupController *controller.TopupController,
	getController *controller.GetController,
	createController *controller.CreateController,
	getUserTransactions *controller.GetUserTransactionsController,
) *Server {
	return &Server{
		conf:                      conf,
		transferController:        transferController,
		processTransferController: processTransferController,
		topupController:           topupController,
		getController:             getController,
		createController:          createController,
		getUserTransactions:       getUserTransactions,
	}
}

func (s Server) Get(ctx context.Context, req *proto.GetWalletReq) (*proto.GetWalletResp, error) {
	return s.getController.Do(req)
}

func (s Server) Transfer(ctx context.Context, req *proto.TransferReq) (*proto.TransferResp, error) {
	return s.transferController.Do(req)
}

func (s Server) ProcessTransfer(ctx context.Context, req *proto.ProcessTransferReq) (*proto.ProcessTransferResp, error) {
	return s.processTransferController.Do(req)
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
	server, err := commonproto.NewServer(s.conf.Addr)
	if err != nil {
		common.L().Fatalw("createServerErr", "err", err)
	}
	proto.RegisterWalletServiceServer(
		server,
		s,
	)

	common.FatalIfErr(server.ListenAndServe())
}
