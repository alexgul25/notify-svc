package app

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"sync"

	"github.com/alexgul25/notify-svc/internal/app/eventpoller"
	"github.com/alexgul25/notify-svc/internal/config"
	"github.com/alexgul25/notify-svc/internal/inbox"
	"github.com/alexgul25/notify-svc/internal/infrastructure/clients/grpc/usersvc"
	"github.com/alexgul25/notify-svc/internal/infrastructure/kafka"
	"github.com/alexgul25/notify-svc/internal/infrastructure/notifier"
	"github.com/alexgul25/notify-svc/internal/infrastructure/serializer"
	"github.com/alexgul25/notify-svc/internal/service/notifylogic"
	"github.com/alexgul25/notify-svc/internal/storage/postgresql"
)

type App struct {
	log       *slog.Logger
	processor *inbox.ConsumerProcessor
	poller    *eventpoller.Poller
	storage   io.Closer
	consumer  io.Closer
	client    io.Closer
}

func New(log *slog.Logger, cfg *config.Config) (*App, error) {
	storage, err := postgresql.NewStorage(
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.DbName,
		cfg.Database.Port)
	if err != nil {
		return nil, fmt.Errorf("failed to init storage: %w", err)
	}

	inboxStorage := postgresql.NewInboxStorage(storage.DB())

	consumer, err := kafka.NewConsumer(cfg.KafkaConsumer.Brokers, inbox.TopicPlaceCreated)
	if err != nil {
		return nil, fmt.Errorf("failed to init consumer: %w", err)
	}

	processor := inbox.NewConsumerProcessor(
		inboxStorage,
		consumer,
		log,
		cfg.InboxProcessor.OpTimeout,
	)

	userClient, err := usersvc.NewClient(
		log,
		cfg.GRPCClient.UserServiceAddr,
		cfg.GRPCClient.UserServiceTimeout,
		cfg.GRPCClient.UserServiceRetriesCount,
		cfg.ServiceName,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to init grpc user client: %w", err)
	}

	logNotifier := notifier.NewLogNotifier(log)

	logic := notifylogic.New(log, userClient, logNotifier)

	poller := eventpoller.New(
		log,
		inbox.TopicPlaceCreated,
		cfg.EventPoller.Limit,
		cfg.EventPoller.NotifyInterval,
		cfg.EventPoller.Timeout,
		logic,
		inboxStorage,
		serializer.JSONSerializer{},
	)

	return &App{
		log:       log,
		processor: processor,
		poller:    poller,
		storage:   storage,
		consumer:  consumer,
		client:    userClient,
	}, nil
}

func (a *App) Start() {
	const op = "App.Start"

	a.log.Info("starting app", slog.String("source", op))

	go a.processor.Start()
	go a.poller.Run()

	a.log.Info("app is running", slog.String("source", op))
}

func (a *App) GracefulShutdown() {
	const op = "App.GracefulShutdown"

	log := a.log.With(slog.String("source", op))

	log.Info("app is shutting down gracefully")

	wg := &sync.WaitGroup{}

	wg.Go(a.processor.Shutdown)
	wg.Go(a.poller.GracefulStop)

	wg.Wait()

	err := errors.Join(
		a.consumer.Close(),
		a.client.Close(),
		a.storage.Close(),
	)
	if err != nil {
		log.Error("failed to close app components", slog.Any("error", err))
	}

	log.Info("app has shut down gracefully")
}
