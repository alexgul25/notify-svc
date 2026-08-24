package inbox

import (
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
