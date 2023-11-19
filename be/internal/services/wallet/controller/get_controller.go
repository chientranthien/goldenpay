package controller

import (
	"github.com/chientranthien/goldenpay/internal/proto"
	"github.com/chientranthien/goldenpay/internal/services/wallet/biz"
)

type GetController struct {
	biz *biz.WalletBiz
}

func NewGetController(biz *biz.WalletBiz) *GetController {
	return &GetController{biz: biz}
}

func (c GetController) Get(req *proto.GetWalletReq) (*proto.GetWalletResp, error) {
	return nil, nil
}
