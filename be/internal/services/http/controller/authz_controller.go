package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/chientranthien/goldenpay/internal/common"
	"github.com/chientranthien/goldenpay/internal/proto"
)

type AuthzReq struct {
	token string
}

func (r AuthzReq) toUserServiceLoginReq() *proto.AuthzReq {
	return &proto.AuthzReq{
		Token: "",
	}
}

type AuthzResp struct {
	Code *common.Code `json:"code"`
}

type AuthzController struct {
	uclient proto.UserServiceClient
}

func NewAuthzController(client proto.UserServiceClient) *AuthzController {
	return &AuthzController{uclient: client}
}

func (c AuthzController) Authz(ctx *gin.Context) {
	token, _ := ctx.Cookie(TokenCookie)
	req := &AuthzReq{
		token: token,
	}

	if code := c.validate(req); !code.IsSuccess() {
		resp := &LoginResp{
			Code: code,
		}
		ctx.JSON(http.StatusOK, resp)
		return
	}

	_, err := c.uclient.Authz(common.Ctx(), &proto.AuthzReq{Token: token})

	code := status.Code(err)
	resp := &LoginResp{Code: common.GetCode(int32(code))}
	ctx.JSON(http.StatusOK, resp)

	if err != nil {
		log.Printf("failed to login, err=%v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
}

func (c AuthzController) validate(req *AuthzReq) *common.Code {
	if req.token == "" {
		return &common.Code{
			Id:  int32(codes.Unauthenticated),
			Msg: "invalid token",
		}
	}

	return common.CodeSuccess
}
