package sqlt

import (
	dbset "github.com/it512/dsds/simple"
	"github.com/jmoiron/sqlx"
)

func New(db *sqlx.DB, loader SqlLoader) *DbOp {
	return &DbOp{dbset: dbset.NewSimpleDbSet(db), l: loader}
}
