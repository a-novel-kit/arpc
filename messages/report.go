package arpcmessages

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/a-novel-kit/quicklog"
)

type Metrics struct {
	Latency time.Duration
}

type reportMessage struct {
	metrics *Metrics
	service string
	err     error

	quicklog.Message
}

func (report *reportMessage) RenderTerminal() string {
	errorMessage := ""
	if report.err != nil {
		errorMessage = "\n" + lipgloss.NewStyle().MarginLeft(2).Foreground(lipgloss.Color("9")).
			Render(report.err.Error())
	}

	latencyMessage := ""
	if report.metrics != nil {
		latencyMessage = lipgloss.NewStyle().Faint(true).Render(fmt.Sprintf(" (%s)", report.metrics.Latency))
	}

	code := status.Code(report.err)

	color := lo.Switch[codes.Code, lipgloss.Color](code).
		Case(codes.OK, "33").
		Case(codes.Unavailable, "202").
		Case(codes.Canceled, "202").
		Case(codes.Unimplemented, "202").
		Default("9")

	prefix := lo.Switch[codes.Code, string](code).
		Case(codes.OK, "✅ ").
		Case(codes.Unavailable, "⚠ ").
		Case(codes.Canceled, "⚠ ").
		Case(codes.Unimplemented, "⚠ ").
		Case(codes.Internal, "👶🔪🩸 ").
		Case(codes.Unknown, "👶🔪🩸 ").
		Default("❌ ")

	return lipgloss.NewStyle().
		Foreground(color).
		Bold(true).
		Render(prefix+code.String()) +
		lipgloss.NewStyle().Foreground(color).Render(fmt.Sprintf(" [%s]", report.service)) +
		latencyMessage +
		errorMessage +
		"\n\n"
}

func (report *reportMessage) RenderJSON() map[string]interface{} {
	code := status.Code(report.err)

	severity := lo.Switch[codes.Code, string](code).
		Case(codes.OK, "INFO").
		Case(codes.Unavailable, "WARNING").
		Case(codes.Canceled, "WARNING").
		Case(codes.Unimplemented, "WARNING").
		Default("ERROR")

	// TODO: check if we can add trace to GRPC requests.
	// TODO: improve formatting of GRPC messages.
	grpcRequest := map[string]interface{}{
		"service": report.service,
		"code":    code,
	}

	output := map[string]interface{}{
		"severity":    severity,
		"grpcRequest": grpcRequest,
	}

	if report.metrics != nil {
		grpcRequest["latency"] = report.metrics.Latency
	}

	if report.err != nil {
		output["error"] = report.err.Error()
	}

	return output
}

// NewReport creates a new report message.
func NewReport(metrics *Metrics, service string, err error) quicklog.Message {
	return &reportMessage{
		metrics: metrics,
		service: service,
		err:     err,
	}
}
