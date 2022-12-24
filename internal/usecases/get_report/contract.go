//go:generate mockgen -source ${GOFILE} -package ${GOPACKAGE}_mocks -destination mocks/${GOPACKAGE}_mocks.go
package get_report

import (
	"context"

	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/entities"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/reports"
)

type reportCache interface {
	Get(ctx context.Context, req *reports.NewReportRequest) (*reports.FormattedReport, error)
	Save(ctx context.Context, report *reports.FormattedReport) error
}

type userStorage interface {
	Get(ctx context.Context, id int64) (entities.User, bool, error)
}

type reportRequester interface {
	Send(ctx context.Context, req *reports.NewReportRequest) error
}
