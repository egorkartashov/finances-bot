package set_limit

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/dtos"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/entities"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/usecases/set_currency"
)

type Usecase struct {
	limitStorage limitStorage
	userStorage  userStorage
	converter    currencyConverter
}

func NewUsecase(ls limitStorage, us userStorage, cc currencyConverter) *Usecase {
	return &Usecase{
		limitStorage: ls,
		userStorage:  us,
		converter:    cc,
	}
}

type Req struct {
	UserID            int64
	SumInUserCurrency decimal.Decimal
	Category          string
	SetAt             time.Time
}

type Resp struct {
	InputSum dtos.SumWithCurrency
	SavedSum dtos.SumWithCurrency
	Category string
}

func (u *Usecase) SetLimit(ctx context.Context, req Req) (resp Resp, err error) {
	user, ok, err := u.userStorage.Get(ctx, req.UserID)
	if err != nil {
		return Resp{}, err
	}
	if !ok {
		return Resp{}, set_currency.NewUserNotFoundErr(req.UserID)
	}

	sumInBaseCurr, baseCurr, err := u.converter.ToBase(ctx, user.Currency, req.SumInUserCurrency, req.SetAt)
	if err != nil {
		return Resp{}, err
	}

	limit := entities.MonthBudgetLimit{
		UserID:   req.UserID,
		Category: req.Category,
		Sum:      sumInBaseCurr,
		SetAt:    req.SetAt,
		Currency: baseCurr,
	}

	if err = u.limitStorage.Save(ctx, limit); err != nil {
		return Resp{}, err
	}

	resp = Resp{
		InputSum: dtos.SumWithCurrency{
			Sum:      req.SumInUserCurrency,
			Currency: user.Currency,
		},
		SavedSum: dtos.SumWithCurrency{
			Sum:      sumInBaseCurr,
			Currency: baseCurr,
		},
		Category: req.Category,
	}

	return resp, nil
}
