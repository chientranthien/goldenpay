package dao

import (
	"gorm.io/gorm"

	"github.com/chientranthien/goldenpay/internal/common"
	"github.com/chientranthien/goldenpay/internal/proto"
)

type ChatDao struct {
	db *gorm.DB
}

func NewChatDao(db *gorm.DB) *ChatDao {
	// Auto-migrate tables
	db.AutoMigrate(&proto.Channel{}, &proto.Message{}, &proto.ChannelMember{}, &proto.UserPresence{})

	return &ChatDao{db: db}
}

func (d *ChatDao) Begin() (*ChatDao, error) {
	tx := d.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &ChatDao{db: tx}, nil
}

func (d *ChatDao) CommitOrRollback(success bool) {
	if success {
		d.db.Commit()
	} else {
		d.db.Rollback()
	}
}

// Channel operations
func (d *ChatDao) CreateChannel(channel *proto.Channel) error {
	result := d.db.Create(channel)
	return result.Error
}

func (d *ChatDao) GetChannel(channelId uint64) (*proto.Channel, error) {
	var channel proto.Channel
	result := d.db.Preload("Members").First(&channel, channelId)
	if result.Error != nil {
		return nil, result.Error
	}
	return &channel, nil
}

func (d *ChatDao) ListChannelsForUser(userId uint64, pagination *proto.Pagination) ([]*proto.Channel, *proto.Pagination, error) {
	var channels []*proto.Channel
	var total int64

	// Get channels where user is a member
	subQuery := d.db.Model(&proto.ChannelMember{}).Select("channel_id").Where("user_id = ?", userId)

	query := d.db.Model(&proto.Channel{}).Where("id IN (?)", subQuery)

	// Count total
	query.Count(&total)

	// Apply pagination
	if pagination != nil && pagination.Val > 0 {
		query = query.Where("id < ?", pagination.Val)
	}

	limit := int(pagination.GetLimit())
	if limit == 0 {
		limit = 20 // default
	}

	result := query.Limit(limit + 1).Preload("Members").Find(&channels)

	hasMore := len(channels) > limit
	if hasMore {
		channels = channels[:limit]
	}

	nextPagination := &proto.Pagination{
		Limit:   pagination.GetLimit(),
		HasMore: hasMore,
	}

	if hasMore && len(channels) > 0 {
		nextPagination.Val = int64(channels[len(channels)-1].Id)
	}

	return channels, nextPagination, result.Error
}

func (d *ChatDao) UpdateChannel(channel *proto.Channel) error {
	return d.db.Save(channel).Error
}

// Channel membership operations
func (d *ChatDao) AddChannelMember(member *proto.ChannelMember) error {
	return d.db.Create(member).Error
}

func (d *ChatDao) RemoveChannelMember(channelId, userId uint64) error {
	return d.db.Where("channel_id = ? AND user_id = ?", channelId, userId).Delete(&proto.ChannelMember{}).Error
}

func (d *ChatDao) GetChannelMember(channelId, userId uint64) (*proto.ChannelMember, error) {
	var member proto.ChannelMember
	result := d.db.Where("channel_id = ? AND user_id = ?", channelId, userId).First(&member)
	if result.Error != nil {
		return nil, result.Error
	}
	return &member, nil
}

func (d *ChatDao) IsUserInChannel(channelId, userId uint64) (bool, error) {
	var count int64
	result := d.db.Model(&proto.ChannelMember{}).Where("channel_id = ? AND user_id = ?", channelId, userId).Count(&count)
	return count > 0, result.Error
}

// Message operations
func (d *ChatDao) CreateMessage(message *proto.Message) error {
	return d.db.Create(message).Error
}

func (d *ChatDao) GetMessage(messageId uint64) (*proto.Message, error) {
	var message proto.Message
	result := d.db.First(&message, messageId)
	if result.Error != nil {
		return nil, result.Error
	}
	return &message, nil
}

func (d *ChatDao) GetMessages(channelId uint64, pagination *proto.Pagination) ([]*proto.Message, *proto.Pagination, error) {
	var messages []*proto.Message

	query := d.db.Model(&proto.Message{}).Where("channel_id = ? AND status != ?", channelId, common.MessageStatusDeleted)

	// Apply pagination
	if pagination != nil && pagination.Val > 0 {
		query = query.Where("id < ?", pagination.Val)
	}

	limit := int(pagination.GetLimit())
	if limit == 0 {
		limit = 50 // default
	}

	result := query.Order("id DESC").Limit(limit + 1).Find(&messages)

	hasMore := len(messages) > limit
	if hasMore {
		messages = messages[:limit]
	}

	nextPagination := &proto.Pagination{
		Limit:   pagination.GetLimit(),
		HasMore: hasMore,
	}

	if hasMore && len(messages) > 0 {
		nextPagination.Val = int64(messages[len(messages)-1].Id)
	}

	return messages, nextPagination, result.Error
}

func (d *ChatDao) UpdateMessage(message *proto.Message) error {
	return d.db.Save(message).Error
}

func (d *ChatDao) DeleteMessage(messageId, userId uint64) error {
	return d.db.Model(&proto.Message{}).Where("id = ? AND user_id = ?", messageId, userId).Update("status", common.MessageStatusDeleted).Error
}

// Direct message operations
func (d *ChatDao) FindDirectMessageChannel(user1Id, user2Id uint64) (*proto.Channel, error) {
	var channel proto.Channel

	// Find channel where both users are members and channel type is "direct"
	subQuery1 := d.db.Model(&proto.ChannelMember{}).Select("channel_id").Where("user_id = ?", user1Id)
	subQuery2 := d.db.Model(&proto.ChannelMember{}).Select("channel_id").Where("user_id = ?", user2Id)

	result := d.db.Where("type = ? AND id IN (?) AND id IN (?)", "direct", subQuery1, subQuery2).First(&channel)

	if result.Error != nil {
		return nil, result.Error
	}

	return &channel, nil
}

func (d *ChatDao) GetDirectMessageChannels(userId uint64, pagination *proto.Pagination) ([]*proto.Channel, *proto.Pagination, error) {
	var channels []*proto.Channel

	subQuery := d.db.Model(&proto.ChannelMember{}).Select("channel_id").Where("user_id = ?", userId)
	query := d.db.Model(&proto.Channel{}).Where("type = ? AND id IN (?)", "direct", subQuery)

	// Apply pagination
	if pagination != nil && pagination.Val > 0 {
		query = query.Where("id < ?", pagination.Val)
	}

	limit := int(pagination.GetLimit())
	if limit == 0 {
		limit = 20 // default
	}

	result := query.Order("mtime DESC").Limit(limit + 1).Preload("Members").Find(&channels)

	hasMore := len(channels) > limit
	if hasMore {
		channels = channels[:limit]
	}

	nextPagination := &proto.Pagination{
		Limit:   pagination.GetLimit(),
		HasMore: hasMore,
	}

	if hasMore && len(channels) > 0 {
		nextPagination.Val = int64(channels[len(channels)-1].Id)
	}

	return channels, nextPagination, result.Error
}

// Presence operations
func (d *ChatDao) UpsertUserPresence(presence *proto.UserPresence) error {
	return d.db.Save(presence).Error
}

func (d *ChatDao) GetUserPresence(userIds []uint64) ([]*proto.UserPresence, error) {
	var presences []*proto.UserPresence
	result := d.db.Where("user_id IN ?", userIds).Find(&presences)
	return presences, result.Error
}

func (d *ChatDao) UpdateLastReadAt(channelId, userId uint64, timestamp uint64) error {
	return d.db.Model(&proto.ChannelMember{}).Where("channel_id = ? AND user_id = ?", channelId, userId).Update("last_read_at", timestamp).Error
}
