package sqlt

import "github.com/jmoiron/sqlx"

var (
	nilParam = make(map[string]interface{})
)

func insertDeleteUpdate(l SqlRenderer, ext sqlx.Ext, id string, param interface{}) (int64, error) {
	sql, _, e := l.RenderSql(id, param)
	oops(e)

	var p = param
	if p == nil {
		p = nilParam
	}
	result, e := sqlx.NamedExec(ext, sql, p)
	if e != nil {
		return 0, e
	}
	return result.RowsAffected()
}

func selectWithRowHandler(l SqlRenderer, ext sqlx.Ext, id string, param interface{}, rh RowHandler) error {
	sql, _, e := l.RenderSql(id, param)
	oops(e)

	var p = param
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
