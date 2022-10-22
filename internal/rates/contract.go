package rates

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/entities"
)

type ratesToRubApi interface {
	FetchRatesToRub(ctx context.Context, currencies []entities.Currency, at time.Time) ([]entities.Rate, error)
}

type ratesStorage interface {
	AddRates(ctx context.Context, rates []entities.Rate) error
	GetRate(ctx context.Context, from entities.Currency, date time.Time) (r decimal.Decimal, ok bool, err error)
}

type cfg interface {
	BaseCurrency() entities.Currency
}
