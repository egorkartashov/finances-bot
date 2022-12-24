package currency

import (
	"context"
	"time"

	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/entities"
)

type (
	cfg interface {
		BaseCurrency() entities.Currency
	}
	ratesProvider interface {
		GetRate(ctx context.Context, from, to entities.Currency, date time.Time) (entities.Rate, error)
	}
	userUc interface {
		Get(ctx context.Context, userID int64) (entities.User, bool, error)
	}
)
