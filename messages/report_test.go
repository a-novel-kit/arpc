package arpcmessages_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	testutils "github.com/a-novel-kit/test-utils"

	arpcmessages "github.com/a-novel-kit/arpc/messages"
)

func TestReport(t *testing.T) {
	testCases := []struct {
		name string

		metrics *arpcmessages.Metrics
		service string
		err     error

		expectConsole string
		expectJSON    interface{}
	}{
		{
			name: "SimpleReport",

			metrics: nil,
			service: "MyService",
			err:     nil,

			expectConsole: "✅ OK [MyService]\n\n",
			expectJSON: map[string]interface{}{
				"grpcRequest": map[string]interface{}{
					"code":    codes.OK,
					"service": "MyService",
				},
				"severity": "INFO",
			},
		},
		{
			name: "Unavailable",

			metrics: nil,
			service: "MyService",
			err:     status.Error(codes.Unavailable, "uwups"),

			expectConsole: "⚠ Unavailable [MyService]\n  rpc error: code = Unavailable desc = uwups\n\n",
			expectJSON: map[string]interface{}{
				"grpcRequest": map[string]interface{}{
					"code":    codes.Unavailable,
					"service": "MyService",
				},
				"error":    "rpc error: code = Unavailable desc = uwups",
				"severity": "WARNING",
			},
		},
		{
			name: "Canceled",

			metrics: nil,
			service: "MyService",
			err:     status.Error(codes.Canceled, "uwups"),

			expectConsole: "⚠ Canceled [MyService]\n  rpc error: code = Canceled desc = uwups\n\n",
			expectJSON: map[string]interface{}{
				"grpcRequest": map[string]interface{}{
					"code":    codes.Canceled,
					"service": "MyService",
				},
				"error":    "rpc error: code = Canceled desc = uwups",
				"severity": "WARNING",
			},
		},
		{
			name: "Unimplemented",

			metrics: nil,
			service: "MyService",
			err:     status.Error(codes.Unimplemented, "uwups"),

			expectConsole: "⚠ Unimplemented [MyService]\n  rpc error: code = Unimplemented desc = uwups\n\n",
			expectJSON: map[string]interface{}{
				"grpcRequest": map[string]interface{}{
					"code":    codes.Unimplemented,
					"service": "MyService",
				},
				"error":    "rpc error: code = Unimplemented desc = uwups",
				"severity": "WARNING",
			},
		},
		{
			name: "Internal",

			metrics: nil,
			service: "MyService",
			err:     status.Error(codes.Internal, "uwups"),

			expectConsole: "👶🔪🩸 Internal [MyService]\n  rpc error: code = Internal desc = uwups\n\n",
			expectJSON: map[string]interface{}{
				"grpcRequest": map[string]interface{}{
					"code":    codes.Internal,
					"service": "MyService",
				},
				"error":    "rpc error: code = Internal desc = uwups",
				"severity": "ERROR",
			},
		},
		{
			name: "Unknown",

			metrics: nil,
			service: "MyService",
			err:     testutils.ErrDummy,

			expectConsole: "👶🔪🩸 Unknown [MyService]\n  uwups\n\n",
			expectJSON: map[string]interface{}{
				"grpcRequest": map[string]interface{}{
					"code":    codes.Unknown,
					"service": "MyService",
				},
				"error":    "uwups",
				"severity": "ERROR",
			},
		},
		{
			name: "OtherError",

			metrics: nil,
			service: "MyService",
			err:     status.Error(codes.NotFound, "uwups"),

			expectConsole: "❌ NotFound [MyService]\n  rpc error: code = NotFound desc = uwups\n\n",
			expectJSON: map[string]interface{}{
				"grpcRequest": map[string]interface{}{
					"code":    codes.NotFound,
					"service": "MyService",
				},
				"error":    "rpc error: code = NotFound desc = uwups",
				"severity": "ERROR",
			},
		},
		{
			name: "WithMetrics",

			metrics: &arpcmessages.Metrics{
				Latency: 2*time.Second + 200*time.Millisecond,
			},
			service: "MyService",
			err:     nil,

			expectConsole: "✅ OK [MyService] (2.2s)\n\n",
			expectJSON: map[string]interface{}{
				"grpcRequest": map[string]interface{}{
					"code":    codes.OK,
					"service": "MyService",
					"latency": 2*time.Second + 200*time.Millisecond,
				},
				"severity": "INFO",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			message := arpcmessages.NewReport(testCase.metrics, testCase.service, testCase.err)
			require.Equal(t, testCase.expectConsole, message.RenderTerminal())
			require.Equal(t, testCase.expectJSON, message.RenderJSON())
		})
	}
}
