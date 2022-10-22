package expenses

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/entities"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/limits"
)

type AddExpenseDto struct {
	Category string
	Sum      int32
	Date     time.Time
}

type AddExpenseResult struct {
	SumInUserCurrency int32
	Expense           entities.Expense
	Rate              entities.Rate
	LimitCheckResult  limits.LimitCheckResult
}

func (u *Usecase) AddExpense(outsideCtx context.Context, userID int64, expDto AddExpenseDto) (AddExpenseResult, error) {
	var res AddExpenseResult
	txErr := u.tx.WithTransaction(
		outsideCtx, func(ctx context.Context) error {
			expense, rate, err := u.convertToExpense(ctx, expDto, userID)
			if err != nil {
				return err
			}

			// Пока просто сохраняем статус проверки лимита, т.к. запрещать вводить трату не будем.
			// Если лимит превышен, просто сообщим об этом пользователю
			limitStatus, err := u.limitUc.CheckLimit(ctx, userID, expense)
			if err != nil {
				return err
			}

			if err = u.expenseStorage.AddExpense(ctx, userID, expense); err != nil {
				return err
			}

			res = AddExpenseResult{
				SumInUserCurrency: expDto.Sum,
				Expense:           expense,
				Rate:              rate,
				LimitCheckResult:  limitStatus,
			}
			return nil
		},
	)

	return res, txErr
}

func (u *Usecase) convertToExpense(
	ctx context.Context, expDto AddExpenseDto, userID int64,
) (exp entities.Expense, rate entities.Rate, err error) {
	sum := decimal.NewFromInt32(expDto.Sum)
	subRub, rate, err := u.currencyConverter.ConvertToBaseCurrency(ctx, sum, userID, expDto.Date)
	if err != nil {
		return
	}

	exp = entities.Expense{
		Category: expDto.Category,
		Sum:      subRub,
		Date:     expDto.Date,
	}
	return
}
