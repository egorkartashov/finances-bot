package users

import (
	"context"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/entities"
)

type userStorage interface {
	Save(ctx context.Context, user entities.User) error
	Get(ctx context.Context, id int64) (u entities.User, ok bool, err error)
}

type cfg interface {
	BaseCurrency() entities.Currency
}

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

func (u *Usecase) Register(ctx context.Context, userID int64) (err error) {
	_, ok, err := u.storage.Get(ctx, userID)
	if err != nil {
		err = errors.WithMessage(err, "Register")
		return
	}
	if ok {
		return
	}

	user := entities.User{
		ID:       userID,
		Currency: u.cfg.BaseCurrency(),
	}
	err = u.storage.Save(ctx, user)
	if err != nil {
		err = errors.WithMessage(err, "Register")
		return
	}
	return
}

func (u *Usecase) SetCurrency(ctx context.Context, userID int64, curr entities.Currency) error {
	user, ok, err := u.storage.Get(ctx, userID)
	if err != nil {
		return err
	}
	if !ok {
		return NewUserNotFoundErr(userID)
	}

	user.Currency = curr
	return u.storage.Save(ctx, user)
}
