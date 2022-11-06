package dtos

import (
	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/entities"
)

type SumWithCurrency struct {
	Sum      decimal.Decimal
	Currency entities.Currency
}
