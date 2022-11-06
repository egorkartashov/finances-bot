package remove_limit

import "context"

type usecase interface {
	RemoveLimit(ctx context.Context, userID int64, category string) error
}
