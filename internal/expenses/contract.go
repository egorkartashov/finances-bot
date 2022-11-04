//go:generate mockgen -source ${GOFILE} -package ${GOPACKAGE}_mocks -destination mocks/${GOPACKAGE}.go
package expenses

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/entities"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/limits"
)

type cfg interface {
	BaseCurrency() entities.Currency
}

type tx interface {
	WithTransaction(ctx context.Context, fn func(context.Context) error) error
}

type expenseStorage interface {
	AddExpense(ctx context.Context, userID int64, exp entities.Expense) error
	GetExpenses(ctx context.Context, userID int64, minTime time.Time) ([]entities.Expense, error)
}

type userStorage interface {
	Get(ctx context.Context, id int64) (entities.User, bool, error)
}

type currencyConverter interface {
	Convert(ctx context.Context, sum decimal.Decimal, from, to entities.Currency, date time.Time) (
		res decimal.Decimal, rate entities.Rate, err error,
	)
	ConvertToBaseCurrency(ctx context.Context, sum decimal.Decimal, userID int64, date time.Time) (
		res decimal.Decimal, rate entities.Rate, err error,
	)
}

type limitUc interface {
	CheckLimit(ctx context.Context, userID int64, expense entities.Expense) (limits.LimitCheckResult, error)
}
