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

	SqlOperator interface {
		InsertUpdateDelete(id string, p interface{}) (int64, error)

		Select(id string, p interface{}) ([]map[string]interface{}, error)
		SelectWithRowHandler(id string, p interface{}, h RowHandler) error
	}

	SqlDbOperator interface {
		Begin() (SqlTxOperator, error)
		SqlOperator
	}

	SqlTxOperator interface {
		Commit() error
		Rollback() error

		SqlOperator
	}
)
