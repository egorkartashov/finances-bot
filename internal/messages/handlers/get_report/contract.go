package get_report

import (
	"context"

	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/entities"
)

type usecase interface {
	GenerateReport(ctx context.Context, userID int64, period entities.ReportPeriod) (*entities.Report, error)
}
