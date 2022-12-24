package limits

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/entities"
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

type userStorage interface {
	Get(ctx context.Context, id int64) (entities.User, bool, error)
}

type currencyConverter interface {
	ToBase(ctx context.Context, from entities.Currency, sum decimal.Decimal, date time.Time) (
		res decimal.Decimal, curr entities.Currency, err error,
	)
}
