package arpcdata

import (
	"time"

	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/durationpb"
)

func DurationOptional(src *time.Duration) *durationpb.Duration {
	if src == nil {
		return nil
	}

	return durationpb.New(*src)
}

func DurationOptionalProto(duration *durationpb.Duration) *time.Duration {
	if duration == nil {
		return nil
	}

	return lo.ToPtr(duration.AsDuration())
}
