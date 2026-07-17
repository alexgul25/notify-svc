package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Env            string `envconfig:"ENV"`
	ServiceName    string `envconfig:"SERVICE_NAME" env-default:"notify-svc"`
	Database       DatabaseConfig
	InboxProcessor InboxProcessorConfig
	KafkaConsumer  KafkaConsumerConfig
	GRPCClient     GRPCClientConfig
	EventPoller    EventPollerConfig
}

type DatabaseConfig struct {
	User     string `envconfig:"DB_USER"`
	Password string `envconfig:"DB_PASSWORD"`
	DbName   string `envconfig:"DB_NAME"`
	Host     string `envconfig:"DB_HOST"`
	Port     int    `envconfig:"DB_PORT"`
}

type InboxProcessorConfig struct {
	OpTimeout time.Duration `envconfig:"INBOX_PROCESSOR_OP_TIMEOUT" env-default:"5s"`
}

type KafkaConsumerConfig struct {
	BrokersRaw string `envconfig:"KAFKA_CONSUMER_BROKERS"`
	Brokers    []string
}

type GRPCClientConfig struct {
	UserServiceAddr         string        `envconfig:"USER_SERVICE_ADDR"`
	UserServiceTimeout      time.Duration `envconfig:"USER_SERVICE_TIMEOUT" env-default:"5s"`
	UserServiceRetriesCount int           `envconfig:"USER_SERVICE_RETRY_COUNT" env-default:"3"`
}

type EventPollerConfig struct {
	Limit          int           `envconfig:"EVENT_POLLER_LIMIT" env-default:"50"`
	NotifyInterval time.Duration `envconfig:"EVENT_POLLER_NOTIFY_INTERVAL" env-default:"10s"`
	Timeout        time.Duration `envconfig:"EVENT_POLLER_TIMEOUT" env-default:"5s"`
}

func load() (*Config, error) {
	const op = "load"

	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var cfg Config
	err = envconfig.Process("", &cfg)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &cfg, nil
}

func LoadNotifyService() (*Config, error) {
	const op = "LoadNotifyService"

	cfg, err := load()
	if err != nil {
		return nil, err
	}

	if cfg.Env == "" {
		return nil, fmt.Errorf("%s: env variable ENV not set", op)
	}
	if cfg.Database.User == "" {
		return nil, fmt.Errorf("%s: env variable DB_USER not set", op)
	}
	if cfg.Database.Password == "" {
		return nil, fmt.Errorf("%s: env variable DB_PASSWORD not set", op)
	}
	if cfg.Database.DbName == "" {
		return nil, fmt.Errorf("%s: env variable DB_NAME not set", op)
	}
	if cfg.Database.Host == "" {
		return nil, fmt.Errorf("%s: env variable DB_HOST not set", op)
	}
	if cfg.Database.Port == 0 {
		return nil, fmt.Errorf("%s: env variable DB_PORT not set", op)
	}
	if cfg.KafkaConsumer.BrokersRaw == "" {
		return nil, fmt.Errorf("%s: env variable KAFKA_CONSUMER_BROKERS not set", op)
	}
	if cfg.GRPCClient.UserServiceAddr == "" {
		return nil, fmt.Errorf("%s: env variable USER_SERVICE_ADDR not set", op)
	}

	cfg.KafkaConsumer.Brokers = strings.Split(cfg.KafkaConsumer.BrokersRaw, ",")

	return cfg, nil
}
