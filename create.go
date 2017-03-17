package sqlt

import "github.com/jmoiron/sqlx"

func New(ext sqlx.Ext, loader SqlLoader) *DbOp {
	return &DbOp{db: ext, l: loader}
}
