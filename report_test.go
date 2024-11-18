package arpc_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/a-novel-kit/quicklog"
	quicklogmocks "github.com/a-novel-kit/quicklog/mocks"

	"github.com/a-novel-kit/arpc"
)

type fakeExecService struct {
	out string
	err error
}

func (f *fakeExecService) Exec(_ context.Context, _ int) (string, error) {
	return f.out, f.err
}

func TestWithReport(t *testing.T) {
	testCases := []struct {
		name string

		service string
		in      int
		out     string
		err     error

		expectLevel quicklog.Level
	}{
		{
			name: "Success",

			service: "foo",
			in:      42,
			out:     "bar",
			err:     nil,

			expectLevel: quicklog.LevelInfo,
		},
		{
			name: "Unavailable",

			service: "foo",
			in:      42,
			out:     "bar",
			err:     status.Error(codes.Unavailable, "foo"),

			expectLevel: quicklog.LevelWarning,
		},
		{
			name: "Canceled",

			service: "foo",
			in:      42,
			out:     "bar",
			err:     status.Error(codes.Canceled, "foo"),

			expectLevel: quicklog.LevelWarning,
		},
		{
			name: "Unimplemented",

			service: "foo",
			in:      42,
			out:     "bar",
			err:     status.Error(codes.Unimplemented, "foo"),

			expectLevel: quicklog.LevelWarning,
		},
		{
			name: "Internal",

			service: "foo",
			in:      42,
			out:     "bar",
			err:     status.Error(codes.Internal, "foo"),

			expectLevel: quicklog.LevelError,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			service := &fakeExecService{out: testCase.out, err: testCase.err}
			logger := quicklogmocks.NewMockLogger(t)

			logger.On("Log", testCase.expectLevel, mock.Anything).Once()

			res, err := arpc.WithReport(testCase.service, service, logger).Exec(context.Background(), testCase.in)
			require.Equal(t, testCase.out, res)
			require.ErrorIs(t, err, testCase.err)
		})
	}
}
