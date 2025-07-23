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
	conf               common.ServiceConfig
	channelController  *controller.ChannelController
	messageController  *controller.MessageController
	presenceController *controller.PresenceController
}

func NewServer(
	conf common.ServiceConfig,
	channelController *controller.ChannelController,
	messageController *controller.MessageController,
	presenceController *controller.PresenceController,
) *Server {
	return &Server{
		conf:               conf,
		channelController:  channelController,
		messageController:  messageController,
		presenceController: presenceController,
	}
}

// Channel management
func (s *Server) CreateChannel(ctx context.Context, req *proto.CreateChannelReq) (*proto.CreateChannelResp, error) {
	// Extract user ID from context (set by HTTP gateway authentication)
	userId := common.GetUserIdFromCtx(ctx)
	return s.channelController.CreateChannel(ctx, req, userId)
}

func (s *Server) GetChannel(ctx context.Context, req *proto.GetChannelReq) (*proto.GetChannelResp, error) {
	userId := common.GetUserIdFromCtx(ctx)
	return s.channelController.GetChannel(ctx, req, userId)
}

func (s *Server) ListChannels(ctx context.Context, req *proto.ListChannelsReq) (*proto.ListChannelsResp, error) {
	return s.channelController.ListChannels(ctx, req)
}

func (s *Server) JoinChannel(ctx context.Context, req *proto.JoinChannelReq) (*proto.JoinChannelResp, error) {
	return s.channelController.JoinChannel(ctx, req)
}

func (s *Server) LeaveChannel(ctx context.Context, req *proto.LeaveChannelReq) (*proto.LeaveChannelResp, error) {
	return s.channelController.LeaveChannel(ctx, req)
}

// Message management
func (s *Server) SendMessage(ctx context.Context, req *proto.SendMessageReq) (*proto.SendMessageResp, error) {
	return s.messageController.SendMessage(ctx, req)
}

func (s *Server) GetMessages(ctx context.Context, req *proto.GetMessagesReq) (*proto.GetMessagesResp, error) {
	userId := common.GetUserIdFromCtx(ctx)
	return s.messageController.GetMessages(ctx, req, userId)
}

func (s *Server) UpdateMessage(ctx context.Context, req *proto.UpdateMessageReq) (*proto.UpdateMessageResp, error) {
	return s.messageController.UpdateMessage(ctx, req)
}

func (s *Server) DeleteMessage(ctx context.Context, req *proto.DeleteMessageReq) (*proto.DeleteMessageResp, error) {
	return s.messageController.DeleteMessage(ctx, req)
}

// Direct messages
func (s *Server) CreateDirectMessage(ctx context.Context, req *proto.CreateDirectMessageReq) (*proto.CreateDirectMessageResp, error) {
	return s.messageController.CreateDirectMessage(ctx, req)
}

func (s *Server) GetDirectMessages(ctx context.Context, req *proto.GetDirectMessagesReq) (*proto.GetDirectMessagesResp, error) {
	return s.messageController.GetDirectMessages(ctx, req)
}

// User presence
func (s *Server) UpdatePresence(ctx context.Context, req *proto.UpdatePresenceReq) (*proto.UpdatePresenceResp, error) {
	return s.presenceController.UpdatePresence(ctx, req)
}

func (s *Server) GetPresence(ctx context.Context, req *proto.GetPresenceReq) (*proto.GetPresenceResp, error) {
	return s.presenceController.GetPresence(ctx, req)
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
