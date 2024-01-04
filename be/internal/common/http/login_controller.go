package http

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"

	"github.com/chientranthien/goldenpay/internal/common"
	"github.com/chientranthien/goldenpay/internal/proto"
)

type (
	LoginBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	LoginData struct {
		UserId uint64 `json:"user_id"`
	}

	LoginController struct {
		uClient proto.UserServiceClient
		ctx     context.Context
		body    *LoginBody
	}
)

func NewLoginController(client proto.UserServiceClient) *LoginController {
	return &LoginController{uClient: client}
}

func (c LoginController) Do() (common.AnyPtr, common.Code) {
	loginResp, err := c.uClient.Login(c.ctx, c.body.toUserServiceLoginReq())

	code := common.GetCodeFromErr(err)
	if !code.Success() {
		common.L().Errorw("loginErr", "body", c.body, "err", err)
		return nil, code
	}

	SetCookie(c.ctx, Cookie{
		Name:     CookieToken,
		Value:    loginResp.Token,
		MaxAge:   int(3 * 24 * time.Hour.Seconds()),
		Path:     "/",
		Domain:   "",
		Secure:   false,
		HttpOnly: false,
	})
	SetCookie(c.ctx, Cookie{
		CookieUserId,
		fmt.Sprintf("%d", loginResp.UserId),
		int(3 * 24 * time.Hour.Seconds()),
		"/",
		"",
		false,
		false,
	})

	return &LoginData{}, common.CodeSuccess
}

func (c *LoginController) Take(ctx context.Context, req Req) common.Code {
	if b, ok := req.Body.(*LoginBody); ok {
		c.body = b
		c.ctx = ctx
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

func (r LoginBody) toUserServiceLoginReq() *proto.LoginReq {
	return &proto.LoginReq{
		Email:    r.Email,
		Password: r.Password,
	}
}

