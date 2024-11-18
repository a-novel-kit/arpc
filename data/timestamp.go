package arpcdata

import (
	"time"

	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TimestampOptional(src *time.Time) *timestamppb.Timestamp {
	if src == nil {
		return nil
	}

	return timestamppb.New(*src)
}

func TimestampOptionalProto(timestamp *timestamppb.Timestamp) *time.Time {
	if timestamp == nil {
		return nil
	}

	return lo.ToPtr(timestamp.AsTime())
}
