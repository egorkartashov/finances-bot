package dtos

import (
	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/entities"
)

type SumWithCurrency struct {
	Sum      decimal.Decimal
	Currency entities.Currency
}
