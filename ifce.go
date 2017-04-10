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

	DbOper interface {
		Query(id string, data interface{}, mrh MultiRowsHandler) error
		Exec(id string, data interface{}) (int64, error)
		ExecReturning(id string, data interface{}, mrh MultiRowsHandler) error
	}
)
