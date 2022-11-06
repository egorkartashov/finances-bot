//go:generate mockgen -source ${GOFILE} -package ${GOPACKAGE}_mocks -destination mocks/${GOPACKAGE}_mocks.go
package add_expense

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/entities"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/limits"
)

type tx interface {
	WithTransaction(ctx context.Context, fn func(context.Context) error) error
}

type expenseStorage interface {
	AddExpense(ctx context.Context, userID int64, exp entities.Expense) error
}

type userStorage interface {
	Get(ctx context.Context, id int64) (entities.User, bool, error)
}

type currencyConverter interface {
	ToBase(ctx context.Context, from entities.Currency, sum decimal.Decimal, date time.Time) (
		res decimal.Decimal, curr entities.Currency, err error,
	)
}

type limitUc interface {
	Check(ctx context.Context, userID int64, expense entities.Expense) (limits.LimitCheckResult, error)
}
