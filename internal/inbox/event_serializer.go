package inbox

type EventSerializer interface {
	Unmarshal(data []byte, v any) error
}
