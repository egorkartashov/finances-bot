package entities

import (
	"time"

	"github.com/shopspring/decimal"
)

type Rate struct {
	From  Currency
	To    Currency
	Value decimal.Decimal
	Date  time.Time
}

func (r Rate) ReverseRate() Rate {
	return Rate{
		From:  r.To,
		To:    r.From,
		Date:  r.Date,
		Value: decimal.NewFromInt32(1).Div(r.Value),
	}
}
