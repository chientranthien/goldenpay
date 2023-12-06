package common

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

var ServerLoggingInterceptor = grpc.ChainUnaryInterceptor(
	func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		L().Infow("incomingReq", "cmd", info.FullMethod, "req", req)
		return handler(ctx, req)
	},
)

var ClientLoggingInterceptor = grpc.WithChainUnaryInterceptor(
	func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		start := time.Now()

		err := invoker(ctx, method, req, reply, cc, opts...)
		L().Infow(
			"outgoingReq",
			"cmd", method,
			"req", req,
			"resp", reply,
			"elapsed", time.Since(start).Microseconds(),
			"err", status.Code(err).String(),
		)
		return err
	},
)
