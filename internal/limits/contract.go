package limits

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/entities"
)

type limitStorage interface {
	Save(ctx context.Context, limit entities.MonthBudgetLimit) (err error)
	Get(ctx context.Context, userID int64, category string) (limit entities.MonthBudgetLimit, ok bool, err error)
	Delete(ctx context.Context, userID int64, category string) (err error)
}

type expenseStorage interface {
	GetSumForCategoryAndPeriod(
		ctx context.Context, userID int64, category string, startDate, endDate time.Time,
	) (decimal.Decimal, error)
}

type converter interface {
	ConvertToBaseCurrency(ctx context.Context, sum decimal.Decimal, userID int64, date time.Time) (
		res decimal.Decimal, rate entities.Rate, err error,
	)
}
