package client

import (
	"log"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/chientranthien/goldenpay/internal/proto"
)

var (
	clientOnce = sync.Once{}
	client     proto.WalletServiceClient
)

func NewWalletServiceClient(addr string) proto.WalletServiceClient {
	clientOnce.Do(func() {
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("failed to dial, add=%v, err=%v", addr, err)
		} else {
			log.Printf("dialed to addr=%v", addr)
		}

		client = proto.NewWalletServiceClient(conn)
	})

	return client
}
