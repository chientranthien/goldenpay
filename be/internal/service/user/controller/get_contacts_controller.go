package controller

import (
	"context"

	"github.com/chientranthien/goldenpay/internal/common"
	"github.com/chientranthien/goldenpay/internal/proto"
	"github.com/chientranthien/goldenpay/internal/service/user/biz"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetContactsController struct {
	biz *biz.UserBiz
}

func NewGetContactsController(biz *biz.UserBiz) *GetContactsController {
	return &GetContactsController{biz: biz}
}

func (c GetContactsController) Do(ctx context.Context, req *proto.GetContactsReq) (*proto.GetContactsResp, error) {
	if err := c.validate(req); err != nil {
		return nil, err
	}

	return c.biz.GetContacts(req)
}

func (c GetContactsController) validate(req *proto.GetContactsReq) error {
	if req.Cond == nil {
		common.L().Errorw("emptyCond", "req", req)
		return status.Error(codes.InvalidArgument, "empty condition")
	}

	if req.Cond.User == nil || req.Cond.User.Eq == 0 {
		common.L().Errorw("invalidUser", "req", req)
		return status.Error(codes.InvalidArgument, "invalid User")
	}

	return nil
}
