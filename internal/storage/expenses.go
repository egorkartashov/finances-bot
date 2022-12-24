package storage

import (
	"context"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/entities"
)

type expense struct {
	ID       int64           `db:"id"`
	UserID   int64           `db:"user_id"`
	Category string          `db:"category"`
	Sum      decimal.Decimal `db:"sum_rub"`
	Date     time.Time       `db:"date"`
}

type Expenses struct {
	db dbTxStorage
}

func NewExpenses(db dbTxStorage) *Expenses {
	return &Expenses{db}
}

func (s *Expenses) AddExpense(ctx context.Context, userID int64, exp entities.Expense) error {
	expModel := expense{
		UserID:   userID,
		Category: exp.Category,
		Sum:      exp.Sum,
		Date:     exp.Date,
	}

	const query = `
INSERT INTO expenses (user_id, category, sum_rub, date) 
VALUES (:user_id, :category, :sum_rub, :date)`

	_, err := s.db.Db(ctx).NamedExecContext(ctx, query, expModel)
	return err
}

func (s *Expenses) GetExpenses(ctx context.Context, userID int64, startDate time.Time) ([]entities.Expense, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "get-expenses-from-db")
	defer span.Finish()

	const query = `SELECT * FROM expenses WHERE user_id = $1 AND date >= $2`

	var expModels []expense
	if err := s.db.Db(ctx).SelectContext(ctx, &expModels, query, userID, startDate); err != nil {
		return nil, err
	}

	expenses := make([]entities.Expense, len(expModels))
	for i, model := range expModels {
		expenses[i] = entities.Expense{
			Category: model.Category,
			Sum:      model.Sum,
			Date:     model.Date,
		}
	}

	return expenses, nil
}

type sumModel struct {
	Sum decimal.Decimal `db:"sum"`
}

func (s *Expenses) GetSumForCategoryAndPeriod(
	ctx context.Context, userID int64, category string, startDate, endDate time.Time,
) (decimal.Decimal, error) {
	const query = `
SELECT COALESCE(SUM(sum_rub),0) AS sum 
FROM expenses 
WHERE user_id = $1 AND category = $2 AND date >= $3 AND date <= $4`

	var model sumModel
	if err := s.db.Db(ctx).GetContext(ctx, &model, query, userID, category, startDate, endDate); err != nil {
		return decimal.Decimal{}, err
	}

	return model.Sum, nil
}
