package controller

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/chientranthien/goldenpay/internal/proto"
	"github.com/chientranthien/goldenpay/internal/services/user/biz"
)

type TransferController struct {
	biz  *biz.TransactionBiz
	uBiz *biz.UserBiz
	wBiz *biz.WalletBiz
}

func (c TransferController) Transfer(req *proto.TransferReq) (*proto.TransferResp, error) {
	if err := c.validate(req); err != nil {
		return nil, err
	}
}

func (c TransferController) getAndValidateUser(req *proto.TransferReq) error {
	user, err := c.uBiz.GetByEmail(req.ToEmail)
	if err != nil {
		return status.New(codes.Internal, "failed to get user").Err()
	}

	if user == nil || user.Id == 0 {
		return status.New(codes.InvalidArgument, "invalid to_email").Err()
	}

}
func (c TransferController) validate(req *proto.TransferReq) error {
	if req.ToEmail == "" || req.Amount <= 0 {
		return status.New(codes.InvalidArgument, "invalid to_email or amount").Err()
	}


	return nil
}
