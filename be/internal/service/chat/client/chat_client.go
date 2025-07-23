package client

import (
	"github.com/chientranthien/goldenpay/internal/common"
	commonproto "github.com/chientranthien/goldenpay/internal/common/proto"
	"github.com/chientranthien/goldenpay/internal/proto"
)

func NewChatServiceClient(addr string) proto.ChatServiceClient {
	conn, err := commonproto.NewDial(addr)
	if err != nil {
		common.L().Fatalw("dialChatServiceErr", "addr", addr, "err", err)
	}

	return proto.NewChatServiceClient(conn)
}
