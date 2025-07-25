package controller

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/chientranthien/goldenpay/internal/proto"
	"github.com/chientranthien/goldenpay/internal/service/chat/biz"
)

type GetMessagesController struct {
	biz *biz.MessageBiz
}

func NewGetMessagesController(biz *biz.MessageBiz) *GetMessagesController {
	return &GetMessagesController{biz: biz}
}

func (c *GetMessagesController) Do(ctx context.Context, req *proto.GetMessagesReq) (*proto.GetMessagesResp, error) {
	if err := c.validate(req); err != nil {
		return nil, err
	}
	return c.biz.GetMessages(req)
}

func (c *GetMessagesController) validate(req *proto.GetMessagesReq) error {
	if req.ChannelId == 0 {
		return status.New(codes.InvalidArgument, "channel id is required").Err()
	}
	return nil
}
