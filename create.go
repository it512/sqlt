package sqlt

import "github.com/jmoiron/sqlx"

func New(driverName, dataSourceName, pattern string) *DbOp {
	db := sqlx.MustOpen(driverName, dataSourceName)
	loader := newTemplateSqlLoader(pattern)
	return &DbOp{db: db, l: loader}
}

func NewWithTemplateLoader(driverName, dataSourceName string, loader SqlLoader) *DbOp {
	db := sqlx.MustOpen(driverName, dataSourceName)
	return &DbOp{db: db, l: loader}
}
