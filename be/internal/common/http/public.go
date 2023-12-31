package http

import (
	"bytes"
	"context"
	"io"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/chientranthien/goldenpay/internal/common"
	"github.com/chientranthien/goldenpay/internal/proto"
)

type (
	Empty interface {
	}

	CommonResp struct {
		Code common.Code `json:"code"`
		Data any         `json:"data"`
	}

	Metadata struct {
		UserId uint64
		Email string
	}

	Req struct {
		Metadata Metadata
		Body     any
	}

	Ctl interface {
		Take(ctx context.Context, req Req) common.Code
		Do() (data common.AnyPtr, code common.Code)
	}

	NewCtlFn func() Ctl

	endpointInfo struct {
		method      string
		ep          string
		newCtlFn    NewCtlFn
		req         any
		resp        any
		preReqHooks []Hook
	}
)

var (
	handler *Server
)

const (
	MethodGet  = "GET"
	MethodPost = "POST"
	MethodPut  = "PUT"

	CookieToken  = "token"
	CookieUserId = "uid"
)

func Init(conf common.ServiceConfig, uClient proto.UserServiceClient) *gin.Engine {
	server := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000", "https://goldenpay.chientran.info"}
	corsConfig.AllowCredentials = true
	server.Use(cors.New(corsConfig))
	server.Use(func(ctx *gin.Context) {
		body, _ := io.ReadAll(ctx.Request.Body)
		common.L().Infow(
			"incomingReq",
			"method", ctx.Request.Method,
			"url", ctx.Request.URL,
			"body", string(body),
		)

		ctx.Request.Body = io.NopCloser(bytes.NewReader(body))
	})

	handler = &Server{
		server: server,
	}

	RegisterPost("api/v1/authz", func() Ctl { return NewAuthzController(uClient) }, &AuthzBody{}, &AuthzData{})
	RegisterPost("api/v1/login", func() Ctl { return NewLoginController(uClient) }, &LoginBody{}, &LoginData{})
	RegisterPost("api/v1/signup", func() Ctl { return NewSignupController(uClient) }, &SignupBody{}, &SignupData{})

	server.Run(conf.Addr)
	return server
}

func RegisterPost(ep string, newCtlFn NewCtlFn, req, resp any) int {
	handler.registry(MethodPost, ep, newCtlFn, req, resp)
	return 0
}

func RegisterGet(ep string, newCtlFn NewCtlFn, req, resp any) {
	handler.registry(MethodGet, ep, newCtlFn, req, resp)
}

func RegisterPut(ep string, newCtlFn NewCtlFn, req, resp any) {
	handler.registry(MethodPut, ep, newCtlFn, req, resp)
}

func Authz(ctx *gin.Context, reqCtx context.Context) (*proto.AuthzResp, common.Code) {
	return handler.authz(ctx, reqCtx)
}
