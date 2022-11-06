package start

import "context"

type userUc interface {
	Register(ctx context.Context, userID int64) error
}
