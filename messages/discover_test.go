package arpcmessages_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	arpcmessages "github.com/a-novel-kit/arpc/messages"
)

func TestDiscoverGRPC(t *testing.T) {
	t.Run("Render", func(t *testing.T) {
		content := arpcmessages.NewDiscover([]grpc.ServiceDesc{
			{
				ServiceName: "Service1",
				Methods: []grpc.MethodDesc{
					{
						MethodName: "Method1",
					},
					{
						MethodName: "Method2",
					},
				},
				Streams: []grpc.StreamDesc{
					{
						StreamName: "Stream1",
					},
					{
						StreamName: "Stream2",
					},
				},
			},
			{
				ServiceName: "Service2",
				Methods: []grpc.MethodDesc{
					{
						MethodName: "Method1",
					},
				},
			},
			{
				ServiceName: "Service3",
				Methods: []grpc.MethodDesc{
					{
						MethodName: "Method1",
					},
					{
						MethodName: "Method2",
					},
				},
			},
		}, 1234)

		expectConsole := "╭────────────────────────────────────────────────────────────────────────────────╮\n" +
			"│ RPC services successfully registered.                                          │\n" +
			"│ 3 services registered on port :1234                                            │\n" +
			"╰────────────────────────────────────────────────────────────────────────────────╯\n" +
			"     Service1\n" +
			"        Method1\n" +
			"        Method2\n" +
			"        [Stream1]\n" +
			"        [Stream2]\n" +
			"     Service2\n" +
			"        Method1\n" +
			"     Service3\n" +
			"        Method1\n" +
			"        Method2\n"
		expectJSON := map[string]interface{}{
			"Service1": map[string]interface{}{
				"methods": []interface{}{"Method1", "Method2"},
				"streams": []interface{}{"Stream1", "Stream2"},
			},
			"Service2": map[string]interface{}{
				"methods": []interface{}{"Method1"},
				"streams": []interface{}{},
			},
			"Service3": map[string]interface{}{
				"methods": []interface{}{"Method1", "Method2"},
				"streams": []interface{}{},
			},
		}

		require.Equal(t, expectConsole, content.RenderTerminal())
		require.Equal(t, expectJSON, content.RenderJSON())
	})
}
