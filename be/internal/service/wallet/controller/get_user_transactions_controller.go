package controller

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/chientranthien/goldenpay/internal/common"
	"github.com/chientranthien/goldenpay/internal/proto"
	"github.com/chientranthien/goldenpay/internal/service/wallet/biz"
)

type GetUserTransactionsController struct {
	biz *biz.WalletBiz
}

func NewGetUserTransactionsController(biz *biz.WalletBiz) *GetUserTransactionsController {
	return &GetUserTransactionsController{biz: biz}
}

func (c GetUserTransactionsController) Do(req *proto.GetUserTransactionsReq) (*proto.GetUserTransactionsResp, error) {
	if err := c.validate(req); err != nil {
		return nil, err
	}

	return c.biz.GetUserTransactions(req)
}
func (c GetUserTransactionsController) validate(req *proto.GetUserTransactionsReq) error {
	if req.Cond == nil{
		common.L().Errorw("emptyCond", "req", req)
		return status.Error(codes.InvalidArgument, "empty condition")
	}

	if req.Cond.User == nil || req.Cond.User.Eq == 0{
		common.L().Errorw("invalidUser", "req", req)
		return status.Error(codes.InvalidArgument, "invalid User")

	}

	return nil
}
