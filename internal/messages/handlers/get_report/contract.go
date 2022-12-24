package get_report

import (
	"context"

	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/reports"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/usecases/get_report"
)

type usecase interface {
	GenerateReport(ctx context.Context, req get_report.ReportRequest) (
		get_report.ReportResponse, error,
	)
	ReportFinished(ctx context.Context, report *reports.FormattedReport) error
}
