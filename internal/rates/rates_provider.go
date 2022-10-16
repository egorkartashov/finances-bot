package rates

import (
	"context"
	"github.com/shopspring/decimal"
	"log"
	"time"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/currency"
)

type ratesToRubApi interface {
	FetchRatesToRub(currencies []currency.Currency, at time.Time) ([]currency.Rate, error)
}

type ratesStorage interface {
	AddRates(rates []currency.Rate) error
	GetRate(from currency.Currency, date time.Time) (r decimal.Decimal, ok bool, err error)
}

type Provider struct {
	l        *log.Logger
	baseCurr currency.Currency
	storage  ratesStorage
	api      ratesToRubApi
}

func NewProvider(l *log.Logger, baseCurr currency.Currency, api ratesToRubApi, storage ratesStorage) *Provider {
	return &Provider{
		l:        l,
		baseCurr: baseCurr,
		storage:  storage,
		api:      api,
	}
}

func (p *Provider) GetRate(from, to currency.Currency, date time.Time) (r currency.Rate, err error) {
	if from == to {
		r = rateToSelf(from, date)
		return
	}

	if from != p.baseCurr && to != p.baseCurr {
		err = errors.New("from and to are equal to base curr")
		return
	}

	if to == p.baseCurr {
		r, err = p.getRate(from, date)
		return
	}

	if r, err = p.getRate(to, date); err != nil {
		return
	}

	r = r.ReverseRate()
	return
}

func (p *Provider) getRate(from currency.Currency, date time.Time) (r currency.Rate, err error) {
	rateVal, ok, err := p.storage.GetRate(from, date)
	if err != nil {
		return
	}
	if ok {
		r = currency.Rate{
			From:  from,
			To:    p.baseCurr,
			Value: rateVal,
			Date:  date,
		}
		return
	}

	var rates []currency.Rate
	rates, err = p.updateRates([]currency.Currency{from}, date)
	if err != nil {
		return
	}

	r = rates[0]
	return
}

func (p *Provider) UpdateRates(ctx context.Context, freq time.Duration, currencies []currency.Currency) {
	p.updateRatesForNow(currencies)

	ticker := time.NewTicker(freq)
	select {
	case <-ticker.C:
		p.updateRatesForNow(currencies)
	case <-ctx.Done():
		log.Println("UpdateRates: exiting due to ctx.Done()")
		ticker.Stop()
		return
	}
}

func (p *Provider) updateRatesForNow(currencies []currency.Currency) {
	_, err := p.updateRates(currencies, time.Now())
	if err != nil {
		p.l.Println(err)
	}
}

func (p *Provider) updateRates(currencies []currency.Currency, date time.Time) ([]currency.Rate, error) {
	rates, err := p.api.FetchRatesToRub(currencies, date)
	if err != nil {
		err = errors.WithMessage(err, "error in updateRates")
		return nil, err
	}

	err = p.storage.AddRates(rates)
	if err != nil {
		err = errors.WithMessage(err, "error in updateRates")
		return nil, err
	}

	return rates, nil
}

func rateToSelf(cur currency.Currency, date time.Time) currency.Rate {
	return currency.Rate{
		From:  cur,
		To:    cur,
		Value: decimal.NewFromInt32(1),
		Date:  date,
	}
}
