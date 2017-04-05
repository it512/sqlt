package sqlt

import "github.com/jmoiron/sqlx"

var (
	defParam = make(map[string]interface{})
)

func checkParamNilWithDef(param interface{}, def interface{}) interface{} {
	if param == nil {
		return param
	}

	return def
}

func insertDeleteUpdate(msql MappedSql, ext sqlx.Ext, param interface{}) (int64, error) {
	sql, e := msql.GetSql()
	if e != nil {
		return -1, e
	}

	result, e := sqlx.NamedExec(ext, sql, checkParamNilWithDef(param, defParam))
	if e != nil {
		return -1, e
	}
	return result.RowsAffected()
}

func selectWithRowHandler(msql MappedSql, ext sqlx.Ext, param interface{}, rh RowHandler) error {
	sql, e := msql.GetSql()
	if e != nil {
		return e
	}

	rows, e := sqlx.NamedQuery(ext, sql, checkParamNilWithDef(param, defParam))
	if e != nil {
		return e
	}

	for rows.Next() {
		rh.HandleRow(rows)
	}

	return rows.Close()
}
