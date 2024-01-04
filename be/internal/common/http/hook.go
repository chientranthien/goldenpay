package http

import (
	"context"

	"github.com/chientranthien/goldenpay/internal/common"
)

type (
	HookResult int
	Hook       func(ctx context.Context) common.Code
)

var (
	EarlySkipIfAuthenticated = func(ctx context.Context) common.Code {
		ginCtx := GetGinCtx(ctx)
		_, code := authz(ginCtx, ctx)
		if code.Success() {
			return common.CodeAuthenticated
		}

		return common.CodeSuccess
	}
)
