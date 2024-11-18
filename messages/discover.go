package arpcmessages

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
	"google.golang.org/grpc"

	"github.com/a-novel-kit/quicklog"
	"github.com/a-novel-kit/quicklog/messages"
)

type discoverMessage struct {
	services []grpc.ServiceDesc
	port     int

	quicklog.Message
}

func (discover *discoverMessage) RenderTerminal() string {
	servicesList := list.New().
		Enumerator(func(_ list.Items, _ int) string {
			return ""
		}).
		Indenter(func(_ list.Items, _ int) string { return "" }).
		EnumeratorStyle(lipgloss.NewStyle().MarginLeft(4)).
		ItemStyle(lipgloss.NewStyle().MarginLeft(1).Foreground(lipgloss.Color("220")))

	for _, service := range discover.services {
		var methodsItems []interface{}

		methods := list.New().
			Enumerator(func(_ list.Items, _ int) string {
				return ""
			}).
			EnumeratorStyle(lipgloss.NewStyle().MarginLeft(4)).
			Indenter(func(_ list.Items, _ int) string { return "" }).
			ItemStyle(lipgloss.NewStyle().Faint(true).Foreground(lipgloss.Color("220")))

		for _, method := range service.Methods {
			methodsItems = append(methodsItems, method.MethodName)
		}
		for _, method := range service.Streams {
			methodsItems = append(methodsItems, "["+method.StreamName+"]")
		}

		methods.Items(methodsItems...)
		servicesList.Items(service.ServiceName, methods)
	}

	description := fmt.Sprintf(
		"%v services registered on port :%v",
		len(discover.services), discover.port,
	)

	return messages.NewTitle("RPC services successfully registered.", description, nil).RenderTerminal() +
		servicesList.String() + "\n"
}

func (discover *discoverMessage) RenderJSON() map[string]interface{} {
	servicesList := map[string]interface{}{}

	for _, service := range discover.services {
		methods := make([]interface{}, 0)
		streams := make([]interface{}, 0)

		for _, method := range service.Methods {
			methods = append(methods, method.MethodName)
		}
		for _, stream := range service.Streams {
			streams = append(streams, stream.StreamName)
		}

		servicesList[service.ServiceName] = map[string]interface{}{
			"methods": methods,
			"streams": streams,
		}
	}

	return servicesList
}

func NewDiscover(services []grpc.ServiceDesc, port int) quicklog.Message {
	return &discoverMessage{
		services: services,
		port:     port,
	}
}
