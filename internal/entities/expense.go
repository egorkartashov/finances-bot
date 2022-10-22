package entities

import (
	"time"

	"github.com/shopspring/decimal"
)

type Expense struct {
	Category string
	Sum      decimal.Decimal
	Date     time.Time
}
