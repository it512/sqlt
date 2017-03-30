package sqlt

type (
	SqlLoader interface {
		LoadSql(id string, data interface{}) (string, int, error)
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

const (
	SQL_TYPE_NORMAL   = 0
	SQL_TYPE_READONLY = 1
)
