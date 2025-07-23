package proto

import (
	"net"
	"os"

	"github.com/chientranthien/goldenpay/internal/common"
	"github.com/chientranthien/goldenpay/internal/common/metric"
	"google.golang.org/grpc"
	"google.golang.org/grpc/xds"
)

var (
	_ grpc.ServiceRegistrar = (*server)(nil)
)

type serve interface {
	Serve(lis net.Listener) error
	RegisterService(desc *grpc.ServiceDesc, impl any)
}

type server struct {
	s    serve
	addr string
}

func NewServer(addr string) (*server, error) {
	env := os.Getenv("G_ENV")
	if env == "local" || env == "" {
		s := grpc.NewServer(
			ServerLoggingInterceptor,
			// TODO(tom): consider to suer go-grpc-middleware
			ServerMetricInterceptor,
		)
		return &server{s: s, addr: addr}, nil
	}

	s, err := xds.NewGRPCServer(
		ServerLoggingInterceptor,
		// TODO(tom): consider to suer go-grpc-middleware
		ServerMetricInterceptor,
	)
	if err != nil {
		common.L().Errorw("newGRPCServerErr", "err", err)
		return nil, err
	}

	return &server{s: s, addr: addr}, nil
}

func (s server) RegisterService(desc *grpc.ServiceDesc, impl any) {
	s.s.RegisterService(desc, impl)
}

func (s server) ListenAndServe() error {
	go metric.ServeDefault()

	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		common.L().Errorw("netListenErr", "addr", s.addr, "err", err)
		return err
	} else {
		common.L().Infow("listening", "addr", s.addr)
	}

	err = s.s.Serve(lis)
	if err != nil {
		common.L().Errorw("serveErr", "err", err)
		return err
	}

	return nil
}
