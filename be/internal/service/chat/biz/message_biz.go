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

type MessageBiz struct {
	messageDao               *dao.MessageDao
	channelDao               *dao.ChannelDao
	newMessageProducer       sarama.SyncProducer
	newMessageProducerConfig common.ProducerConfig
}

func NewMessageBiz(
	messageDao *dao.MessageDao,
	channelDao *dao.ChannelDao,
	newMessageProducer sarama.SyncProducer,
	newMessageProducerConfig common.ProducerConfig,
) *MessageBiz {
	return &MessageBiz{
		messageDao:               messageDao,
		channelDao:               channelDao,
		newMessageProducer:       newMessageProducer,
		newMessageProducerConfig: newMessageProducerConfig,
	}
}

// Message operations
func (b *MessageBiz) SendMessage(req *proto.SendMessageReq) (*proto.SendMessageResp, error) {
	// Check if user is member of the channel
	isMember, err := b.channelDao.IsUserInChannel(req.ChannelId, req.UserId)
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

	err = b.messageDao.CreateMessage(message)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create message: %v", err)
	}

	// Send new message event
	b.sendNewMessageEvent(message.Id, req.ChannelId, req.UserId)

	// Update last read timestamp for the sender
	b.channelDao.UpdateLastReadAt(req.ChannelId, req.UserId, message.Ctime)

	return &proto.SendMessageResp{MessageId: message.Id}, nil
}

func (b *MessageBiz) GetMessages(req *proto.GetMessagesReq) (*proto.GetMessagesResp, error) {
	// Check if user is member of the channel
	isMember, err := b.channelDao.IsUserInChannel(req.ChannelId, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check membership: %v", err)
	}
	if !isMember {
		return nil, status.Errorf(codes.PermissionDenied, "user is not a member of this channel")
	}

	messages, pagination, err := b.messageDao.GetMessages(req.ChannelId, req.Pagination)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get messages: %v", err)
	}

	return &proto.GetMessagesResp{
		Messages:       messages,
		NextPagination: pagination,
	}, nil
}

// Event publishing helpers
func (b *MessageBiz) sendNewMessageEvent(messageId, channelId, userId uint64) {
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
