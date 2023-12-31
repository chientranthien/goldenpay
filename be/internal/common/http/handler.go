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
		uClient proto.UserServiceClient
		server  *gin.Engine
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

func (s Server) Ctx(ctx *gin.Context) context.Context {
	reqCtx := common.Ctx()
	return context.WithValue(reqCtx, CtxHeaderGin, ctx)
}

func GetGinCtx(reqCtx context.Context) *gin.Context {
	return reqCtx.Value(CtxHeaderGin).(*gin.Context)
}

func SetCookie(reqCtx context.Context, c Cookie) {
	GetGinCtx(reqCtx).SetCookie(c.Name, c.Value, c.MaxAge, c.Path, c.Domain, c.Secure, c.HttpOnly)
}

func (s Server) newHandler(epInfo endpointInfo) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		reqCtx := common.Ctx()

		req := Req{}
		if s.needAuthentication(epInfo) {
			authzResp, code := s.authz(ctx, reqCtx)

			if !code.Success() {
				ctx.JSON(http.StatusOK, &CommonResp{Code: code})
			}

			metadata := authzResp.GetMetadata()
			req.Metadata = Metadata{UserId: metadata.GetUserId(), Email: metadata.GetEmail()}
		}

		c := epInfo.newCtlFn()
		var body any
		if epInfo.req != nil {
			t := reflect.TypeOf(epInfo.req)
			body = reflect.New(t).Interface()
		}
		req.Body = body

		for _, hook := range epInfo.preReqHooks {
			if code := hook(reqCtx); !code.Success() {
				ctx.JSON(http.StatusOK, &CommonResp{Code: code})
				return
			}
		}
		if code := c.Take(reqCtx, req); !code.Success() {
			resp := &CommonResp{
				Code: code,
			}
			ctx.JSON(http.StatusOK, resp)
			return
		}

		data, code := c.Do()
		resp := CommonResp{Code: code, Data: data}

		ctx.JSON(http.StatusOK, resp)
	}
}

// TODO(tom): make this meth more generic
func (s Server) needAuthentication(info endpointInfo) bool {
	ep := strings.ToLower(info.ep)
	if strings.Contains(ep, "login") || strings.Contains(ep, "signup") {
		return false
	}

	return true
}

func (s Server) authz(ctx *gin.Context, reqCtx context.Context) (*proto.AuthzResp, common.Code) {
	token, _ := ctx.Cookie(CookieToken)
	if token == "" {
		return nil, common.CodeUnauthenticated
	}

	authzResp, err := s.uClient.Authz(reqCtx, &proto.AuthzReq{Token: token})
	if err != nil {
		return nil, common.GetCodeFromErr(err)
	}

	return authzResp, common.CodeSuccess
}

func (s Server) registry(method, ep string, newCtlFn NewCtlFn, req, resp common.AnyPtr) {
	epInfo := endpointInfo{
		method:   method,
		ep:       ep,
		newCtlFn: newCtlFn,
		req:      req,
		resp:     resp,
	}

	s.server.Handle(method, ep, handler.newHandler(epInfo))
}
