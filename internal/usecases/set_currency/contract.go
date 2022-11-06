package set_currency

import (
	"context"

	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/entities"
)

type userStorage interface {
	Save(ctx context.Context, user entities.User) error
	Get(ctx context.Context, id int64) (u entities.User, ok bool, err error)
}

type cfg interface {
	BaseCurrency() entities.Currency
}
