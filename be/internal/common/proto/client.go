package proto

import (
	"github.com/chientranthien/goldenpay/internal/common"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	_ "google.golang.org/grpc/xds"
)

func NewDial(addr string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		ClientLoggingInterceptor,
		ClientMetricInterceptor,
	)
	if err != nil {
		common.L().Errorw("dialErr", "addr", addr, "err", err)
		return nil, err
	} else {
		common.L().Infow("dialed", "addr", addr)
	}

	return conn, nil
}
