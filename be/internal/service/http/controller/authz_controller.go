package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/chientranthien/goldenpay/internal/common"
	"github.com/chientranthien/goldenpay/internal/proto"
)

type (
	AuthzReq struct {
		token string
	}

	AuthzResp struct {
		Code *common.Code `json:"code"`
	}

	AuthzController struct {
		uclient proto.UserServiceClient
	}
)

func (r AuthzReq) toUserServiceLoginReq() *proto.AuthzReq {
	return &proto.AuthzReq{
		Token: "",
	}
}

func NewAuthzController(client proto.UserServiceClient) *AuthzController {
	return &AuthzController{uclient: client}
}

func (c AuthzController) Do(ctx *gin.Context) {
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
		common.L().Errorw("authzErr", "req", req, "err", err)
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
