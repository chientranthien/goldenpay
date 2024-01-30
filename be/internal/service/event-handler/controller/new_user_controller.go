package controller

import (
	"encoding/json"
	"time"

	"github.com/IBM/sarama"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/chientranthien/goldenpay/internal/common"
	"github.com/chientranthien/goldenpay/internal/proto"
	"github.com/chientranthien/goldenpay/internal/service/event_handler/config"
)

type NewUserController struct {
	wClient proto.WalletServiceClient
}

func NewNewUserController(wClient proto.WalletServiceClient) *NewUserController {
	return &NewUserController{wClient: wClient}
}

func (c NewUserController) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c NewUserController) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c NewUserController) ConsumeClaim(
	session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
) error {
	for {
		select {
		case m, ok := <-claim.Messages():
			if !ok {
				common.L().Infow(
					"messageChannelBroke",
					"topic",
					claim.Topic(),
					"partition",
					claim.Partition(),
				)
				return nil
			}
			c.handleMessage(m)
			session.MarkMessage(m, "")
		case <-session.Context().Done():
			common.L().Infow(
				"sessionDone",
				"topic",
				claim.Topic(),
				"partition",
				claim.Partition(),
			)
			return nil
		}
	}
}

func (c NewUserController) handleMessage(m *sarama.ConsumerMessage) {
	event := &proto.NewUserEvent{}
	start := time.Now()
	defer func() {
		elapsedSinceEventCreated := common.NowMicro() - event.EventTime
		elapsed := time.Since(start).Microseconds()
		common.L().Infow(
			"handleNewTransaction",
			"partition", m.Partition,
			"offset", m.Offset,
			"elapsed", elapsed,
			"elapsedSinceEventCreated", elapsedSinceEventCreated,
		)
	}()

	ctx := common.Ctx()
	err := json.Unmarshal(m.Value, event)
	if err != nil {
		common.L().Warnw("marshalNewUserEventErr", "value", string(m.Value), "err", err)
		return
	}

	_, err = c.wClient.Create(ctx, &proto.CreateWalletReq{
		UserId:         event.UserId,
		InitialBalance: config.Get().General.InitialBalance,
	})

	if err != nil && status.Code(err) != codes.AlreadyExists {
		common.L().Warnw("createWalletErr", "err", err)
		return
	}
}
