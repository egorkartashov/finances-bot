package set_limit

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/entities"
)

type limitStorage interface {
	Save(ctx context.Context, limit entities.MonthBudgetLimit) (err error)
}

type userStorage interface {
	Get(ctx context.Context, id int64) (entities.User, bool, error)
}

type currencyConverter interface {
	ToBase(ctx context.Context, from entities.Currency, sum decimal.Decimal, date time.Time) (
		res decimal.Decimal, curr entities.Currency, err error,
	)
}
