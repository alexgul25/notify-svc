package kafka

import (
	"errors"
	"fmt"
	"sync"

	"github.com/IBM/sarama"
	"github.com/alexgul25/notify-svc/internal/inbox"
)

type Consumer struct {
	consumer          sarama.Consumer
	partitionConsumer sarama.PartitionConsumer
	messages          chan inbox.Message
	wg                *sync.WaitGroup
	done              chan struct{}
}

func NewConsumer(brokers []string, topic string) (*Consumer, error) {
	const op = "kafka.NewConsumer"

	config := sarama.NewConfig()

	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Return.Errors = false

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	messages := make(chan inbox.Message)
	wg := &sync.WaitGroup{}
	done := make(chan struct{})

	wg.Go(func() {
		for {
			select {
			case <-done:
				return
			case msg := <-partitionConsumer.Messages():
				msgID := getHeaderValue(msg.Headers, "message_id")
				inboxMsg := inbox.Message{
					ID:        msgID,
					Topic:     msg.Topic,
					Partition: msg.Partition,
					Offset:    msg.Offset,
					Key:       msg.Key,
					Value:     msg.Value,
				}
				messages <- inboxMsg
			}
		}
	})

	return &Consumer{
		consumer:          consumer,
		partitionConsumer: partitionConsumer,
		messages:          messages,
		wg:                wg,
		done:              done,
	}, nil
}

func getHeaderValue(headers []*sarama.RecordHeader, key string) string {
	for _, h := range headers {
		if string(h.Key) == key {
			return string(h.Value)
		}
	}
	return ""
}

func (c *Consumer) Messages() <-chan inbox.Message {
	return c.messages
}

func (c *Consumer) Close() error {
	const op = "kafka.Consumer.Close"

	select {
	case <-c.done:
		return nil
	default:
		close(c.done)
		c.wg.Wait()

		err := errors.Join(c.partitionConsumer.Close(), c.consumer.Close())

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		return nil
	}
}
