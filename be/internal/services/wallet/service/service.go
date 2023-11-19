package service

import (
	"context"

	"github.com/chientranthien/goldenpay/internal/proto"
	"github.com/chientranthien/goldenpay/internal/services/wallet/controller"
)

type Service struct {
	proto.UnimplementedWalletServiceServer
	transferController *controller.TransferController
	getController      *controller.GetController
}

func NewService(
	transferController *controller.TransferController,
	getController *controller.GetController,
) *Service {
	return &Service{transferController: transferController, getController: getController}
}

func (s Service) Get(ctx context.Context, req *proto.GetWalletReq) (*proto.GetWalletResp, error) {
	return s.getController.Get(req)
}

func (s Service) Transfer(ctx context.Context, req *proto.TransferReq) (*proto.TransferResp, error) {
	return s.transferController.Transfer(req)
}
