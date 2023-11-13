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
	if code := c.validate(req); code.Id != common.CodeSuccess.Id {
		resp := &SignupResp{
			Code: code,
		}
		ctx.JSON(http.StatusOK, resp)
		return
	}

	_, err := c.uclient.Signup(common.Ctx(), (*proto.SignupReq)(req))

	code := status.Code(err)
	resp := &SignupResp{
		Code: common.GetCode(int32(code)),
	}
	ctx.JSON(http.StatusOK, resp)
	if err != nil {
		log.Printf("failed to signup, req=%v, err=%v", err)
		return
	}
}
func (c SignupController) validate(req *SignupReq) *proto.Code {
	if common.ValidateEmail(req.Email) != nil {
		return &proto.Code{
			Id:  int32(codes.InvalidArgument),
			Msg: "invalid email",
		}
	}

	return common.CodeSuccess
}
