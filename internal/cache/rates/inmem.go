package rates

import (
	"context"
	"fmt"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/cache/lru"
	"time"

	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/entities"
)

type ratesProvider interface {
	GetRate(ctx context.Context, from, to entities.Currency, date time.Time) (entities.Rate, error)
}

type InMem struct {
	cache    *lru.Lru
	provider ratesProvider
}

func NewInMemCacheDecorator(rp ratesProvider) *InMem {
	return &InMem{
		cache:    lru.New(100),
		provider: rp,
	}
}

func (i *InMem) GetRate(ctx context.Context, from, to entities.Currency, date time.Time) (entities.Rate, error) {
	key := getKey(from, to, date)
	if item, ok := i.cache.Get(key); ok {
		rate := item.(*entities.Rate)
		return *rate, nil
	}

	rate, err := i.provider.GetRate(ctx, from, to, date)
	if err != nil {
		return entities.Rate{}, err
	}

	i.cache.Set(key, &rate)
	return rate, nil
}

func getKey(from, to entities.Currency, date time.Time) string {
	dateIso := date.Format("2006-01-02")
	return fmt.Sprintf("%s_%s_%s", from, to, dateIso)
}
