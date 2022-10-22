package storage

import (
	"context"
	"database/sql"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/entities"
)

type Rates struct {
	db dbTxStorage
}

func NewRates(db dbTxStorage) *Rates {
	return &Rates{db}
}

type rateModel struct {
	ID           int32           `db:"id"`
	FromCurrency string          `db:"from_currency"`
	Rate         decimal.Decimal `db:"rate"`
	Date         string          `db:"date"`
}

func (r *Rates) AddRates(ctx context.Context, rates []entities.Rate) (err error) {
	models := make([]rateModel, len(rates))
	for i, rate := range rates {
		models[i] = rateModel{
			FromCurrency: string(rate.From),
			Rate:         rate.Value,
			Date:         rate.Date.Format("2006-01-02"),
		}
	}

	const insertQuery = `
INSERT INTO exchange_rates_to_rub(rate, from_currency, date) 
VALUES (:rate, :from_currency, :date)
ON CONFLICT DO NOTHING`

	if _, err = r.db.Db(ctx).NamedExecContext(ctx, insertQuery, models); err != nil {
		return err
	}

	return err
}

func (r *Rates) GetRate(ctx context.Context, from entities.Currency, date time.Time) (
	rate decimal.Decimal, ok bool, err error,
) {
	const query = `SELECT * FROM exchange_rates_to_rub WHERE from_currency = $1 AND date = $2`

	var model rateModel
	if err = r.db.Db(ctx).GetContext(ctx, &model, query, from, date); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			ok = false
		}
		return
	}

	rate = model.Rate
	ok = true

	return
}
