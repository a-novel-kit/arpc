package arpcdata_test

import (
	"testing"
	"time"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/durationpb"

	arpcdata "github.com/a-novel-kit/arpc/data"
)

func TestDurationOptional(t *testing.T) {
	duration := arpcdata.DurationOptional(lo.ToPtr(5 * time.Second))
	require.NotNil(t, duration)
	require.Equal(t, 5*time.Second, duration.AsDuration())

	duration = arpcdata.DurationOptional(nil)
	require.Nil(t, duration)
}

func TestDurationOptionalProto(t *testing.T) {
	duration := arpcdata.DurationOptionalProto(durationpb.New(5 * time.Second))
	require.NotNil(t, duration)
	require.Equal(t, 5*time.Second, *duration)

	duration = arpcdata.DurationOptionalProto(nil)
	require.Nil(t, duration)
}
