package biz

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/chientranthien/goldenpay/internal/common"
	"github.com/chientranthien/goldenpay/internal/proto"
	"github.com/chientranthien/goldenpay/internal/services/user/config"
	"github.com/chientranthien/goldenpay/internal/services/user/dao"
	"github.com/chientranthien/goldenpay/internal/services/user/model"
)

type UserBiz struct {
	jwtConfig config.JWTConfig
	dao       *dao.UserDao
}

func NewUserBiz(
	jwtConfig config.JWTConfig,
	dao *dao.UserDao,
) *UserBiz {
	return &UserBiz{jwtConfig: jwtConfig, dao: dao}
}

func (b UserBiz) Signup(req *proto.SignupReq) (*proto.SignupResp, error) {
	user := &model.User{
		Email:          req.Email,
		HashedPassword: b.HashPassword(req.Password),
		Name:           req.Name,
		Status:         model.StatusActive,
		Version:        model.VersionFirst,
		Ctime:          common.NowMillis(),
		Mtime:          common.NowMillis(),
	}

	err := b.dao.Insert(user)
	if err != nil {
		return nil, err
	}

	return &proto.SignupResp{}, nil
}

func (b UserBiz) Get(ctx context.Context, id uint64) (*model.User, error) {
	user, err := b.dao.Get(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
func (b UserBiz) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	user, err := b.dao.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (b UserBiz) HashPassword(password string) string {
	sum256 := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", sum256)
}

func (b UserBiz) GenerateToken(user *model.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = user.Email
	claims["uid"] = user.Id
	claims["exp"] = (time.Duration(b.jwtConfig.DurationInMin) * time.Minute).Seconds()

	tokenString, err := token.SignedString([]byte(b.jwtConfig.Secret))
	if err != nil {
		fmt.Println("failed to create token, err=%v", err)
		return "", err
	}

	return tokenString, nil
}

func (b UserBiz) Login(ctx context.Context, username string) (interface{}, interface{}) {
	return nil, nil
}
