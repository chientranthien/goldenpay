package controller

import (
	"context"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/chientranthien/goldenpay/internal/proto"
	"github.com/chientranthien/goldenpay/internal/service/user/biz"
)

type (
	AuthzController struct {
		biz *biz.UserBiz
	}
)

func NewAuthzController(biz *biz.UserBiz) *AuthzController {
	return &AuthzController{biz: biz}
}

func (c AuthzController) Do(ctx context.Context, req *proto.AuthzReq) (*proto.AuthzResp, error) {
	claims, err := c.biz.ParseToken(req.Token)
	if err != nil {
		return nil, err
	}

	uid, err := strconv.ParseUint(claims.Audience, 10, 64)
	if err != nil {
		return nil, status.New(codes.InvalidArgument, "invalid token").Err()
	}

	return &proto.AuthzResp{Metadata: &proto.AuthzMetadata{UserId: uid}}, nil
}

func (c AuthzController) validate(req *proto.AuthzReq) error {
	if len(req.Token) == 0 || len(req.Resource) == 0 {
		return status.New(codes.InvalidArgument, "invalid token or resource").Err()
	}

	return nil
}
