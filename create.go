package sqlt

import "github.com/jmoiron/sqlx"

func New(db *sqlx.DB, loader SqlLoader) *DbOp {
	return &DbOp{db: db, l: loader}
}
