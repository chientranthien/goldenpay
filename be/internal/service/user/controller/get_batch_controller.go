package controller

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/chientranthien/goldenpay/internal/proto"
	"github.com/chientranthien/goldenpay/internal/service/user/biz"
)

type (
	GetBatchController struct {
		biz *biz.UserBiz
	}
)

func NewGetBatchController(biz *biz.UserBiz) *GetBatchController {
	return &GetBatchController{biz: biz}
}

func (c GetBatchController) Do(ctx context.Context, req *proto.GetBatchReq) (*proto.GetBatchResp, error) {
	if err := c.validate(req); err != nil {
		return nil, err
	}

	return c.biz.GetBatch(req.Ids)
}

func (c GetBatchController) validate(req *proto.GetBatchReq) error {
	if len(req.Ids) == 0 {
		return status.New(codes.InvalidArgument, "empty ids").Err()
	}

	if len(req.Ids) > 100 {
		return status.New(codes.InvalidArgument, "exceed max len").Err()
	}

	return nil
}
