package controller

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/chientranthien/goldenpay/internal/proto"
	"github.com/chientranthien/goldenpay/internal/service/chat/biz"
)

type GetChannelController struct {
	biz *biz.ChannelBiz
}

func NewGetChannelController(biz *biz.ChannelBiz) *GetChannelController {
	return &GetChannelController{biz: biz}
}

func (c *GetChannelController) Do(ctx context.Context, req *proto.GetChannelReq) (*proto.GetChannelResp, error) {
	if err := c.validate(req); err != nil {
		return nil, err
	}
	return c.biz.GetChannel(req)
}

func (c *GetChannelController) validate(req *proto.GetChannelReq) error {
	if req.ChannelId == 0 {
		return status.New(codes.InvalidArgument, "channel id is required").Err()
	}
	return nil
}
