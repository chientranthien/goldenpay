package server

import (
	"context"

	"github.com/chientranthien/goldenpay/internal/common"
	commonproto "github.com/chientranthien/goldenpay/internal/common/proto"
	"github.com/chientranthien/goldenpay/internal/proto"
	"github.com/chientranthien/goldenpay/internal/service/chat/controller"
)

type Server struct {
	proto.UnimplementedChatServiceServer
	conf                                 common.ServiceConfig
	getChannelController                 *controller.GetChannelController
	sendMessageController                *controller.SendMessageController
	getMessagesController                *controller.GetMessagesController
	createDirectMessageChannelController *controller.CreateDirectMessageChannelController
	getDirectMessagesController          *controller.GetDirectMessagesController
}

func NewServer(
	conf common.ServiceConfig,
	getChannelController *controller.GetChannelController,
	sendMessageController *controller.SendMessageController,
	getMessagesController *controller.GetMessagesController,
	createDirectMessageChannelController *controller.CreateDirectMessageChannelController,
	getDirectMessagesController *controller.GetDirectMessagesController,
) *Server {
	return &Server{
		conf:                                 conf,
		getChannelController:                 getChannelController,
		sendMessageController:                sendMessageController,
		getMessagesController:                getMessagesController,
		createDirectMessageChannelController: createDirectMessageChannelController,
		getDirectMessagesController:          getDirectMessagesController,
	}
}

// Channel management
func (s *Server) GetChannel(ctx context.Context, req *proto.GetChannelReq) (*proto.GetChannelResp, error) {
	return s.getChannelController.Do(ctx, req)
}

// Message management
func (s *Server) SendMessage(ctx context.Context, req *proto.SendMessageReq) (*proto.SendMessageResp, error) {
	return s.sendMessageController.Do(ctx, req)
}

func (s *Server) GetMessages(ctx context.Context, req *proto.GetMessagesReq) (*proto.GetMessagesResp, error) {
	return s.getMessagesController.Do(ctx, req)
}

// Direct messages
func (s *Server) CreateDirectMessageChannel(ctx context.Context, req *proto.CreateDirectMessageChannelReq) (*proto.CreateDirectMessageChannelResp, error) {
	return s.createDirectMessageChannelController.Do(ctx, req)
}

func (s *Server) GetDirectMessages(ctx context.Context, req *proto.GetDirectMessagesReq) (*proto.GetDirectMessagesResp, error) {
	return s.getDirectMessagesController.Do(ctx, req)
}

func (s *Server) Serve() {
	server, err := commonproto.NewServer(s.conf.Addr)
	if err != nil {
		common.L().Fatalw("createServerErr", "err", err)
	}
	proto.RegisterChatServiceServer(
		server,
		s,
	)

	common.FatalIfErr(server.ListenAndServe())
}
