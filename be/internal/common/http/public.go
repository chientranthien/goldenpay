package http

import (
	"bytes"
	"context"
	"io"

	"github.com/chientranthien/goldenpay/internal/common/metric"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/chientranthien/goldenpay/internal/common"
	"github.com/chientranthien/goldenpay/internal/proto"
)

type (
	Resp struct {
		Code common.Code `json:"code"`
		Data any         `json:"data"`
	}

	Metadata struct {
		UserId uint64
		Email  string
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

	EndpointInfoPost struct {
		EndpointInfo
	}

	EndpointInfo struct {
		method      string
		EP          string
		NewCtlFn    NewCtlFn
		Req         common.AnyPtr
		Resp        common.AnyPtr
		PreReqHooks []Hook
	}

	PostEndpointInfo EndpointInfo
	GetEndpointInfo  EndpointInfo
	PutEndpointInfo  EndpointInfo
)

func (i *PostEndpointInfo) Ensure() {
	i.method = MethodPost
}

func (i *GetEndpointInfo) Ensure() {
	i.method = MethodGet
}

func (i *PutEndpointInfo) Ensure() {
	i.method = MethodPut
}

var (
	server *Server
)

const (
	MethodGet  = "GET"
	MethodPost = "POST"
	MethodPut  = "PUT"

	CookieToken  = "token"
	CookieUserId = "uid"
)

func Init(conf common.ServiceConfig, uClient proto.UserServiceClient) {
	ginServer := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000", "https://goldenpay.chientran.info"}
	corsConfig.AllowCredentials = true
	ginServer.Use(cors.New(corsConfig))
	ginServer.Use(func(ctx *gin.Context) {
		body, _ := io.ReadAll(ctx.Request.Body)
		common.L().Infow(
			"incomingReq",
			"method", ctx.Request.Method,
			"url", ctx.Request.URL,
			"body", string(body),
		)

		ctx.Request.Body = io.NopCloser(bytes.NewReader(body))
	})

	server = &Server{
		conf:      conf,
		ginServer: ginServer,
		uClient:   uClient,
	}

	RegisterPost(PostEndpointInfo{
		EP:       "api/v1/authz",
		NewCtlFn: func() Ctl { return NewAuthzController(uClient) },
		Req:      &AuthzBody{},
		Resp:     &AuthzData{},
	})
	RegisterPost(PostEndpointInfo{
		EP:          "api/v1/login",
		NewCtlFn:    func() Ctl { return NewLoginController(uClient) },
		Req:         &LoginBody{},
		Resp:        &LoginData{},
		PreReqHooks: []Hook{EarlySkipIfAuthenticated},
	})
	RegisterPost(PostEndpointInfo{
		EP:          "api/v1/signup",
		NewCtlFn:    func() Ctl { return NewSignupController(uClient) },
		Req:         &SignupBody{},
		Resp:        &SignupData{},
		PreReqHooks: []Hook{EarlySkipIfAuthenticated},
	},
	)

}

func Run() {
	go metric.ServeDefault()
	server.ginServer.Run(server.conf.Addr)
}

func RegisterPost(epInfo PostEndpointInfo) {
	epInfo.Ensure()
	server.register(EndpointInfo(epInfo))
}

func RegisterGet(epInfo GetEndpointInfo) {
	epInfo.Ensure()
	server.register(EndpointInfo(epInfo))
}

func RegisterPut(epInfo PutEndpointInfo) {
	epInfo.Ensure()
	server.register(EndpointInfo(epInfo))
}

func authz(ginCtx *gin.Context, ctx context.Context) (*proto.AuthzResp, common.Code) {
	return server.authz(ginCtx, ctx)
}
