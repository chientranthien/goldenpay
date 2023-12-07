package biz

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/chientranthien/goldenpay/internal/common"
	"github.com/chientranthien/goldenpay/internal/proto"
	"github.com/chientranthien/goldenpay/internal/service/user/config"
	"github.com/chientranthien/goldenpay/internal/service/user/dao"
)

type UserBiz struct {
	jwtConfig      config.JWTConfig
	dao            *dao.UserDao
	producer       sarama.SyncProducer
	producerConfig common.ProducerConfig
}

func NewUserBiz(
	jwtConfig config.JWTConfig,
	dao *dao.UserDao,
	producer sarama.SyncProducer,
	producerConfig common.ProducerConfig,
) *UserBiz {
	return &UserBiz{
		jwtConfig:      jwtConfig,
		dao:            dao,
		producer:       producer,
		producerConfig: producerConfig,
	}
}

func (b UserBiz) Signup(req *proto.SignupReq) (*proto.SignupResp, error) {
	user := &proto.User{
		Email:          req.Email,
		HashedPassword: b.HashPassword(req.Password),
		Name:           req.Name,
		Status:         common.UserStatusActive,
		Version:        common.FirstVersion,
		Ctime:          common.NowMillis(),
		Mtime:          common.NowMillis(),
	}

	err := b.dao.Insert(user)
	if err != nil {
		return nil, err
	}

	e := &proto.NewUserEvent{
		UserId:    user.Id,
		EventTime: common.NowMicro(),
	}

	msg := &sarama.ProducerMessage{
		Topic: b.producerConfig.Topic,
		Value: sarama.ByteEncoder(common.ToJsonIgnoreErr(e)),
	}

	partition, offset, produceErr := b.producer.SendMessage(msg)
	common.L().Infow(
		"sentNewUserEvent",
		"partition", partition,
		"offset", offset,
		"err", produceErr,
	)

	return &proto.SignupResp{}, nil
}
func (b UserBiz) Login(req *proto.LoginReq) (*proto.LoginResp, error) {
	getResp, err := b.GetByEmail(&proto.GetByEmailReq{Email: req.Email})
	if err != nil {
		return nil, status.New(codes.Internal, "unable to get user").Err()
	}
	user := getResp.User

	if user.Id == 0 {
		return nil, status.New(codes.NotFound, "not found").Err()
	}

	if user.HashedPassword != b.HashPassword(req.Password) {
		return nil, status.New(codes.InvalidArgument, "incorrect password").Err()
	}

	token, err := b.GenerateToken(user)
	if err != nil {
		return nil, status.New(codes.Internal, "unable to generate token").Err()
	}

	return &proto.LoginResp{Token: token, UserId: getResp.User.Id}, nil
}

func (b UserBiz) Get(id uint64) (*proto.User, error) {
	user, err := b.dao.Get(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (b UserBiz) GetByEmail(req *proto.GetByEmailReq) (*proto.GetByEmailResp, error) {
	user, err := b.dao.GetByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	return &proto.GetByEmailResp{User: user}, nil
}

func (b UserBiz) HashPassword(password string) string {
	sum256 := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", sum256)
}

func (b UserBiz) GenerateToken(user *proto.User) (string, error) {
	claims := &jwt.StandardClaims{
		Audience:  fmt.Sprintf("%d", user.Id),
		ExpiresAt: time.Now().Add(time.Duration(b.jwtConfig.DurationInMin) * time.Minute).Unix(),
		IssuedAt:  time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(b.jwtConfig.Secret))
	if err != nil {
		common.L().Errorw(
			"generateTokenErr",
			"userID", user.Id,
			"jwtConfig", b.jwtConfig,
			"err", err,
		)
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
		common.L().Errorw(
			"parseTokenErr",
			"tokenStr", tokenStr,
			"err", err,
		)
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

func (b UserBiz) GetBatch(ids []uint64) (*proto.GetBatchResp, error) {
	usersFromDB, err := b.dao.GetBatch(ids)
	if err != nil {
		return nil, status.New(codes.Internal, "unable to query DB").Err()
	}

	userMap := make(map[uint64]*proto.User)
	for _, user := range usersFromDB {
		userMap[user.Id] = user
	}
	users := make([]*proto.User, len(ids))
	for i, id := range ids {
		users[i] = userMap[id]
	}

	return &proto.GetBatchResp{Users: users}, nil
}
