package common

import (
	"context"
	"time"
)

// Ctx to get default context
func Ctx() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), time.Minute)
	return ctx
}

// GetUserIdFromCtx extracts user ID from context (set by HTTP gateway authentication)
func GetUserIdFromCtx(ctx context.Context) uint64 {
	if userId, ok := ctx.Value("userId").(uint64); ok {
		return userId
	}
	return 0
}
