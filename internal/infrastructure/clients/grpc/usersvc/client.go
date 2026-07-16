package usersvc

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/alexgul25/notify-svc/internal/domain/models"
	grpcclient "github.com/alexgul25/notify-svc/internal/infrastructure/clients/grpc"
	userv1 "github.com/alexgul25/protos/gen/go/user/v1"
)

type Client struct {
	api  userv1.UserServiceClient
	conn *grpc.ClientConn
}

func NewClient(log *slog.Logger, addr string, timeout time.Duration, retriesCount int, serviceName string) (*Client, error) {
	const op = "usersvc.NewClient"

	kvToAdd := []string{grpcclient.HeaderServiceName, serviceName}

	dialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			grpcclient.NewAddingHeadersInterceptor(kvToAdd),
			grpcclient.NewLoggingInterceptor(log, []string{}),
			grpcclient.NewRetryInterceptor(retriesCount, timeout),
		),
	}

	conn, err := grpc.NewClient(addr, dialOpts...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	grpcClient := userv1.NewUserServiceClient(conn)

	return &Client{
		api:  grpcClient,
		conn: conn,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) GetFollowers(ctx context.Context, userID string) ([]models.Follower, error) {
	const op = "usersvc.Client.GetFollowers"

	resp, err := c.api.GetFollowers(ctx, &userv1.GetFollowersRequest{UserId: userID})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	followers := make([]models.Follower, len(resp.Followers))
	for i, follower := range resp.Followers {
		followers[i] = models.Follower{ID: follower.GetUserId(), Email: follower.GetEmail(), DisplayName: follower.GetDisplayName()}
	}

	return followers, nil
}
