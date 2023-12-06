package controller

import (
	"encoding/json"
	"time"

	"github.com/IBM/sarama"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/chientranthien/goldenpay/internal/common"
	"github.com/chientranthien/goldenpay/internal/proto"
)

type NewTransactionController struct {
	wClient proto.WalletServiceClient
}

func NewNewTransactionController(wClient proto.WalletServiceClient) *NewTransactionController {
	return &NewTransactionController{wClient: wClient}
}

func (c NewTransactionController) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c NewTransactionController) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c NewTransactionController) ConsumeClaim(
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

func (c NewTransactionController) handleMessage(m *sarama.ConsumerMessage) {
	start := time.Now()
	event := &proto.NewTransactionEvent{}
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
		common.L().Warnw("marshalNewTransactionEventErr", "value", string(m.Value), "err", err)
		return
	}

	_, err = c.wClient.ProcessTransfer(ctx, &proto.ProcessTransferReq{
		TransactionId: event.TransactionId,
	})

	if err != nil && status.Code(err) != codes.AlreadyExists {
		common.L().Warnw("createWalletErr", "err", err)
		return
	}

	if err != nil {
		common.L().Errorw("processTransferErr", "event", event, "err", err)
	}

}
