package controller

import (
	"github.com/chientranthien/goldenpay/internal/proto"
	"github.com/chientranthien/goldenpay/internal/service/user/biz"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CreateContactIfNotExistController struct {
	biz *biz.UserBiz
}

func NewCreateContactIfNotExistController(biz *biz.UserBiz) *CreateContactIfNotExistController {
	return &CreateContactIfNotExistController{biz: biz}
}

func (c CreateContactIfNotExistController) Do(req *proto.CreateContactIfNotExistReq) (
	*proto.CreateContactIfNotExistResp, error) {
	if err := c.validate(req); err != nil {
		return nil, err
	}

	return c.biz.CreateContactIfNotExist(req)
}

func (c CreateContactIfNotExistController) validate(req *proto.CreateContactIfNotExistReq) error {
	if req.UserId == 0 || req.ContactUserId == 0 {
		return status.Error(codes.InvalidArgument, "invalid user or contact user")
	}

	return nil
}
