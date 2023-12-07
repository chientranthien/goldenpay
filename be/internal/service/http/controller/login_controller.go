package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"

	"github.com/chientranthien/goldenpay/internal/common"
	"github.com/chientranthien/goldenpay/internal/proto"
)

const (
	TokenCookie  = "token"
	UserIdCookie = "uid"
)

type Controller interface {
	Validate() *common.Code
	Do(ctx *gin.Context)
	GetReq() any
	GetResp() any
}

type Wrapper struct {
	c Controller
}

func (c Wrapper) Do(ctx *gin.Context) {
	if ctx.BindJSON(c.c.GetReq()) != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	if code := c.c.Validate(); !code.IsSuccess() {
		resp := &LoginResp{
			Code: code,
		}
		ctx.JSON(http.StatusOK, resp)
		return
	}

	c.c.Do(ctx)
}

type (
	LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	LoginData struct {
		UserId uint64 `json:"user_id"`
	}

	LoginResp struct {
		Code *common.Code `json:"code"`
		Data LoginData    `json:"data"`
	}

	LoginController struct {
		uClient proto.UserServiceClient
	}
)

func (r LoginReq) toUserServiceLoginReq() *proto.LoginReq {
	return &proto.LoginReq{
		Email:    r.Email,
		Password: r.Password,
	}
}

func NewLoginController(client proto.UserServiceClient) *LoginController {
	return &LoginController{uClient: client}
}

func (c LoginController) Do(ctx *gin.Context) {
	req := &LoginReq{}
	if ctx.BindJSON(req) != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	if code := c.validate(req); !code.IsSuccess() {
		resp := &LoginResp{
			Code: code,
		}
		ctx.JSON(http.StatusOK, resp)
		return
	}

	token, _ := ctx.Cookie(TokenCookie)
	if token != "" {
		_, err := c.uClient.Authz(common.Ctx(), &proto.AuthzReq{Token: token})
		if err == nil {
			ctx.JSON(http.StatusOK, &LoginResp{Code: common.CodeSuccess})
			return
		}
	}

	loginResp, err := c.uClient.Login(common.Ctx(), req.toUserServiceLoginReq())

	resp := &LoginResp{Code: common.GetCodeFromErr(err)}
	if resp.Code.IsSuccess() {
		ctx.SetCookie(TokenCookie, loginResp.Token, int(3*24*time.Hour.Seconds()), "/", "", false, false)
		ctx.SetCookie(
			UserIdCookie,
			fmt.Sprintf("%d", loginResp.UserId),
			int(3*24*time.Hour.Seconds()),
			"/",
			"",
			false,
			false,
		)
	}

	ctx.JSON(http.StatusOK, resp)

	if err != nil {
		common.L().Errorw("loginErr", "req", req, "err", err)
		return
	}
}

func (c LoginController) validate(req *LoginReq) *common.Code {
	if common.ValidateEmail(req.Email) != nil {
		return &common.Code{
			Id:  int32(codes.InvalidArgument),
			Msg: "invalid email",
		}
	}

	return common.CodeSuccess
}
