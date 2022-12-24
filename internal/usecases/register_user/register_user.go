package register_user

import (
	"context"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/entities"
)

type Usecase struct {
	cfg     cfg
	storage userStorage
}

func NewUsecase(c cfg, s userStorage) *Usecase {
	return &Usecase{
		cfg:     c,
		storage: s,
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
