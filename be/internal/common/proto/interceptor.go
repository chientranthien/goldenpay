package proto

import (
	"context"
	"strconv"
	"time"

	"github.com/chientranthien/goldenpay/internal/common"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

var (
	L = common.L
)

var ServerLoggingInterceptor = grpc.ChainUnaryInterceptor(
	func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		L().Infow("incomingReq", "cmd", info.FullMethod, "req", req)
		return handler(ctx, req)
	},
)

var ServerMetricInterceptor = grpc.ChainUnaryInterceptor(
	func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		start := time.Now()
		defer func() {
			api := info.FullMethod
			c := strconv.Itoa(int(status.Code(err)))
			serverRequestLatency.WithLabelValues(api, c).Observe(float64(time.Now().Sub(start).Milliseconds()))
			serverRequestCount.WithLabelValues(api, c).Inc()
		}()

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
var ClientMetricInterceptor = grpc.WithChainUnaryInterceptor(
	func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var err error
		start := time.Now()
		defer func() {
			c := strconv.Itoa(int(status.Code(err)))
			serverRequestLatency.WithLabelValues(method, c).Observe(float64(time.Now().Sub(start).Milliseconds()))
			serverRequestCount.WithLabelValues(method, c).Inc()
		}()

		err = invoker(ctx, method, req, reply, cc, opts...)
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
