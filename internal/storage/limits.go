package storage

import (
	"context"
	"database/sql"

	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/entities"
)

type Limits struct {
	db           dbTxStorage
	baseCurrency entities.Currency
}

type limitModel struct {
	UserID   int64           `db:"user_id"`
	Category string          `db:"category"`
	Sum      decimal.Decimal `db:"sum"`
}

func (l *Limits) Save(ctx context.Context, limit entities.MonthBudgetLimit) (err error) {

	const query = `
INSERT INTO month_limits(user_id, category, sum) 
VALUES (:user_id, :category, :sum)
ON CONFLICT ON CONSTRAINT month_limits_pkey DO UPDATE SET sum = excluded.sum`

	model := limitModel{
		UserID:   limit.UserID,
		Category: limit.Category,
		Sum:      limit.Sum,
	}

	_, err = l.db.Db(ctx).NamedExecContext(ctx, query, model)
	return
}

func (l *Limits) Get(ctx context.Context, userID int64, category string) (
	limit entities.MonthBudgetLimit, ok bool, err error,
) {
	const query = `SELECT user_id, category, sum FROM month_limits WHERE user_id = $1 AND category = $2`

	model := limitModel{}
	err = l.db.Db(ctx).GetContext(ctx, &model, query, userID, category)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
			ok = false
		}
		return
	}

	ok = true
	limit = entities.MonthBudgetLimit{
		UserID:   model.UserID,
		Category: model.Category,
		Sum:      model.Sum,
		Currency: l.baseCurrency,
	}
	return
}

func (l *Limits) Delete(ctx context.Context, userID int64, category string) (err error) {
	const query = `DELETE FROM month_limits WHERE user_id = $1 AND category = $2`
	_, err = l.db.Db(ctx).ExecContext(ctx, query, userID, category)
	return
}

func NewLimits(db dbTxStorage, baseCurrency entities.Currency) *Limits {
	return &Limits{
		db:           db,
		baseCurrency: baseCurrency,
	}
}
