package sqlt

type (
	MappedSql interface {
		GetSql() (string, error)
		GetParam() interface{}
		GetId() string
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
