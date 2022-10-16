package currency

import (
	"github.com/shopspring/decimal"
	"time"
)

type (
	ratesProvider interface {
		GetRate(from, to Currency, date time.Time) (Rate, error)
	}
)

type Converter struct {
	ratesProvider ratesProvider
}

func NewConverter(ratesProvider ratesProvider) *Converter {
	return &Converter{ratesProvider}
}

func (c *Converter) Convert(sum decimal.Decimal, from, to Currency, date time.Time) (res decimal.Decimal, err error) {
	rate, err := c.ratesProvider.GetRate(from, to, date)
	if err != nil {
		return
	}

	res = sum.Mul(rate.Value)
	return
}
