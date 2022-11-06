package set_limit

import (
	"context"

	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/usecases/set_limit"
)

type usecase interface {
	SetLimit(ctx context.Context, req set_limit.Req) (set_limit.Resp, error)
}
