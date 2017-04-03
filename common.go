package sqlt

import "github.com/jmoiron/sqlx"

var (
	nilParam = make(map[string]interface{})
)

func insertDeleteUpdate(msql MappedSql, ext sqlx.Ext) (int64, error) {
	sql, e := msql.GetSql()
	if e != nil {
		return -1, e
	}

	var p = msql.GetParam()
	if p == nil {
		p = nilParam
	}
	result, e := sqlx.NamedExec(ext, sql, p)
	if e != nil {
		return -1, e
	}
	return result.RowsAffected()
}

func selectWithRowHandler(msql MappedSql, ext sqlx.Ext, rh RowHandler) error {
	sql, e := msql.GetSql()
	if e != nil {
		return e
	}

	var p = msql.GetParam()
	if p == nil {
		p = nilParam
	}
	rows, e := sqlx.NamedQuery(ext, sql, p)
	if e != nil {
		return e
	}

	for rows.Next() {
		rh.HandleRow(rows)
	}

	return rows.Close()
}
