package controller

import (
	"context"

	"github.com/chientranthien/goldenpay/internal/proto"
	"github.com/chientranthien/goldenpay/internal/service/chat/biz"
)

type MessageController struct {
	biz *biz.ChatBiz
}

func NewMessageController(biz *biz.ChatBiz) *MessageController {
	return &MessageController{biz: biz}
}

func (c *MessageController) SendMessage(ctx context.Context, req *proto.SendMessageReq) (*proto.SendMessageResp, error) {
	return c.biz.SendMessage(req)
}

func (c *MessageController) GetMessages(ctx context.Context, req *proto.GetMessagesReq, userId uint64) (*proto.GetMessagesResp, error) {
	return c.biz.GetMessages(req, userId)
}

func (c *MessageController) UpdateMessage(ctx context.Context, req *proto.UpdateMessageReq) (*proto.UpdateMessageResp, error) {
	return c.biz.UpdateMessage(req)
}

func (c *MessageController) DeleteMessage(ctx context.Context, req *proto.DeleteMessageReq) (*proto.DeleteMessageResp, error) {
	return c.biz.DeleteMessage(req)
}

func (c *MessageController) CreateDirectMessage(ctx context.Context, req *proto.CreateDirectMessageReq) (*proto.CreateDirectMessageResp, error) {
	return c.biz.CreateDirectMessage(req)
}

func (c *MessageController) GetDirectMessages(ctx context.Context, req *proto.GetDirectMessagesReq) (*proto.GetDirectMessagesResp, error) {
	return c.biz.GetDirectMessages(req)
}
