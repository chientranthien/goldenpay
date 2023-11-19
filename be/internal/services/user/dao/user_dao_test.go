package dao

import (
	"testing"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/stretchr/testify/assert"

	"github.com/chientranthien/goldenpay/internal/proto"
	"github.com/chientranthien/goldenpay/internal/services/user/config"
)

func TestInsert(t *testing.T) {
	testCases := []struct {
		name string
		user *proto.User
	}{
		{
			name: "success",
			user: &proto.User{
				Email:          "test1",
				HashedPassword: "abcd",
				Status:         1,
				Version:        1,
				Ctime:          uint64(time.Now().UnixMilli()),
				Mtime:          uint64(time.Now().UnixMilli()),
			},
		},
	}

	for _, tc := range testCases {
		ltc := tc
		t.Parallel()
		t.Run(ltc.name, func(t *testing.T) {
			db, err := gorm.Open(mysql.Open(config.Get().DB.GetDSN()))
			assert.Nil(t, err)
			dao := NewUserDao(db)
			dao.Insert(ltc.user)
		})
	}
}

func TestUpdate(t *testing.T) {
	testCases := []struct {
		name string
		user *proto.User
	}{
		{
			name: "success",
			user: &proto.User{
				Id:      1,
				Email:   "test1",
				Status:  2,
				Version: 2,
				Ctime:   uint64(time.Now().UnixMilli()),
				Mtime:   uint64(time.Now().UnixMilli()),
			},
		},
	}

	for _, tc := range testCases {
		ltc := tc
		t.Parallel()
		t.Run(ltc.name, func(t *testing.T) {
			db, err := gorm.Open(mysql.Open(config.Get().DB.GetDSN()))
			assert.Nil(t, err)
			dao := NewUserDao(db)
			dao.Update(ltc.user)
		})
	}
}
