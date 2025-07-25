package controller

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/chientranthien/goldenpay/internal/proto"
	"github.com/chientranthien/goldenpay/internal/service/chat/biz"
)

type SendMessageController struct {
	biz *biz.MessageBiz
}

func NewSendMessageController(biz *biz.MessageBiz) *SendMessageController {
	return &SendMessageController{biz: biz}
}

func (c *SendMessageController) Do(ctx context.Context, req *proto.SendMessageReq) (*proto.SendMessageResp, error) {
	if err := c.validate(req); err != nil {
		return nil, err
	}
	return c.biz.SendMessage(req)
}

func (c *SendMessageController) validate(req *proto.SendMessageReq) error {
	if req.ChannelId == 0 {
		return status.New(codes.InvalidArgument, "channel id is required").Err()
	}
	if req.UserId == 0 {
		return status.New(codes.InvalidArgument, "user id is required").Err()
	}
	if len(req.Content) == 0 {
		return status.New(codes.InvalidArgument, "message content is required").Err()
	}
	if len(req.MessageType) == 0 {
		return status.New(codes.InvalidArgument, "message type is required").Err()
	}
	return nil
}
