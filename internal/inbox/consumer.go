package inbox

type Message struct {
	Topic     string
	Partition int32
	Offset    int64
	Key       []byte
	Value     []byte
}

type MessageConsumer interface {
	Messages() <-chan Message
	Close() error
}
