package remove_limit

import "context"

type limitUc interface {
	RemoveLimit(ctx context.Context, userID int64, category string) error
}
