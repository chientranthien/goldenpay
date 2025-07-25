package controller

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/chientranthien/goldenpay/internal/proto"
	"github.com/chientranthien/goldenpay/internal/service/chat/biz"
)

type GetDirectMessagesController struct {
	biz *biz.ChannelBiz
}

func NewGetDirectMessagesController(biz *biz.ChannelBiz) *GetDirectMessagesController {
	return &GetDirectMessagesController{biz: biz}
}

func (c *GetDirectMessagesController) Do(ctx context.Context, req *proto.GetDirectMessagesReq) (*proto.GetDirectMessagesResp, error) {
	if err := c.validate(req); err != nil {
		return nil, err
	}
	return c.biz.GetDirectMessages(req)
}

func (c *GetDirectMessagesController) validate(req *proto.GetDirectMessagesReq) error {
	if req.UserId == 0 {
		return status.New(codes.InvalidArgument, "user id is required").Err()
	}
	return nil
}
