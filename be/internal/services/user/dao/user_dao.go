package dao

import (
	"gorm.io/gorm"

	"github.com/chientranthien/goldenpay/internal/services/user/model"
)

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{db: db}
}

func (d UserDao) getDB() *gorm.DB {
	return d.db.Table("user_tab")
}

func (d *UserDao) Insert(user *model.User) error {
	return d.getDB().Create(user).Error
}

func (d *UserDao) Update(user *model.User) error {
	return d.getDB().Updates(user).Commit().Error
}

func (d *UserDao) GetByEmail(email string) (*model.User, error) {
	u := &model.User{
		Email: email,
	}

	if err := d.getDB().Where(u).First(u).Error; err != nil {
		return nil, err
	}

	return u, nil
}

func (d *UserDao) Get(userID uint64) (*model.User, error) {
	u := &model.User{
		Id: userID,
	}

	if err := d.getDB().First(u).Error; err != nil {
		return nil, err
	}

	return u, nil
}
