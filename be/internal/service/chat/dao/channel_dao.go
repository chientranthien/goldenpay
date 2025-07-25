package dao

import (
	"gorm.io/gorm"

	"github.com/chientranthien/goldenpay/internal/common"
	"github.com/chientranthien/goldenpay/internal/proto"
)

type ChannelDao struct {
	db *gorm.DB
}

func NewChannelDao(db *gorm.DB) *ChannelDao {
	return &ChannelDao{db: db}
}

func (d *ChannelDao) getDB() *gorm.DB {
	return d.db.Table("channel_tab")
}

func (d *ChannelDao) getChannelMemberDB() *gorm.DB {
	return d.db.Table("channel_member_tab")
}

// Channel operations
func (d *ChannelDao) InsertChannel(channel *proto.Channel) error {
	result := d.getDB().Create(channel)
	return result.Error
}

func (d *ChannelDao) Get(channelId uint64) (*proto.Channel, error) {
	var channel proto.Channel

	result := d.getDB().First(&channel, channelId)
	if result.Error != nil {
		return nil, result.Error
	}

	return &channel, nil
}

func (d *ChannelDao) GetChannelMembersByChannelId(channelId uint64) ([]*proto.ChannelMember, error) {
	var members []*proto.ChannelMember
	result := d.getChannelMemberDB().Where("channel_id = ?", channelId).Find(&members)
	return members, result.Error
}

func (d *ChannelDao) List(userId uint64, pagination *proto.Pagination) ([]*proto.Channel, *proto.Pagination, error) {
	var channels []*proto.Channel
	var total int64

	// Get channels where user is a member
	subQuery := d.getChannelMemberDB().Select("channel_id").Where("user_id = ?", userId)

	query := d.getDB().Where("id IN (?)", subQuery)

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

	result := query.Limit(limit + 1).Find(&channels)

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

func (d *ChannelDao) Update(channel *proto.Channel) error {
	return d.getDB().Save(channel).Error
}

// Channel membership operations
func (d *ChannelDao) AddChannelMember(member *proto.ChannelMember) error {
	return d.getChannelMemberDB().Create(member).Error
}

func (d *ChannelDao) IsUserInChannel(channelId, userId uint64) (bool, error) {
	var count int64
	result := d.getChannelMemberDB().
		Where("channel_id = ? AND user_id = ?", channelId, userId).
		Count(&count)
	return count > 0, result.Error
}

func (d *ChannelDao) UpdateLastReadAt(channelId, userId uint64, timestamp uint64) error {
	return d.getChannelMemberDB().
		Where("channel_id = ? AND user_id = ?", channelId, userId).
		Update("last_read_at", timestamp).Error
}

// Direct message operations
func (d *ChannelDao) FindDirectMessageChannel(user1Id, user2Id uint64) (*proto.Channel, error) {
	var channel proto.Channel

	// Find channel where both users are members and channel type is "direct"
	subQuery1 := d.getChannelMemberDB().Select("channel_id").Where("user_id = ?", user1Id)
	subQuery2 := d.getChannelMemberDB().Select("channel_id").Where("user_id = ?", user2Id)

	result := d.getDB().
		Where("type = ? AND id IN (?) AND id IN (?)", "direct", subQuery1, subQuery2).
		First(&channel)

	if result.Error != nil {
		return nil, result.Error
	}

	return &channel, nil
}

func (d *ChannelDao) GetDirectMessageChannels(
	userId uint64,
	pagination *proto.Pagination,
) ([]*proto.Channel, *proto.Pagination, error) {
	var channels []*proto.Channel

	subQuery := d.getChannelMemberDB().Select("channel_id").Where("user_id = ?", userId)
	query := d.getDB().Where("type = ? AND id IN (?)", "direct", subQuery)

	// Apply pagination
	if pagination != nil && pagination.Val > 0 {
		query = query.Where("id < ?", pagination.Val)
	}

	limit := int(pagination.GetLimit())
	if limit == 0 {
		limit = 20 // default
	}

	result := query.Order("mtime DESC").Limit(limit + 1).Find(&channels)

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

func (d *ChannelDao) Begin() (*ChannelDao, error) {
	trx := d.getDB().Begin()
	if err := trx.Error; err != nil {
		common.L().Errorw("beginErr", "err", err)
		return nil, err
	}
	return &ChannelDao{db: trx}, nil
}

func (d *ChannelDao) CommitOrRollback(success bool) {
	if success {
		d.db.Commit()
	} else {
		d.db.Rollback()
	}
}
