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
	SignupReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}

	SignupResp struct {
		Code *common.Code `json:"code"`
	}

	SignupController struct {
		uclient proto.UserServiceClient
	}
)

func (r SignupReq) toUserServiceSignupReq() *proto.SignupReq {
	return &proto.SignupReq{
		Email:    r.Email,
		Password: r.Password,
		Name:     r.Name,
	}
}

func NewSignupController(client proto.UserServiceClient) *SignupController {
	return &SignupController{uclient: client}
}

func (c SignupController) Do(ctx *gin.Context) {
	req := &SignupReq{}
	if ctx.BindJSON(req) != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	if code := c.validate(req); !code.IsSuccess() {
		resp := &SignupResp{
			Code: code,
		}
		ctx.JSON(http.StatusOK, resp)
		return
	}

	_, err := c.uclient.Signup(common.Ctx(), req.toUserServiceSignupReq())

	code := status.Code(err)
	resp := &SignupResp{
		Code: common.GetCode(int32(code)),
	}

	ctx.JSON(http.StatusOK, resp)
	if err != nil {
		common.L().Errorw("signupErr", "req", req, "err", err)
		return
	}
}

func (c SignupController) validate(req *SignupReq) *common.Code {
	if common.ValidateEmail(req.Email) != nil {
		return &common.Code{
			Id:  int32(codes.InvalidArgument),
			Msg: "invalid email",
		}
	}

	return common.CodeSuccess
}
