package sqlt

import (
	"context"
)

type (
	SqlDescriber interface {
		GetSql(c context.Context) (string, context.Context, error)
		Release()
	}

	SqlAssembler interface {
		AssembleSql(id string, data interface{}) (SqlDescriber, error)
	}

	RowScanner interface {
		Columns() ([]string, error)
		Scan(dest ...interface{}) error
		Err() error
	}

	RowsHandler interface {
		HandleRow(r RowScanner)
	}

	MultiRowsHandler interface {
		RowsHandler
		AddResultSet()
	}
)
