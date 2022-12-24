package rates

import (
	"context"
	"log"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/entities"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/logger"
	"go.uber.org/zap"
)

type Provider struct {
	cfg     cfg
	storage ratesStorage
	api     ratesToRubApi
}

func NewProvider(cfg cfg, api ratesToRubApi, storage ratesStorage) *Provider {
	return &Provider{
		cfg:     cfg,
		storage: storage,
		api:     api,
	}
}

func (p *Provider) GetRate(ctx context.Context, from, to entities.Currency, date time.Time) (
	r entities.Rate, err error,
) {
	if from == to {
		r = rateToSelf(from, date)
		return
	}

	baseCurr := p.cfg.BaseCurrency()
	if from != baseCurr && to != baseCurr {
		err = errors.New("neither from nor to are equal to base curr")
		return
	}

	if to == baseCurr {
		r, err = p.getRate(ctx, from, date)
		return
	}

	if r, err = p.getRate(ctx, to, date); err != nil {
		return
	}

	r = r.ReverseRate()
	return
}

func (p *Provider) getRate(ctx context.Context, from entities.Currency, date time.Time) (r entities.Rate, err error) {
	rateVal, ok, err := p.storage.GetRate(ctx, from, date)
	if err != nil {
		return
	}
	if ok {
		r = entities.Rate{
			From:  from,
			To:    p.cfg.BaseCurrency(),
			Value: rateVal,
			Date:  date,
		}
		return
	}

	var rates []entities.Rate
	rates, err = p.updateRates(ctx, []entities.Currency{from}, date)
	if err != nil {
		return
	}

	r = rates[0]
	return
}

func (p *Provider) UpdateRates(ctx context.Context, freq time.Duration, currencies []entities.Currency) {
	p.updateRatesForNow(ctx, currencies)

	ticker := time.NewTicker(freq)
	select {
	case <-ticker.C:
		p.updateRatesForNow(ctx, currencies)
	case <-ctx.Done():
		log.Println("UpdateRates: exiting due to ctx.Done()")
		ticker.Stop()
		return
	}
}

func (p *Provider) updateRatesForNow(ctx context.Context, currencies []entities.Currency) {
	_, err := p.updateRates(ctx, currencies, time.Now())
	if err != nil {
		logger.Error("error updating current rates", zap.Error(err))
	}
}

func (p *Provider) updateRates(ctx context.Context, currencies []entities.Currency, date time.Time) (
	[]entities.Rate, error,
) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "updateRates")
	defer span.Finish()

	span.SetTag("date", date)
	span.SetTag("currencies", currencies)

	logger.Info("updateRates: fetching rates...")
	rates, err := p.api.FetchRatesToRub(ctx, currencies, date)
	if err != nil {
		err = errors.WithMessage(err, "error fetching rates in updateRates")
		return nil, err
	}

	logger.Info("updateRates: rates fetched successfully, saving to storage...")
	err = p.storage.AddRates(ctx, rates)
	if err != nil {
		logger.Error("error fetching rates", zap.Error(err))
		err = errors.WithMessage(err, "error in updateRates")
		return nil, err
	}

	logger.Info("updateRates: rates saved to storage")

	return rates, nil
}

func rateToSelf(cur entities.Currency, date time.Time) entities.Rate {
	return entities.Rate{
		From:  cur,
		To:    cur,
		Value: decimal.NewFromInt32(1),
		Date:  date,
	}
}
