package arpc

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/a-novel-kit/quicklog"

	arpcmessages "github.com/a-novel-kit/arpc/messages"
)

type ExecService[In any, Out any] interface {
	Exec(ctx context.Context, data In) (Out, error)
}

type GRPCCallback[In any, Out any] func(ctx context.Context, in In) (Out, error)

type wrappedService[In any, Out any] struct {
	service GRPCCallback[In, Out]
}

func (s *wrappedService[In, Out]) Exec(ctx context.Context, data In) (Out, error) {
	return s.service(ctx, data)
}

func WithReport[In any, Out any](
	name string, service ExecService[In, Out], logger quicklog.Logger,
) ExecService[In, Out] {
	return &wrappedService[In, Out]{
		service: func(ctx context.Context, in In) (Out, error) {
			start := time.Now()
			out, err := service.Exec(ctx, in)
			end := time.Now()

			level := quicklog.LevelInfo
			if err != nil {
				code := status.Code(err)
				if code == codes.Unavailable || code == codes.Canceled || code == codes.Unimplemented {
					level = quicklog.LevelWarning
				} else {
					level = quicklog.LevelError
				}
			}

			logger.Log(level, arpcmessages.NewReport(
				&arpcmessages.Metrics{Latency: end.Sub(start)},
				name,
				err,
			))

			return out, err //nolint:wrapcheck
		},
	}
}
