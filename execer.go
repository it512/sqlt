package sqlt

import (
	"github.com/it512/dsds"
	"github.com/jmoiron/sqlx"
)

type (
	prepareNameder interface {
		PrepareNamed(query string) (*sqlx.NamedStmt, error)
	}
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

func processSqlDescriberWithDbManager(m dsds.DbManager, assembler SqlAssembler, id string, data interface{}) (prepareNameder, SqlDescriber, error) {
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

func processSqlDescriberWithTx(tx *sqlx.Tx, assembler SqlAssembler, id string, param interface{}) (prepareNameder, SqlDescriber, error) {
	desc, e := assembler.AssembleSql(id, param)
	if e != nil {
		return nil, nil, e
	}

	return tx, desc, nil
}

func processNamedStmt(p prepareNameder, desc SqlDescriber) (*sqlx.NamedStmt, error) {
	sql, e := desc.GetSql()
	if e != nil {
		return nil, e
	}
	return p.PrepareNamed(sql)
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

func query(p prepareNameder, desc SqlDescriber, data interface{}, mrh MultiRowsHandler) error {
	st, e := processNamedStmt(p, desc)
	if e != nil {
		return e
	}

	rows, e := st.Queryx(dataIsNilWithDef(data))
	if e != nil {
		return e
	}

	processMultiRowsHander(rows, mrh)

	defer rows.Close()
	defer st.Close()

	return nil
}

func exec(p prepareNameder, desc SqlDescriber, data interface{}) (int64, error) {
	st, e := processNamedStmt(p, desc)
	if e != nil {
		return -1, e
	}
	defer st.Close()
	r, e := st.Exec(dataIsNilWithDef(data))

	return r.RowsAffected()
}

func execRetruning(p prepareNameder, desc SqlDescriber, data interface{}, mrh MultiRowsHandler) error {
	return query(p, desc, data, mrh)
}
