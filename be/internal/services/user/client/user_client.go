package client

import (
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/chientranthien/goldenpay/internal/common"
	"github.com/chientranthien/goldenpay/internal/proto"
)

var (
	clientOnce = sync.Once{}
	client     proto.UserServiceClient
)

func NewUserServiceClient(addr string) proto.UserServiceClient {
	clientOnce.Do(func() {
		conn, err := grpc.Dial(
			addr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			common.ClientLoggingInterceptor,
		)
		if err != nil {
			common.L().Errorw("dialErr", "addr", addr, "err", err)
		} else {
			common.L().Infow("dialed", "addr", addr)
		}

		client = proto.NewUserServiceClient(conn)
	})

	return client
}
