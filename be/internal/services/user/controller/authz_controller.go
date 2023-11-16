package controller

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/chientranthien/goldenpay/internal/proto"
	"github.com/chientranthien/goldenpay/internal/services/user/biz"
)

type AuthzController struct {
	biz *biz.UserBiz
}

func NewAuthzController(biz *biz.UserBiz) *AuthzController {
	return &AuthzController{biz: biz}
}

func (c AuthzController) Authz(ctx context.Context, req *proto.AuthzReq) (*proto.AuthzResp, error) {
	_, err := c.biz.ParseToken(req.Token)
	if err != nil {
		return nil, err
	}

	return &proto.AuthzResp{}, nil
}

func (c AuthzController) validate(req *proto.AuthzReq) error {
	if len(req.Token) == 0 || len(req.Resource) == 0 {
		return status.New(codes.InvalidArgument, "invalid token or resource").Err()
	}

	return nil
}
