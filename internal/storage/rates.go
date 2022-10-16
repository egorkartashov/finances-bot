package storage

import (
	"sync"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/currency"
)

type Rates struct {
	mu       *sync.RWMutex
	ratesMap map[currency.Currency]map[time.Time]decimal.Decimal
	nextID   int32
}

func NewRates() (*Rates, error) {
	r := &Rates{
		ratesMap: make(map[currency.Currency]map[time.Time]decimal.Decimal),
		nextID:   1,
		mu:       new(sync.RWMutex),
	}

	err := r.AddRate(currency.Rate{
		From:  currency.RUB,
		To:    currency.RUB,
		Value: decimal.NewFromInt32(1),
	})

	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *Rates) AddRates(rates []currency.Rate) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, rate := range rates {
		if err := r.addRateNoLock(rate); err != nil {
			return err
		}
	}
	return nil
}

func (r *Rates) AddRate(newRate currency.Rate) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.addRateNoLock(newRate)
}

func (r *Rates) addRateNoLock(rate currency.Rate) (err error) {
	err = nil
	_, ok := r.ratesMap[rate.From]
	if ok {
		r.ratesMap[rate.From][rate.Date] = rate.Value
	} else {
		r.ratesMap[rate.From] = map[time.Time]decimal.Decimal{rate.Date: rate.Value}
	}
	return
}

func (r *Rates) GetRate(from currency.Currency, date time.Time) (rate decimal.Decimal, ok bool, err error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	err = nil

	if _, ok = r.ratesMap[from]; !ok {
		ok = false
		return
	}

	rate, ok = r.ratesMap[from][date]
	return
}
