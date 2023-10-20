package dao

import "github.com/chientranthien/goldenpay/internal/proto"

type UserDao struct {
}

func (u UserDao) Insert(user proto.User) error {

}

func (u UserDao) Update(user proto.User) error {

}

func (u UserDao) GetByUsername(username string) (proto.User, error) {

}

func (u UserDao) GetByUserId(userId uint64) (proto.User, error) {

}
