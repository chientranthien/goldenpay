package biz

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/chientranthien/goldenpay/internal/common"
	"github.com/chientranthien/goldenpay/internal/proto"
	"github.com/chientranthien/goldenpay/internal/service/chat/dao"
)

type ChannelBiz struct {
	channelDao *dao.ChannelDao
}

func NewChannelBiz(channelDao *dao.ChannelDao) *ChannelBiz {
	return &ChannelBiz{
		channelDao: channelDao,
	}
}

func (b *ChannelBiz) GetChannel(req *proto.GetChannelReq) (*proto.GetChannelResp, error) {
	// Check if user is member of the channel
	isMember, err := b.channelDao.IsUserInChannel(req.ChannelId, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check membership: %v", err)
	}
	if !isMember {
		return nil, status.Errorf(codes.PermissionDenied, "user is not a member of this channel")
	}

	channel, err := b.channelDao.Get(req.ChannelId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "channel not found: %v", err)
	}

	members, err := b.channelDao.GetChannelMembersByChannelId(req.ChannelId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get channel members: %v", err)
	}

	return &proto.GetChannelResp{Channel: channel, Members: members}, nil
}

func (b *ChannelBiz) CreateDirectMessageChannel(req *proto.CreateDirectMessageChannelReq) (*proto.CreateDirectMessageChannelResp, error) {
	// Check if direct message channel already exists
	existingChannel, err := b.channelDao.FindDirectMessageChannel(req.FromUser, req.ToUser)
	if err == nil {
		// Channel already exists
		return &proto.CreateDirectMessageChannelResp{ChannelId: existingChannel.Id}, nil
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

	txChannelDao, err := b.channelDao.Begin()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to start transaction: %v", err)
	}
	defer txChannelDao.CommitOrRollback(err == nil)

	err = txChannelDao.InsertChannel(channel)
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
		err = txChannelDao.AddChannelMember(member)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to add member: %v", err)
		}
	}

	return &proto.CreateDirectMessageChannelResp{ChannelId: channel.Id}, nil
}

func (b *ChannelBiz) GetDirectMessages(req *proto.GetDirectMessagesReq) (*proto.GetDirectMessagesResp, error) {
	channels, pagination, err := b.channelDao.GetDirectMessageChannels(req.UserId, req.Pagination)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get direct messages: %v", err)
	}

	return &proto.GetDirectMessagesResp{
		Channels:       channels,
		NextPagination: pagination,
	}, nil
}
