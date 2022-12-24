package generator

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/entities"
)

type currencyConverter interface {
	FromBase(ctx context.Context, to entities.Currency, sum decimal.Decimal, date time.Time) (
		res decimal.Decimal, err error,
	)
	ToBase(ctx context.Context, from entities.Currency, sum decimal.Decimal, date time.Time) (
		res decimal.Decimal, curr entities.Currency, err error,
	)
}

type expenseStorage interface {
	GetExpenses(ctx context.Context, userID int64, startDate time.Time) ([]entities.Expense, error)
}

type ReportPresenter interface {
	Format() entities.ReportFormat
	Present(ctx context.Context, report *entities.Report) (payload string, err error)
}
