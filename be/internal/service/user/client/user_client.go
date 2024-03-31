package client

import (
	"sync"

	commonproto "github.com/chientranthien/goldenpay/internal/common/proto"
	"github.com/chientranthien/goldenpay/internal/proto"
)

var (
	clientOnce = sync.Once{}
	client     proto.UserServiceClient
)

func NewUserServiceClient(addr string) proto.UserServiceClient {
	clientOnce.Do(func() {
		conn, _ := commonproto.NewDial(addr)
		client = proto.NewUserServiceClient(conn)
	})

	return client
}
