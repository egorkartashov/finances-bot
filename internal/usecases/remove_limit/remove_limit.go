package remove_limit

import (
	"context"

	"github.com/pkg/errors"
)

type Usecase struct {
	limitStorage limitStorage
}

func NewUsecase(ls limitStorage) *Usecase {
	return &Usecase{
		limitStorage: ls,
	}
}

func (u *Usecase) RemoveLimit(ctx context.Context, userID int64, category string) error {
	err := u.limitStorage.Delete(ctx, userID, category)
	if err != nil {
		return errors.WithMessage(err, "RemoveLimit")
	}
	return nil
}
