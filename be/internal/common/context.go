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
