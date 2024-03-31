package client

import (
	"sync"

	commonproto "github.com/chientranthien/goldenpay/internal/common/proto"
	"github.com/chientranthien/goldenpay/internal/proto"
)

var (
	clientOnce = sync.Once{}
	client     proto.WalletServiceClient
)

func NewWalletServiceClient(addr string) proto.WalletServiceClient {
	clientOnce.Do(func() {
		conn, _ := commonproto.NewDial(addr)
		client = proto.NewWalletServiceClient(conn)
	})

	return client
}
