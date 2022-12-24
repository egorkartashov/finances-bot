package limits

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/entities"
)

type Checker struct {
	limitStorage   limitStorage
	expenseStorage expenseStorage
	userStorage    userStorage
	converter      currencyConverter
}

func NewChecker(ls limitStorage, es expenseStorage, us userStorage, cc currencyConverter) *Checker {
	return &Checker{
		limitStorage:   ls,
		expenseStorage: es,
		userStorage:    us,
		converter:      cc,
	}
}

type LimitCheckResult struct {
	Status            LimitCheckStatus
	Limit             decimal.Decimal
	SumWithNewExpense decimal.Decimal
	Currency          entities.Currency
}

type LimitCheckStatus byte

const (
	StatusLimitNotSet LimitCheckStatus = iota
	StatusLimitSatisfied
	StatusLimitExceeded
)

func (u *Checker) Check(ctx context.Context, userID int64, expense entities.Expense) (
	res LimitCheckResult, err error,
) {
	limit, ok, err := u.limitStorage.Get(ctx, userID, expense.Category)
	if err != nil {
		return
	}
	if !ok {
		res.Status = StatusLimitNotSet
		return
	}

	start, end := getMonthStartAndEndDate(expense.Date)
	expensesSum, err := u.expenseStorage.GetSumForCategoryAndPeriod(ctx, userID, expense.Category, start, end)
	if err != nil {
		return
	}

	sumWithNewExpense := expensesSum.Add(expense.Sum)

	res = LimitCheckResult{
		Limit:             limit.Sum,
		SumWithNewExpense: sumWithNewExpense,
		Currency:          limit.Currency,
	}

	if sumWithNewExpense.LessThanOrEqual(limit.Sum) {
		res.Status = StatusLimitSatisfied
	} else {
		res.Status = StatusLimitExceeded
	}
	return
}

func getMonthStartAndEndDate(expenseDate time.Time) (start, end time.Time) {
	year := expenseDate.Year()
	month := expenseDate.Month()
	start = time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	end = start.AddDate(0, 1, 0)
	return
}
