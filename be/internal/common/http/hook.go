package http

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/chientranthien/goldenpay/internal/common"
)

type (
	HookResult int
	Hook       func(ctx context.Context) common.Code
)

var (
	EarlySkipIfAuthenticated = func(ctx *gin.Context, reqCtx context.Context) common.Code {
		_, code := Authz(ctx, reqCtx)
		if code.Success() {
			return common.CodeAuthenticated
		}

		return common.CodeSuccess
	}
)
