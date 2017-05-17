package sqlt

import (
	"context"

	"github.com/it512/dsds"
	"github.com/jmoiron/sqlx"
)

var (
	defData = struct{}{}
)

func dataIsNilWithDef(data interface{}) interface{} {
	if data == nil {
		return defData
	}
	return data
}

func processSqlDescriberWithDbManager(m dsds.DbManager, assembler SqlAssembler, id string, data interface{}) (*sqlx.DB, SqlDescriber, error) {
	desc, e := assembler.AssembleSql(id, data)
	if e != nil {
		return nil, nil, e
	}

	p, e := m.GetDb(desc)
	if e != nil {
		return nil, nil, e
	}

	return p, desc, nil
}

func processSqlDescriberWithTx(tx *sqlx.Tx, assembler SqlAssembler, id string, param interface{}) (*sqlx.Tx, SqlDescriber, error) {
	desc, e := assembler.AssembleSql(id, param)
	if e != nil {
		return nil, nil, e
	}

	return tx, desc, nil
}

func processMultiRowsHander(rows *sqlx.Rows, mrh MultiRowsHandler) {
	mrh.AddResultSet()
	for rows.Next() {
		mrh.HandleRow(rows)
	}

	for rows.NextResultSet() {
		mrh.AddResultSet()
		for rows.Next() {
			mrh.HandleRow(rows)
		}
	}
}

func query(c context.Context, ext sqlx.ExtContext, descr SqlDescriber, data interface{}, mrh MultiRowsHandler) error {
	sql, ctx, e := descr.GetSql(c)
	if e != nil {
		return e
	}

	rows, e := sqlx.NamedQueryContext(ctx, ext, sql, dataIsNilWithDef(data))
	if e != nil {
		return e
	}

	processMultiRowsHander(rows, mrh)
	return rows.Close()
}

func exec(c context.Context, ext sqlx.ExtContext, descr SqlDescriber, data interface{}) (int64, error) {
	sql, ctx, e := descr.GetSql(c)
	if e != nil {
		return -1, e
	}

	r, e := sqlx.NamedExecContext(ctx, ext, sql, dataIsNilWithDef(data))
	if e != nil {
		return -1, e
	}

	return r.RowsAffected()
}
