package sqlt

type (
	SqlDescriber interface {
		GetSql() (string, error)
	}

	SqlAssembler interface {
		AssembleSql(id string, data interface{}) (SqlDescriber, error)
	}

	RowScanner interface {
		Columns() ([]string, error)
		Scan(dest ...interface{}) error
		Err() error
	}

	MultiRowsHandler interface {
		AddResultSet()
		HandleRow(r RowScanner)
	}
)
