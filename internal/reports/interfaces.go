package reports

import (
	"context"
	"time"

	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/entities"
)

type ReportRequester interface {
	Send(ctx context.Context, req *NewReportRequest) error
}

type ReportGenerator interface {
	Generate(ctx context.Context, req *NewReportRequest) (*GeneratedReportResponse, error)
}

type FinishedReportSender interface {
	Send(ctx context.Context, req *FormattedReport) error
}

type NewReportRequest struct {
	UserID   int64
	Currency entities.Currency
	Format   entities.ReportFormat
	Period   entities.ReportPeriod
	Date     time.Time
}

type GeneratedReportResponse struct {
	Payload string
}

type FormattedReport struct {
	NewReportRequest
	Payload string
}
