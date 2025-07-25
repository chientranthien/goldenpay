package dao

import (
	"gorm.io/gorm"

	"github.com/chientranthien/goldenpay/internal/common"
	"github.com/chientranthien/goldenpay/internal/proto"
)

type MessageDao struct {
	db *gorm.DB
}

func NewMessageDao(db *gorm.DB) *MessageDao {
	return &MessageDao{db: db}
}

func (d *MessageDao) getDB() *gorm.DB {
	return d.db.Table("message_tab")
}

// Message operations
func (d *MessageDao) CreateMessage(message *proto.Message) error {
	return d.getDB().Create(message).Error
}

func (d *MessageDao) Get(messageId uint64) (*proto.Message, error) {
	var message proto.Message
	result := d.getDB().First(&message, messageId)
	if result.Error != nil {
		return nil, result.Error
	}
	return &message, nil
}

func (d *MessageDao) GetMessages(channelId uint64, pagination *proto.Pagination) ([]*proto.Message, *proto.Pagination, error) {
	var messages []*proto.Message

	query := d.getDB().Model(&proto.Message{}).Where("channel_id = ? AND status != ?", channelId, common.MessageStatusDeleted)

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

func (d *MessageDao) Begin() (*MessageDao, error) {
	trx := d.getDB().Begin()
	if err := trx.Error; err != nil {
		common.L().Errorw("beginErr", "err", err)
		return nil, err
	}
	return &MessageDao{db: trx}, nil
}

func (d *MessageDao) CommitOrRollback(success bool) {
	if success {
		d.db.Commit()
	} else {
		d.db.Rollback()
	}
}
