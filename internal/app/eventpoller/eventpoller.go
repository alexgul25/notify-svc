package eventpoller

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/alexgul25/notify-svc/internal/domain/events"
	"github.com/alexgul25/notify-svc/internal/inbox"
	"github.com/alexgul25/notify-svc/internal/service/notifylogic"
)

type InboxSelector interface {
	SelectPending(ctx context.Context, topic string, limit int) ([]inbox.Record, error)
	MarkAsProcessed(ctx context.Context, id string) error
}

type Poller struct {
	log            *slog.Logger
	topic          string
	limit          int
	notifyInterval time.Duration
	timeout        time.Duration
	logic          *notifylogic.NotifyLogic
	selector       InboxSelector
	serializer     inbox.EventSerializer
	wg             *sync.WaitGroup
	done           chan struct{}
}

func New(
	log *slog.Logger,
	topic string,
	limit int,
	notifyInterval time.Duration,
	timeout time.Duration,
	logic *notifylogic.NotifyLogic,
	selector InboxSelector,
	serializer inbox.EventSerializer,
) *Poller {
	return &Poller{
		log:            log,
		topic:          topic,
		limit:          limit,
		notifyInterval: notifyInterval,
		timeout:        timeout,
		logic:          logic,
		selector:       selector,
		serializer:     serializer,
		wg:             &sync.WaitGroup{},
		done:           make(chan struct{}),
	}
}

func (p *Poller) Run() {
	const op = "Poller.Run"

	log := p.log.With(slog.String("source", op), slog.String("topic", p.topic))

	log.Info("notify poller has started")

	p.wg.Add(1)
	defer p.wg.Done()

	ticker := time.NewTicker(p.notifyInterval)
	defer ticker.Stop()

	for {
		select {
		case <-p.done:
			return
		case <-ticker.C:
			log.Info("notify poller new cycle")

			selectCtx, selectCancel := context.WithTimeout(context.Background(), p.timeout)
			records, err := p.selector.SelectPending(selectCtx, p.topic, p.limit)
			if err != nil {
				log.Error("failed to select records", slog.Any("error", err))
			}
			selectCancel()

			for _, record := range records {
				var event events.PlaceCreated
				err := p.serializer.Unmarshal(record.Payload, &event)
				if err != nil {
					log.Error("failed to unmarshal record", slog.String("record_id", record.ID), slog.Any("error", err))
				} else {
					logicCtx, logicCancel := context.WithTimeout(context.Background(), p.timeout)
					err := p.logic.PlaceCreated(logicCtx, event)
					if err != nil {
						log.Error("failed to send notification", slog.String("record_id", record.ID), slog.Any("error", err))
					} else {
						markCtx, markCancel := context.WithTimeout(context.Background(), p.timeout)
						err := p.selector.MarkAsProcessed(markCtx, record.ID)
						if err != nil {
							log.Error("failed to mark record as processed", slog.String("record_id", record.ID), slog.Any("error", err))
						}
						markCancel()
					}
					logicCancel()
				}
			}
		}
	}
}

func (p *Poller) GracefulStop() {
	const op = "Poller.GracefulStop"

	log := p.log.With(slog.String("source", op), slog.String("topic", p.topic))

	select {
	case <-p.done:
		return
	default:
		log.Info("stopping event poller")

		close(p.done)
		p.wg.Wait()

		log.Info("event poller gracefully stopped")
	}
}
