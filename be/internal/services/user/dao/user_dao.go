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
	_ = d.getDB().Create(user)
	return nil
}

func (d *UserDao) Update(user *model.User) error {
	_ = d.getDB().Updates(user).Commit()
	return nil
}

func (d *UserDao) GetByEmail(email string) (*model.User, error) {
	u := &model.User{
		Email: email,
	}
	d.getDB().Where(u).First(u)

	return u, nil
}

func (d *UserDao) Get(userID uint64) (*model.User, error) {
	u := &model.User{
		Id: userID,
	}
	d.getDB().First(u)

	return u, nil
}
