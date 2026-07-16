package grpcclient

import (
	"context"
	"log/slog"
	"time"

	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

func NewAddingHeadersInterceptor(kv []string) grpc.UnaryClientInterceptor {
	return grpc.UnaryClientInterceptor(
		func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			ctx = metadata.AppendToOutgoingContext(ctx, kv...)

			return invoker(ctx, method, req, reply, cc, opts...)
		},
	)
}

func interceptorLogger(log *slog.Logger) grpclog.Logger {
	return grpclog.LoggerFunc(func(ctx context.Context, level grpclog.Level, msg string, fields ...any) {
		var lvl slog.Level
		switch level {
		case grpclog.LevelInfo:
			lvl = slog.LevelInfo
		case grpclog.LevelDebug:
			lvl = slog.LevelDebug
		case grpclog.LevelWarn:
			lvl = slog.LevelWarn
		case grpclog.LevelError:
			lvl = slog.LevelError
		default:
			lvl = slog.LevelInfo
		}

		log.Log(ctx, lvl, msg, fields...)
	})
}

func NewLoggingInterceptor(log *slog.Logger, headersToLog []string) grpc.UnaryClientInterceptor {
	logOpts := []grpclog.Option{
		grpclog.WithFieldsFromContext(func(ctx context.Context) grpclog.Fields {
			md, ok := metadata.FromOutgoingContext(ctx)
			if !ok {
				return nil
			}

			fields := grpclog.Fields{}
			for _, header := range headersToLog {
				if values := md.Get(header); len(values) != 0 {
					fields = append(fields, header, values[0])
				}
			}

			return fields
		}),
	}

	return grpclog.UnaryClientInterceptor(interceptorLogger(log), logOpts...)
}

func NewRetryInterceptor(retriesCount int, timeout time.Duration) grpc.UnaryClientInterceptor {
	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.Unavailable, codes.ResourceExhausted, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(retriesCount)),
		grpcretry.WithPerRetryTimeout(timeout),
	}

	return grpcretry.UnaryClientInterceptor(retryOpts...)
}
