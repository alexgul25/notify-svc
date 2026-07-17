package notifier

import (
	"log/slog"
)

type LogNotifier struct {
	log *slog.Logger
}

func NewLogNotifier(log *slog.Logger) *LogNotifier {
	return &LogNotifier{log: log}
}

func (ln *LogNotifier) Notify(recipient string, msg string) {
	ln.log.Info("new notification sent", slog.String("recipient", recipient), slog.String("message", msg))
}
