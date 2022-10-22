package limits

import "github.com/shopspring/decimal"

type LimitCheckResult struct {
	Status                    LimitCheckStatus
	Limit                     decimal.Decimal
	TotalSumWithoutNewExpense decimal.Decimal
	TotalSumWithNewExpense    decimal.Decimal
}

type LimitCheckStatus byte

const (
	StatusLimitNotSet    LimitCheckStatus = 1
	StatusLimitSatisfied LimitCheckStatus = 2
	StatusLimitExceeded                   = 3
)