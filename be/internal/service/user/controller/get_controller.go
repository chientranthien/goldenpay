package controller

import (
	"context"

	"github.com/chientranthien/goldenpay/internal/proto"
	"github.com/chientranthien/goldenpay/internal/service/user/biz"
)

type GetController struct {
	biz *biz.UserBiz
}

func NewGetController(biz *biz.UserBiz) *GetController {
	return &GetController{biz: biz}
}

func (c GetController) Get(_ context.Context, req *proto.GetUserReq) (*proto.GetUserResp, error) {
	user, err := c.biz.Get(req.Id)
	if err != nil {
		return nil, err
	}
	resp := &proto.GetUserResp{
		User: (*proto.User)(user),
	}

	return resp, nil
}
