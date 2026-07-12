package inbox

import (
	"strconv"
	"time"
)

const (
	StatusPending   = "pending"
	StatusProcessed = "processed"
)

type Record struct {
	ID            string
	Topic         string
	ProcessStatus string
	Payload       []byte
	CreatedAt     time.Time
}

func ToRecordID(topic string, partition int32, offset int64) string {
	return topic + "-" + strconv.FormatInt(int64(partition), 10) + "-" + strconv.FormatInt(offset, 10)
}
