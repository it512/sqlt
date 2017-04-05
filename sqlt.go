package sqlt

type (
	MappedSql interface {
		GetSql() (string, error)
	}

	SqlLoader interface {
		LoadSql(id string, data interface{}) (MappedSql, error)
	}

	ColScanner interface {
		Columns() ([]string, error)
		Scan(dest ...interface{}) error
		Err() error
	}

	RowHandler interface {
		HandleRow(r ColScanner)
	}
)
