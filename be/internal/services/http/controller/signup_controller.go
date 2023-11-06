package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/chientranthien/goldenpay/internal/common"
	"github.com/chientranthien/goldenpay/internal/proto"
)

type SignupReq proto.SignupReq
type SignupResp proto.SignupResp

type SignupController struct {
	uclient proto.UserServiceClient
}

func NewSignupController(client proto.UserServiceClient) *SignupController {
	return &SignupController{uclient: client}
}

func (c SignupController) Signup(ctx *gin.Context) {
	req := &SignupReq{}
	if ctx.BindJSON(req) != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	_, err := c.uclient.Signup(common.Ctx(), (*proto.SignupReq)(req))
	if err != nil {
		log.Printf("failed to signup, err=%v",err)
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
}
