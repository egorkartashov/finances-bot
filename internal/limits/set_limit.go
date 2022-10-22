package limits

import (
	"context"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/entities"
)

type SetLimitDto struct {
	UserID            int64
	SumInUserCurrency decimal.Decimal
	Category          string
}

type SetLimitResult struct {
	Limit             entities.MonthBudgetLimit
	SumInUserCurrency decimal.Decimal
	ExchangeRate      entities.Rate
}

func (u *Usecase) SetLimit(ctx context.Context, limit entities.MonthBudgetLimit) (result SetLimitResult, err error) {
	if err = u.limitStorage.Save(ctx, limit); err != nil {
		err = errors.WithMessage(err, "SetLimit")
		return
	}

	result = SetLimitResult{
		SumInUserCurrency: limit.Sum,
	}

	convertedSum, rate, err := u.converter.ConvertToBaseCurrency(ctx, limit.Sum, limit.UserID, limit.SetAt)
	if err != nil {
		return
	}

	limit.Sum = convertedSum
	if err = u.limitStorage.Save(ctx, limit); err != nil {
		return
	}

	result.Limit = limit
	result.ExchangeRate = rate
	return
}
