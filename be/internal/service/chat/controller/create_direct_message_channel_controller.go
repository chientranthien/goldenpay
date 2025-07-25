package controller

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/chientranthien/goldenpay/internal/proto"
	"github.com/chientranthien/goldenpay/internal/service/chat/biz"
)

type CreateDirectMessageChannelController struct {
	biz *biz.ChannelBiz
}

func NewCreateDirectMessageChannelController(biz *biz.ChannelBiz) *CreateDirectMessageChannelController {
	return &CreateDirectMessageChannelController{biz: biz}
}

func (c *CreateDirectMessageChannelController) Do(ctx context.Context, req *proto.CreateDirectMessageChannelReq) (*proto.CreateDirectMessageChannelResp, error) {
	if err := c.validate(req); err != nil {
		return nil, err
	}
	return c.biz.CreateDirectMessageChannel(req)
}

func (c *CreateDirectMessageChannelController) validate(req *proto.CreateDirectMessageChannelReq) error {
	if req.FromUser == 0 {
		return status.New(codes.InvalidArgument, "from user id is required").Err()
	}
	if req.ToUser == 0 {
		return status.New(codes.InvalidArgument, "to user id is required").Err()
	}
	if req.FromUser == req.ToUser {
		return status.New(codes.InvalidArgument, "cannot create direct message to yourself").Err()
	}
	return nil
}
