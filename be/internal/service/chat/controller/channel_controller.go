package controller

import (
	"context"

	"github.com/chientranthien/goldenpay/internal/proto"
	"github.com/chientranthien/goldenpay/internal/service/chat/biz"
)

type ChannelController struct {
	biz *biz.ChatBiz
}

func NewChannelController(biz *biz.ChatBiz) *ChannelController {
	return &ChannelController{biz: biz}
}

func (c *ChannelController) CreateChannel(ctx context.Context, req *proto.CreateChannelReq, userId uint64) (*proto.CreateChannelResp, error) {
	return c.biz.CreateChannel(req, userId)
}

func (c *ChannelController) GetChannel(ctx context.Context, req *proto.GetChannelReq, userId uint64) (*proto.GetChannelResp, error) {
	return c.biz.GetChannel(req, userId)
}

func (c *ChannelController) ListChannels(ctx context.Context, req *proto.ListChannelsReq) (*proto.ListChannelsResp, error) {
	return c.biz.ListChannels(req)
}

func (c *ChannelController) JoinChannel(ctx context.Context, req *proto.JoinChannelReq) (*proto.JoinChannelResp, error) {
	return c.biz.JoinChannel(req)
}

func (c *ChannelController) LeaveChannel(ctx context.Context, req *proto.LeaveChannelReq) (*proto.LeaveChannelResp, error) {
	return c.biz.LeaveChannel(req)
}
