package tx

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/storage"
)

type db interface {
	storage.DbTx
	BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
}
