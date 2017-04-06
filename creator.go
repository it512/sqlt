package sqlt

import (
	dbset "github.com/it512/dsds"
)

func New(dbset dbset.DbSet, loader SqlLoader) *DbOp {
	return &DbOp{dbset: dbset, l: loader}
}
