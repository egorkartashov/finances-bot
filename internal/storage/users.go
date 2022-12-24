package storage

import (
	"context"
	"database/sql"

	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/entities"
)

type Users struct {
	db dbTxStorage
}

func NewUsers(db dbTxStorage) *Users {
	return &Users{db}
}

type userModel struct {
	ID       int64  `db:"id"`
	Currency string `db:"currency"`
}

func (s *Users) Save(ctx context.Context, user entities.User) (err error) {
	model := userModel{
		ID:       user.ID,
		Currency: string(user.Currency),
	}
	const query = `
INSERT INTO users(id, currency) 
VALUES (:id, :currency) 
ON CONFLICT ON CONSTRAINT users_pkey DO UPDATE 
SET currency = EXCLUDED.currency`

	_, err = s.db.Db(ctx).NamedExecContext(ctx, query, model)
	return
}

func (s *Users) Get(ctx context.Context, id int64) (user entities.User, ok bool, err error) {
	const query = `SELECT id, currency FROM users WHERE id = $1`
	var model userModel
	if err = s.db.Db(ctx).GetContext(ctx, &model, query, id); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			ok = false
		}
		return
	}

	ok = true
	user = entities.User{
		ID:       model.ID,
		Currency: entities.Currency(model.Currency),
	}

	return
}
