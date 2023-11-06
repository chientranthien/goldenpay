package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/chientranthien/goldenpay/internal/common"
	"github.com/chientranthien/goldenpay/internal/proto"
)

type LoginReq proto.LoginReq
type LoginResp proto.LoginResp

type LoginController struct {
	uclient proto.UserServiceClient
}

func NewLoginController(client proto.UserServiceClient) *SignupController {
	return &SignupController{uclient: client}
}

func (c SignupController) Login(ctx *gin.Context) {
	req := &LoginReq{}
	if ctx.BindJSON(req) != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	resp, err := c.uclient.Login(common.Ctx(), (*proto.LoginReq)(req))
	if err != nil {
		log.Printf("failed to signup, err=%v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	ctx.SetCookie("token", resp.Token, 3600, "", "", false, false)
}
