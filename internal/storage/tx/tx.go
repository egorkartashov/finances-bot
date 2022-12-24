package tx

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/logger"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/storage"
	"go.uber.org/zap"
)

type contextTxKeyType string

const ContextTxKey contextTxKeyType = "psqltx"

type Storage struct {
	db db
}

func New(db db) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	if tx, ok := ctx.Value(ContextTxKey).(*sqlx.Tx); ok {
		err := fn(ctx)
		if err != nil {
			tryRollback(tx)
		}
		return err
	}

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("can't begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tryRollback(tx)
			panic(p)
		} else if err != nil {
			tryRollback(tx)
		} else {
			err = tx.Commit()
		}
	}()

	ctx = context.WithValue(ctx, ContextTxKey, tx)

	err = fn(ctx)

	return err
}

func tryRollback(tx *sqlx.Tx) {
	if err := tx.Rollback(); err != nil {
		logger.Error("failed to rollback", zap.Error(err))
	}
}

func (s *Storage) Db(ctx context.Context) storage.DbTx {
	if tx, ok := ctx.Value(ContextTxKey).(*sqlx.Tx); ok {
		return tx
	}
	return s.db
}
