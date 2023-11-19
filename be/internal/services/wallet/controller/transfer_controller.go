package controller

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/chientranthien/goldenpay/internal/proto"
	"github.com/chientranthien/goldenpay/internal/services/wallet/biz"
)

type TransferController struct {
	biz *biz.WalletBiz
}

func NewTransferController(biz *biz.WalletBiz) *TransferController {
	return &TransferController{biz: biz}
}


func (c TransferController) Transfer(req *proto.TransferReq) (*proto.TransferResp, error) {
	if err := c.validate(req); err != nil {
		return nil, err
	}


	return c.biz.Transfer(req)
}


func (c TransferController) validate(req *proto.TransferReq) error {
	if req.FromUser <= 0 || req.ToUser <= 0 {
		return status.New(codes.InvalidArgument, "invalid from_user or to_user").Err()
	}

	if req.Amount <= 0 {
		return status.New(codes.InvalidArgument, "invalid amount").Err()
	}

	return nil
}
