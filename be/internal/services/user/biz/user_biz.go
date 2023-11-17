package biz

import (
	"crypto/sha256"
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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

func (b UserBiz) Get(id uint64) (*model.User, error) {
	user, err := b.dao.Get(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
func (b UserBiz) GetByEmail(email string) (*model.User, error) {
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
	claims := &jwt.StandardClaims{
		Audience:  fmt.Sprintf("%d", user.Id),
		ExpiresAt: time.Now().Add(time.Duration(b.jwtConfig.DurationInMin) * time.Minute).Unix(),
		IssuedAt:  time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(b.jwtConfig.Secret))
	if err != nil {
		log.Printf("failed to create token, err=%v", err)
		return "", err
	}

	return tokenString, nil
}

type Claim struct {
	jwt.StandardClaims
}

func (b UserBiz) ParseToken(tokenStr string) (*jwt.StandardClaims, error) {
	claims := &jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(b.jwtConfig.Secret), nil
	})

	if err != nil {
		log.Printf("failed to parse token, tokenStr=%v, err=%v", tokenStr, err)
		return nil, err
	}

	if !token.Valid {
		c := codes.Unauthenticated
		return nil, status.Error(c, c.String())
	}

	expiredAt := time.Unix(claims.ExpiresAt, 0)
	if time.Now().After(expiredAt) {
		c := codes.DeadlineExceeded
		return nil, status.Error(c, c.String())
	}

	return claims, nil
}
