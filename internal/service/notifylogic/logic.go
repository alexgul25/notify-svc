package notifylogic

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/alexgul25/notify-svc/internal/domain/events"
	"github.com/alexgul25/notify-svc/internal/domain/models"
)

type UserClient interface {
	GetFollowers(ctx context.Context, userID string) (followers []models.Follower, err error)
}

type Notifier interface {
	Notify(recipient string, msg string)
}

type NotifyLogic struct {
	log        *slog.Logger
	userClient UserClient
	notifier   Notifier
}

func New(log *slog.Logger, userClient UserClient, notifier Notifier) *NotifyLogic {
	return &NotifyLogic{
		log:        log,
		userClient: userClient,
		notifier:   notifier,
	}
}

func (nl *NotifyLogic) PlaceCreated(ctx context.Context, event events.PlaceCreated) error {
	const op = "NotifyLogic.PlaceCreated"

	log := nl.log.With(
		slog.String("source", op),
		slog.String("place_id", event.PlaceID),
	)

	log.Info("attempting to get followers")

	followers, err := nl.userClient.GetFollowers(ctx, event.UserID)
	if err != nil {
		log.Error("failed to get followers", slog.Any("error", err))

		return err
	}

	log.Info("followers are gotten successfully")

	log.Info("start notification sending")

	for _, follower := range followers {
		msg := fmt.Sprintf(
			"Hello, %s! User from your subscription list has added a new place, rather invite him on a date!",
			follower.DisplayName,
		)
		nl.notifier.Notify(follower.Email, msg)
	}

	log.Info("finish notification sending")

	return nil
}
