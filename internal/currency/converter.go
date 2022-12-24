package currency

import (
	"context"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/entities"
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
) (res decimal.Decimal, curr entities.Currency, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "convert-currency")
	defer span.Finish()
	span.SetTag("from", from)
	span.SetTag("to", to)
	span.SetTag("date", date)

	rate, err := c.ratesProvider.GetRate(ctx, from, to, date)
	if err != nil {
		return
	}

	curr = rate.To
	res = sum.Mul(rate.Value)
	return
}

func (c *Converter) FromBase(
	ctx context.Context, to entities.Currency, sum decimal.Decimal, date time.Time,
) (res decimal.Decimal, err error) {
	res, _, err = c.Convert(ctx, sum, c.cfg.BaseCurrency(), to, date)
	return
}

func (c *Converter) ToBase(
	ctx context.Context, from entities.Currency, sum decimal.Decimal, date time.Time,
) (res decimal.Decimal, curr entities.Currency, err error) {
	return c.Convert(ctx, sum, from, c.cfg.BaseCurrency(), date)
}
