//go:generate mockgen -source ${GOFILE} -package ${GOPACKAGE}_mocks -destination mocks/${GOPACKAGE}_mocks.go
package get_expenses_report

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/entities"
)

type expenseStorage interface {
	GetExpenses(ctx context.Context, userID int64, startDate time.Time) ([]entities.Expense, error)
}

type userStorage interface {
	Get(ctx context.Context, id int64) (entities.User, bool, error)
}

type currencyConverter interface {
	FromBase(ctx context.Context, to entities.Currency, sum decimal.Decimal, date time.Time) (
		res decimal.Decimal, err error,
	)
	ToBase(ctx context.Context, from entities.Currency, sum decimal.Decimal, date time.Time) (
		res decimal.Decimal, curr entities.Currency, err error,
	)
}

type reportCache interface {
	Get(ctx context.Context, userID int64, period entities.ReportPeriod) (*entities.Report, error)
	Save(ctx context.Context, userID int64, period entities.ReportPeriod, report *entities.Report) error
}
