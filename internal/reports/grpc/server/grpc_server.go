package server

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/entities"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/grpcapi/send_report"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/reports"
)

type reportMessageSender interface {
	SendFinishedReport(ctx context.Context, report *reports.FormattedReport) error
}

type GrpcServer struct {
	send_report.UnimplementedReportSenderServer
	reportMessageSender reportMessageSender
}

func NewGrpc(rms reportMessageSender) *GrpcServer {
	return &GrpcServer{
		reportMessageSender: rms,
	}
}

func (g *GrpcServer) SendReport(ctx context.Context, req *send_report.SendReportRequest) (
	*send_report.SendReportResponse, error,
) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "receive-report-grpc")
	defer span.Finish()

	report := reports.FormattedReport{
		NewReportRequest: reports.NewReportRequest{
			UserID:   req.UserID,
			Currency: entities.Currency(req.Currency),
			Format:   entities.ReportFormat(req.Format),
			Period:   entities.ReportPeriod(req.Period),
			Date:     req.Date.AsTime(),
		},
		Payload: req.Payload,
	}
	err := g.reportMessageSender.SendFinishedReport(ctx, &report)
	return &send_report.SendReportResponse{}, err
}
