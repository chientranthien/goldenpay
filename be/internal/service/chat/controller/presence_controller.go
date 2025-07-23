package controller

import (
	"context"

	"github.com/chientranthien/goldenpay/internal/proto"
	"github.com/chientranthien/goldenpay/internal/service/chat/biz"
)

type PresenceController struct {
	biz *biz.ChatBiz
}

func NewPresenceController(biz *biz.ChatBiz) *PresenceController {
	return &PresenceController{biz: biz}
}

func (c *PresenceController) UpdatePresence(ctx context.Context, req *proto.UpdatePresenceReq) (*proto.UpdatePresenceResp, error) {
	return c.biz.UpdatePresence(req)
}

func (c *PresenceController) GetPresence(ctx context.Context, req *proto.GetPresenceReq) (*proto.GetPresenceResp, error) {
	return c.biz.GetPresence(req)
}
