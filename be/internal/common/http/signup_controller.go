package http

import (
	"context"

	"google.golang.org/grpc/codes"

	"github.com/chientranthien/goldenpay/internal/common"
	"github.com/chientranthien/goldenpay/internal/proto"
)

type (
	SignupBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}

	SignupData struct {
	}

	SignupController struct {
		uclient proto.UserServiceClient
		ctx     context.Context
		body    *SignupBody
	}
)

func (r SignupBody) toUserServiceSignupReq() *proto.SignupReq {
	return &proto.SignupReq{
		Email:    r.Email,
		Password: r.Password,
		Name:     r.Name,
	}
}

func NewSignupController(client proto.UserServiceClient) *SignupController {
	return &SignupController{uclient: client}
}

func (c SignupController) Do() (common.AnyPtr, common.Code) {
	_, err := c.uclient.Signup(c.ctx, c.body.toUserServiceSignupReq())

	code := common.GetCodeFromErr(err)
	if !code.Success() {
		common.L().Errorw("signupErr", "body", c.body, "err", err)
		return nil, code
	}

	return &SignupData{}, code
}

func (c *SignupController) Take(ctx context.Context, req Req) common.Code {
	if b, ok := req.Body.(*SignupBody); ok {
		c.ctx = ctx
		c.body = b
	} else {
		return common.CodeBody
	}

	if common.ValidateEmail(c.body.Email) != nil {
		return common.Code{
			Id:  int32(codes.InvalidArgument),
			Msg: "invalid email",
		}
	}

	return common.CodeSuccess
}
