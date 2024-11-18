package arpcdata_test

import (
	"testing"
	"time"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	arpcdata "github.com/a-novel-kit/arpc/data"
)

func TestTimestampOptional(t *testing.T) {
	timestamp := arpcdata.TimestampOptional(lo.ToPtr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)))
	require.NotNil(t, timestamp)
	require.Equal(t, time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), timestamp.AsTime())

	timestamp = arpcdata.TimestampOptional(nil)
	require.Nil(t, timestamp)
}

func TestTimestampOptionalProto(t *testing.T) {
	timestamp := arpcdata.TimestampOptionalProto(timestamppb.New(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)))
	require.NotNil(t, timestamp)
	require.Equal(t, time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), *timestamp)

	timestamp = arpcdata.TimestampOptionalProto(nil)
	require.Nil(t, timestamp)
}
