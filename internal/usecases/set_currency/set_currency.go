package set_currency

import (
	"context"

	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/entities"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/users"
)

type Usecase struct {
	cfg     cfg
	storage userStorage
}

func NewUsecase(cfg cfg, storage userStorage) *Usecase {
	return &Usecase{
		cfg:     cfg,
		storage: storage,
	}
}

func (u *Usecase) SetCurrency(ctx context.Context, userID int64, curr entities.Currency) error {
	user, ok, err := u.storage.Get(ctx, userID)
	if err != nil {
		return err
	}
	if !ok {
		return users.NewUserNotFoundErr(userID)
	}

	user.Currency = curr
	return u.storage.Save(ctx, user)
}
