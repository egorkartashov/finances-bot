package entities

import (
	"time"

	"github.com/shopspring/decimal"
)

type MonthBudgetLimit struct {
	UserID   int64
	Category string
	Sum      decimal.Decimal
	SetAt    time.Time
	Currency Currency
}
