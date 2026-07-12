package inbox

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/alexgul25/notify-svc/internal/domain"
)

type InboxRepository interface {
	InsertRecord(ctx context.Context, record Record) error
}

type ConsumerProcessor struct {
	repo      InboxRepository
	consumer  MessageConsumer
	log       *slog.Logger
	opTimeout time.Duration
	done      chan struct{}
	wg        *sync.WaitGroup
}

func NewConsumerProcessor(
	repo InboxRepository,
	consumer MessageConsumer,
	log *slog.Logger,
	opTimeout time.Duration,
) *ConsumerProcessor {
	return &ConsumerProcessor{
		repo:      repo,
		consumer:  consumer,
		log:       log,
		opTimeout: opTimeout,
		done:      make(chan struct{}),
		wg:        &sync.WaitGroup{},
	}
}

func (cp *ConsumerProcessor) Start() {
	const op = "ConsumerProcessor.Start"

	log := cp.log.With(
		slog.String("source", op),
	)

	log.Info("consumer processor is started", slog.String("source", op))

	cp.wg.Add(1)
	defer cp.wg.Done()

	for {
		select {
		case <-cp.done:
			return
		case msg := <-cp.consumer.Messages():
			log.Info(
				"get msg from consumer",
				slog.String("topic", msg.Topic),
				slog.Int64("partition", int64(msg.Partition)),
				slog.Int64("offset", msg.Offset),
			)

			record := Record{
				ID:            ToRecordID(msg.Topic, msg.Partition, msg.Offset),
				Topic:         msg.Topic,
				ProcessStatus: StatusPending,
				Payload:       msg.Value,
				CreatedAt:     time.Now(),
			}

			log.Info("attempting to insert record about msg")

			insertCtx, insertCancel := context.WithTimeout(context.Background(), cp.opTimeout)
			err := cp.repo.InsertRecord(insertCtx, record)
			if errors.Is(err, domain.ErrMsgDoubleSend) {
				log.Warn("get duplicate msg")
			} else if err != nil {
				log.Error(
					"failed to insert record about msg",
					slog.Any("error", err),
				)
			} else {
				log.Info("record about msg is inserted successfully")
			}
			insertCancel()
		}
	}
}

func (cp *ConsumerProcessor) Shutdown() {
	const op = "ConsumerProcessor.Shutdown"

	select {
	case <-cp.done:
		return
	default:
		cp.log.Info("consumer processor is shutting down", slog.String("source", op))

		close(cp.done)
		cp.wg.Wait()

		cp.log.Info("consumer processor is shutted down successfully", slog.String("source", op))
	}
}
