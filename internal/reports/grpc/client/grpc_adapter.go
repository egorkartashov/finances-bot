package client

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/grpcapi/send_report"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/reports"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GrpcAdapter struct {
	client send_report.ReportSenderClient
}

func NewGrpcAdapter(client send_report.ReportSenderClient) *GrpcAdapter {
	return &GrpcAdapter{client: client}
}

func (g *GrpcAdapter) Send(ctx context.Context, req *reports.FormattedReport) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "send-report-grpc")
	defer span.Finish()

	grpcReq := send_report.SendReportRequest{
		UserID:   req.UserID,
		Format:   string(req.Format),
		Currency: string(req.Currency),
		Period:   int32(req.Period),
		Payload:  req.Payload,
		Date:     timestamppb.New(req.Date),
	}
	_, err := g.client.SendReport(ctx, &grpcReq)
	return err
}
