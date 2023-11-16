package controller

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/chientranthien/goldenpay/internal/common"
	"github.com/chientranthien/goldenpay/internal/proto"
)

const (
	TokenCookie = "token"
)

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r LoginReq) toUserServiceLoginReq() *proto.LoginReq {
	return &proto.LoginReq{
		Email:    r.Email,
		Password: r.Password,
	}
}

type LoginResp struct {
	Code *common.Code `json:"code"`
}

type LoginController struct {
	uClient proto.UserServiceClient
}

func NewLoginController(client proto.UserServiceClient) *LoginController {
	return &LoginController{uClient: client}
}

func (c LoginController) Login(ctx *gin.Context) {
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

	serviceResp, err := c.uClient.Login(common.Ctx(), req.toUserServiceLoginReq())

	code := status.Code(err)
	resp := &LoginResp{Code: common.GetCode(int32(code))}
	if resp.Code.IsSuccess() {
		ctx.SetCookie(TokenCookie, serviceResp.Token, int(3*24*time.Hour.Seconds()), "/", "", false, false)
	}

	ctx.JSON(http.StatusOK, resp)

	if err != nil {
		log.Printf("failed to login, err=%v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{})
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
