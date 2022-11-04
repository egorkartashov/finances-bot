package currency

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/entities"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/users"
)

type Converter struct {
	cfg           cfg
	ratesProvider ratesProvider
	userUc        userUc
}

func NewConverter(cfg cfg, ratesProvider ratesProvider, userUc userUc) *Converter {
	return &Converter{
		cfg:           cfg,
		ratesProvider: ratesProvider,
		userUc:        userUc,
	}
}

func (c *Converter) Convert(
	ctx context.Context, sum decimal.Decimal, from, to entities.Currency, date time.Time,
) (res decimal.Decimal, rate entities.Rate, err error) {
	rate, err = c.ratesProvider.GetRate(ctx, from, to, date)
	if err != nil {
		return
	}

	res = sum.Mul(rate.Value)
	return
}

func (c *Converter) ConvertToBaseCurrency(
	ctx context.Context, sum decimal.Decimal, userID int64, date time.Time,
) (res decimal.Decimal, rate entities.Rate, err error) {
	user, ok, err := c.userUc.Get(ctx, userID)
	if err != nil {
		return
	}
	if !ok {
		err = users.NewUserNotFoundErr(userID)
		return
	}

	res, rate, err = c.Convert(ctx, sum, user.Currency, c.cfg.BaseCurrency(), date)
	return
}
