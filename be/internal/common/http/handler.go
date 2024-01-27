package http

import (
	"context"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/chientranthien/goldenpay/internal/common"
	"github.com/chientranthien/goldenpay/internal/proto"
)

type (
	Server struct {
		conf      common.ServiceConfig
		ginServer *gin.Engine
		uClient   proto.UserServiceClient
	}

	Cookie struct {
		Name     string
		Value    string
		MaxAge   int
		Path     string
		Domain   string
		Secure   bool
		HttpOnly bool
	}
)

const (
	CtxHeaderGin = "gin"
)

func NewServer(uClient proto.UserServiceClient) *Server {
	return &Server{uClient: uClient}
}

func (s Server) Ctx(ginCtx *gin.Context) context.Context {
	ctx := common.Ctx()
	return context.WithValue(ctx, CtxHeaderGin, ginCtx)
}

func GetGinCtx(ctx context.Context) *gin.Context {
	return ctx.Value(CtxHeaderGin).(*gin.Context)
}

func SetCookie(ctx context.Context, c Cookie) {
	GetGinCtx(ctx).SetCookie(c.Name, c.Value, c.MaxAge, c.Path, c.Domain, c.Secure, c.HttpOnly)
}

func (s Server) newHandler(epInfo EndpointInfo) func(ctx *gin.Context) {
	return func(ginCtx *gin.Context) {
		ctx := s.Ctx(ginCtx)

		req := Req{}
		if s.needAuthentication(epInfo) {
			authzResp, code := s.authz(ginCtx, ctx)

			if !code.Success() {
				ginCtx.JSON(http.StatusOK, &Resp{Code: code})
			}

			metadata := authzResp.GetMetadata()
			req.Metadata = Metadata{UserId: metadata.GetUserId(), Email: metadata.GetEmail()}
		}

		c := epInfo.NewCtlFn()
		var body any
		if epInfo.Req != nil {
			t := reflect.TypeOf(epInfo.Req).Elem()
			body = reflect.New(t).Interface()
		}
		req.Body = body
		ginCtx.BindJSON(req.Body)

		for _, hook := range epInfo.PreReqHooks {
			if code := hook(ctx); !code.Success() {
				ginCtx.JSON(http.StatusOK, &Resp{Code: code})
				return
			}
		}
		if code := c.Take(ctx, req); !code.Success() {
			resp := &Resp{
				Code: code,
			}
			ginCtx.JSON(http.StatusOK, resp)
			return
		}

		data, code := c.Do()
		resp := Resp{Code: code, Data: data}

		ginCtx.JSON(http.StatusOK, resp)
	}
}

// TODO(tom): make this meth more generic
func (s Server) needAuthentication(info EndpointInfo) bool {
	ep := strings.ToLower(info.EP)
	if strings.Contains(ep, "login") || strings.Contains(ep, "signup") {
		return false
	}

	return true
}

func (s Server) authz(ginCtx *gin.Context, ctx context.Context) (*proto.AuthzResp, common.Code) {
	token, _ := ginCtx.Cookie(CookieToken)
	if token == "" {
		return nil, common.CodeUnauthenticated
	}

	authzResp, err := s.uClient.Authz(ctx, &proto.AuthzReq{Token: token})
	if err != nil {
		return nil, common.GetCodeFromErr(err)
	}

	return authzResp, common.CodeSuccess
}

func (s Server) register(epInfo EndpointInfo) {
	s.ginServer.Handle(epInfo.method, epInfo.EP, server.newHandler(epInfo))
}
