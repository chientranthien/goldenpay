package biz

import (
	"fmt"

	"github.com/IBM/sarama"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/chientranthien/goldenpay/internal/common"
	"github.com/chientranthien/goldenpay/internal/proto"
	"github.com/chientranthien/goldenpay/internal/service/chat/dao"
)

type ChatBiz struct {
	dao                      *dao.ChatDao
	newMessageProducer       sarama.SyncProducer
	newMessageProducerConfig common.ProducerConfig
	presenceProducer         sarama.SyncProducer
	presenceProducerConfig   common.ProducerConfig
	membershipProducer       sarama.SyncProducer
	membershipProducerConfig common.ProducerConfig
}

func NewChatBiz(
	dao *dao.ChatDao,
	newMessageProducer sarama.SyncProducer,
	newMessageProducerConfig common.ProducerConfig,
	presenceProducer sarama.SyncProducer,
	presenceProducerConfig common.ProducerConfig,
	membershipProducer sarama.SyncProducer,
	membershipProducerConfig common.ProducerConfig,
) *ChatBiz {
	return &ChatBiz{
		dao:                      dao,
		newMessageProducer:       newMessageProducer,
		newMessageProducerConfig: newMessageProducerConfig,
		presenceProducer:         presenceProducer,
		presenceProducerConfig:   presenceProducerConfig,
		membershipProducer:       membershipProducer,
		membershipProducerConfig: membershipProducerConfig,
	}
}

// Channel operations
func (b *ChatBiz) CreateChannel(req *proto.CreateChannelReq, creatorId uint64) (*proto.CreateChannelResp, error) {
	channel := &proto.Channel{
		Name:        req.Name,
		Description: req.Description,
		Type:        common.ChannelTypeDirect,
		CreatorId:   creatorId,
		Status:      common.ChannelStatusActive,
		Version:     common.FirstVersion,
		Ctime:       common.NowMillis(),
		Mtime:       common.NowMillis(),
	}

	dao, err := b.dao.Begin()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to start transaction: %v", err)
	}
	defer dao.CommitOrRollback(err == nil)

	err = dao.CreateChannel(channel)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create channel: %v", err)
	}

	// Add creator as admin member
	creatorMember := &proto.ChannelMember{
		ChannelId:  channel.Id,
		UserId:     creatorId,
		Role:       "admin",
		JoinedAt:   common.NowMillis(),
		LastReadAt: common.NowMillis(),
	}

	err = dao.AddChannelMember(creatorMember)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add creator as member: %v", err)
	}

	// Add other members
	for _, memberUserId := range req.MemberIds {
		if memberUserId == creatorId {
			continue // skip creator, already added
		}

		member := &proto.ChannelMember{
			ChannelId:  channel.Id,
			UserId:     memberUserId,
			Role:       "member",
			JoinedAt:   common.NowMillis(),
			LastReadAt: common.NowMillis(),
		}

		err = dao.AddChannelMember(member)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to add member %d: %v", memberUserId, err)
		}

		// Send membership event
		b.sendMembershipEvent(channel.Id, memberUserId, "joined")
	}

	return &proto.CreateChannelResp{ChannelId: channel.Id}, nil
}

func (b *ChatBiz) GetChannel(req *proto.GetChannelReq, userId uint64) (*proto.GetChannelResp, error) {
	// Check if user is member of the channel
	isMember, err := b.dao.IsUserInChannel(req.ChannelId, userId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check membership: %v", err)
	}
	if !isMember {
		return nil, status.Errorf(codes.PermissionDenied, "user is not a member of this channel")
	}

	channel, err := b.dao.GetChannel(req.ChannelId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "channel not found: %v", err)
	}

	return &proto.GetChannelResp{Channel: channel}, nil
}

func (b *ChatBiz) ListChannels(req *proto.ListChannelsReq) (*proto.ListChannelsResp, error) {
	channels, pagination, err := b.dao.ListChannelsForUser(req.UserId, req.Pagination)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list channels: %v", err)
	}

	return &proto.ListChannelsResp{
		Channels:       channels,
		NextPagination: pagination,
	}, nil
}

func (b *ChatBiz) JoinChannel(req *proto.JoinChannelReq) (*proto.JoinChannelResp, error) {
	// Check if channel exists and is public
	channel, err := b.dao.GetChannel(req.ChannelId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "channel not found: %v", err)
	}

	if channel.Type == common.ChannelTypePrivate {
		return nil, status.Errorf(codes.PermissionDenied, "cannot join private channel")
	}

	// Check if user is already a member
	isMember, err := b.dao.IsUserInChannel(req.ChannelId, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check membership: %v", err)
	}
	if isMember {
		return nil, status.Errorf(codes.AlreadyExists, "user is already a member")
	}

	member := &proto.ChannelMember{
		ChannelId:  req.ChannelId,
		UserId:     req.UserId,
		Role:       "member",
		JoinedAt:   common.NowMillis(),
		LastReadAt: common.NowMillis(),
	}

	err = b.dao.AddChannelMember(member)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add member: %v", err)
	}

	// Send membership event
	b.sendMembershipEvent(req.ChannelId, req.UserId, "joined")

	return &proto.JoinChannelResp{}, nil
}

func (b *ChatBiz) LeaveChannel(req *proto.LeaveChannelReq) (*proto.LeaveChannelResp, error) {
	// Check if user is a member
	isMember, err := b.dao.IsUserInChannel(req.ChannelId, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check membership: %v", err)
	}
	if !isMember {
		return nil, status.Errorf(codes.NotFound, "user is not a member of this channel")
	}

	err = b.dao.RemoveChannelMember(req.ChannelId, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to remove member: %v", err)
	}

	// Send membership event
	b.sendMembershipEvent(req.ChannelId, req.UserId, "left")

	return &proto.LeaveChannelResp{}, nil
}

// Message operations
func (b *ChatBiz) SendMessage(req *proto.SendMessageReq) (*proto.SendMessageResp, error) {
	// Check if user is member of the channel
	isMember, err := b.dao.IsUserInChannel(req.ChannelId, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check membership: %v", err)
	}
	if !isMember {
		return nil, status.Errorf(codes.PermissionDenied, "user is not a member of this channel")
	}

	message := &proto.Message{
		ChannelId:   req.ChannelId,
		UserId:      req.UserId,
		Content:     req.Content,
		MessageType: req.MessageType,
		Metadata:    req.Metadata,
		Status:      common.MessageStatusActive,
		Version:     common.FirstVersion,
		Ctime:       common.NowMillis(),
		Mtime:       common.NowMillis(),
	}

	if req.MessageType == "" {
		message.MessageType = common.MessageTypeText
	}

	err = b.dao.CreateMessage(message)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create message: %v", err)
	}

	// Send new message event
	b.sendNewMessageEvent(message.Id, req.ChannelId, req.UserId)

	// Update last read timestamp for the sender
	b.dao.UpdateLastReadAt(req.ChannelId, req.UserId, message.Ctime)

	return &proto.SendMessageResp{MessageId: message.Id}, nil
}

func (b *ChatBiz) GetMessages(req *proto.GetMessagesReq, userId uint64) (*proto.GetMessagesResp, error) {
	// Check if user is member of the channel
	isMember, err := b.dao.IsUserInChannel(req.ChannelId, userId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check membership: %v", err)
	}
	if !isMember {
		return nil, status.Errorf(codes.PermissionDenied, "user is not a member of this channel")
	}

	messages, pagination, err := b.dao.GetMessages(req.ChannelId, req.Pagination)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get messages: %v", err)
	}

	return &proto.GetMessagesResp{
		Messages:       messages,
		NextPagination: pagination,
	}, nil
}

func (b *ChatBiz) UpdateMessage(req *proto.UpdateMessageReq) (*proto.UpdateMessageResp, error) {
	message, err := b.dao.GetMessage(req.MessageId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "message not found: %v", err)
	}

	if message.UserId != req.UserId {
		return nil, status.Errorf(codes.PermissionDenied, "user can only update their own messages")
	}

	message.Content = req.Content
	message.Status = common.MessageStatusEdited
	message.Mtime = common.NowMillis()

	err = b.dao.UpdateMessage(message)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update message: %v", err)
	}

	return &proto.UpdateMessageResp{}, nil
}

func (b *ChatBiz) DeleteMessage(req *proto.DeleteMessageReq) (*proto.DeleteMessageResp, error) {
	message, err := b.dao.GetMessage(req.MessageId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "message not found: %v", err)
	}

	if message.UserId != req.UserId {
		return nil, status.Errorf(codes.PermissionDenied, "user can only delete their own messages")
	}

	err = b.dao.DeleteMessage(req.MessageId, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete message: %v", err)
	}

	return &proto.DeleteMessageResp{}, nil
}

// Direct message operations
func (b *ChatBiz) CreateDirectMessage(req *proto.CreateDirectMessageReq) (*proto.CreateDirectMessageResp, error) {
	// Check if direct message channel already exists
	existingChannel, err := b.dao.FindDirectMessageChannel(req.FromUser, req.ToUser)
	if err == nil {
		// Channel already exists
		return &proto.CreateDirectMessageResp{ChannelId: existingChannel.Id}, nil
	}

	// Create new direct message channel
	channelName := fmt.Sprintf("dm_%d_%d", req.FromUser, req.ToUser)
	channel := &proto.Channel{
		Name:      channelName,
		Type:      common.ChannelTypeDirect,
		CreatorId: req.FromUser,
		Status:    common.ChannelStatusActive,
		Version:   common.FirstVersion,
		Ctime:     common.NowMillis(),
		Mtime:     common.NowMillis(),
	}

	dao, err := b.dao.Begin()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to start transaction: %v", err)
	}
	defer dao.CommitOrRollback(err == nil)

	err = dao.CreateChannel(channel)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create direct message channel: %v", err)
	}

	// Add both users as members
	members := []*proto.ChannelMember{
		{
			ChannelId:  channel.Id,
			UserId:     req.FromUser,
			Role:       "member",
			JoinedAt:   common.NowMillis(),
			LastReadAt: common.NowMillis(),
		},
		{
			ChannelId:  channel.Id,
			UserId:     req.ToUser,
			Role:       "member",
			JoinedAt:   common.NowMillis(),
			LastReadAt: common.NowMillis(),
		},
	}

	for _, member := range members {
		err = dao.AddChannelMember(member)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to add member: %v", err)
		}
	}

	return &proto.CreateDirectMessageResp{ChannelId: channel.Id}, nil
}

func (b *ChatBiz) GetDirectMessages(req *proto.GetDirectMessagesReq) (*proto.GetDirectMessagesResp, error) {
	channels, pagination, err := b.dao.GetDirectMessageChannels(req.UserId, req.Pagination)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get direct messages: %v", err)
	}

	return &proto.GetDirectMessagesResp{
		Channels:       channels,
		NextPagination: pagination,
	}, nil
}

// Presence operations
func (b *ChatBiz) UpdatePresence(req *proto.UpdatePresenceReq) (*proto.UpdatePresenceResp, error) {
	presence := &proto.UserPresence{
		UserId:       req.UserId,
		Status:       req.Status,
		StatusText:   req.StatusText,
		LastActivity: common.NowMillis(),
	}

	err := b.dao.UpsertUserPresence(presence)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update presence: %v", err)
	}

	// Send presence update event
	b.sendPresenceEvent(req.UserId, req.Status, req.StatusText)

	return &proto.UpdatePresenceResp{}, nil
}

func (b *ChatBiz) GetPresence(req *proto.GetPresenceReq) (*proto.GetPresenceResp, error) {
	presences, err := b.dao.GetUserPresence(req.UserIds)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get presence: %v", err)
	}

	return &proto.GetPresenceResp{Presences: presences}, nil
}

// Event publishing helpers
func (b *ChatBiz) sendNewMessageEvent(messageId, channelId, userId uint64) {
	event := &proto.NewMessageEvent{
		MessageId: messageId,
		ChannelId: channelId,
		UserId:    userId,
		EventTime: common.NowMicro(),
	}

	msg := &sarama.ProducerMessage{
		Topic: b.newMessageProducerConfig.Topic,
		Key:   sarama.StringEncoder(fmt.Sprintf("%d", channelId)),
		Value: sarama.ByteEncoder(common.ToJsonIgnoreErr(event)),
	}

	partition, offset, err := b.newMessageProducer.SendMessage(msg)
	common.L().Infow(
		"sentNewMessageEvent",
		"messageId", messageId,
		"partition", partition,
		"offset", offset,
		"err", err,
	)
}

func (b *ChatBiz) sendMembershipEvent(channelId, userId uint64, action string) {
	event := &proto.ChannelMembershipEvent{
		ChannelId: channelId,
		UserId:    userId,
		Action:    action,
		EventTime: common.NowMicro(),
	}

	msg := &sarama.ProducerMessage{
		Topic: b.membershipProducerConfig.Topic,
		Key:   sarama.StringEncoder(fmt.Sprintf("%d", channelId)),
		Value: sarama.ByteEncoder(common.ToJsonIgnoreErr(event)),
	}

	partition, offset, err := b.membershipProducer.SendMessage(msg)
	common.L().Infow(
		"sentMembershipEvent",
		"channelId", channelId,
		"userId", userId,
		"action", action,
		"partition", partition,
		"offset", offset,
		"err", err,
	)
}

func (b *ChatBiz) sendPresenceEvent(userId uint64, status, statusText string) {
	event := &proto.UserPresenceEvent{
		UserId:     userId,
		Status:     status,
		StatusText: statusText,
		EventTime:  common.NowMicro(),
	}

	msg := &sarama.ProducerMessage{
		Topic: b.presenceProducerConfig.Topic,
		Key:   sarama.StringEncoder(fmt.Sprintf("%d", userId)),
		Value: sarama.ByteEncoder(common.ToJsonIgnoreErr(event)),
	}

	partition, offset, err := b.presenceProducer.SendMessage(msg)
	common.L().Infow(
		"sentPresenceEvent",
		"userId", userId,
		"status", status,
		"partition", partition,
		"offset", offset,
		"err", err,
	)
}
