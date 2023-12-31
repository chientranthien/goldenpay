package http

import (
	"context"

	"github.com/chientranthien/goldenpay/internal/common"
	"github.com/chientranthien/goldenpay/internal/proto"
)

type (
	AuthzBody struct {
	}

	AuthzData struct {
	}

	AuthzController struct {
		uclient proto.UserServiceClient
	}
)

func (r AuthzBody) toUserServiceLoginReq() *proto.AuthzReq {
	return &proto.AuthzReq{
		Token: "",
	}
}

func NewAuthzController(client proto.UserServiceClient) *AuthzController {
	return &AuthzController{uclient: client}
}

func (c AuthzController) Do() (common.AnyPtr, common.Code) {
	return &AuthzBody{}, common.CodeSuccess
}

func (c *AuthzController) Take(_ context.Context, _ Req) common.Code {
	return common.CodeSuccess
}
