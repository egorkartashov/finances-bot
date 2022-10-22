package limits

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/entities"
)

type Usecase struct {
	limitStorage   limitStorage
	expenseStorage expenseStorage
	converter      converter
}

func NewUsecase(limitStorage limitStorage, expenseStorage expenseStorage, converter converter) *Usecase {
	return &Usecase{
		limitStorage:   limitStorage,
		expenseStorage: expenseStorage,
		converter:      converter,
	}
}

func (u *Usecase) CheckLimit(ctx context.Context, userID int64, expense entities.Expense) (
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

	start, end := getStartAndEndDate(expense.Date)
	expensesSum, err := u.expenseStorage.GetSumForCategoryAndPeriod(ctx, userID, expense.Category, start, end)
	if err != nil {
		return
	}

	sumWithNewExpense := expensesSum.Add(expense.Sum)

	res.Limit = limit.Sum
	res.TotalSumWithoutNewExpense = expensesSum
	res.TotalSumWithNewExpense = sumWithNewExpense

	if sumWithNewExpense.LessThanOrEqual(limit.Sum) {
		res.Status = StatusLimitSatisfied
	} else {
		res.Status = StatusLimitExceeded
	}
	return
}

func (u *Usecase) RemoveLimit(ctx context.Context, userID int64, category string) error {
	err := u.limitStorage.Delete(ctx, userID, category)
	if err != nil {
		return errors.WithMessage(err, "RemoveLimit")
	}
	return nil
}

func getStartAndEndDate(expenseDate time.Time) (start, end time.Time) {
	year := expenseDate.Year()
	month := expenseDate.Month()
	start = time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	end = start.AddDate(0, 1, 0)
	return
}
