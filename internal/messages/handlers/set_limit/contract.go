package set_limit

import (
	"context"

	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/entities"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/limits"
)

type limitUc interface {
	SetLimit(ctx context.Context, limit entities.MonthBudgetLimit) (limits.SetLimitResult, error)
}
