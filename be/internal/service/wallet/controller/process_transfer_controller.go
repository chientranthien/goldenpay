package controller

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/chientranthien/goldenpay/internal/common"
	"github.com/chientranthien/goldenpay/internal/proto"
	"github.com/chientranthien/goldenpay/internal/service/wallet/biz"
)

type ProcessTransferController struct {
	biz *biz.WalletBiz
}

func NewProcessTransferController(biz *biz.WalletBiz) *ProcessTransferController {
	return &ProcessTransferController{biz: biz}
}

func (c ProcessTransferController) Do(req *proto.ProcessTransferReq) (*proto.ProcessTransferResp, error) {
	if err := c.validate(req); err != nil {
		return nil, err
	}

	return c.biz.ProcessTransfer(req)
}

func (c ProcessTransferController) validate(req *proto.ProcessTransferReq) error {
	if req.TransactionId <= 0 {
		common.L().Errorw("invalidTransactionId", "transactionId", req.TransactionId)
		return status.Error(codes.InvalidArgument, "invalid transactionID")
	}

	return nil
}
