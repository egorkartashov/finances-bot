package remove_limit

import (
	"context"
)

type limitStorage interface {
	Delete(ctx context.Context, userID int64, category string) (err error)
}
