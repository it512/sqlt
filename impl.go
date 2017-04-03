package sqlt

import (
	"github.com/it512/dsds"
	"github.com/jmoiron/sqlx"
)

type (
	DbOp struct {
		l     SqlLoader
		dbset dsds.DbSet
	}

	defRowHandler struct {
		vs []map[string]interface{}
	}
)

func (rh *defRowHandler) HandleRow(r ColScanner) {
	v := make(map[string]interface{})
	if sqlx.MapScan(r, v) == nil {
		rh.vs = append(rh.vs, v)
	}
}

func (c *DbOp) processMapped(id string, param interface{}) (sqlx.Ext, MappedSql, error) {
	mappedSql, e := c.l.LoadSql(id, param)
	if e != nil {
		return nil, nil, e
	}

	excuter, e := c.dbset.GetDb(mappedSql)
	if e != nil {
		return nil, nil, e
	}

	return excuter, mappedSql, nil
}

func (c *DbOp) SelectWithRowHandler(id string, param interface{}, rh RowHandler) error {
	excuter, mappedSql, e := c.processMapped(id, param)
	if e != nil {
		return e
	}

	return selectWithRowHandler(mappedSql, excuter, rh)
}

func (c *DbOp) Select(id string, param interface{}) ([]map[string]interface{}, error) {
	rh := &defRowHandler{vs: make([]map[string]interface{}, 0, 10)}
	e := c.SelectWithRowHandler(id, param, rh)
	return rh.vs, e
}

func (c *DbOp) InsertDeleteUpdate(id string, param interface{}) (int64, error) {
	excuter, mappedSql, e := c.processMapped(id, param)
	if e != nil {
		return -1, e
	}
	return insertDeleteUpdate(mappedSql, excuter)
}

func (c *DbOp) BeginTransWithDb(i interface{}) (*TxOp, error) {
	db, e := c.dbset.GetDb(i)
	if e != nil {
		return nil, e
	}

	tx, e := db.Beginx()
	if e != nil {
		return nil, e
	}

	return &TxOp{tx: tx, l: c.l}, nil
}

func (c *DbOp) BeginTrans() (*TxOp, error) {
	return c.BeginTransWithDb(nil)
}

type (
	TxOp struct {
		tx *sqlx.Tx
		l  SqlLoader
	}
)

func (t *TxOp) processMapped(id string, param interface{}) (sqlx.Ext, MappedSql, error) {
	mappedSql, e := t.l.LoadSql(id, param)
	if e != nil {
		return nil, nil, e
	}
	return t.tx, mappedSql, nil
}

func (t *TxOp) SelectWithRowHandler(id string, param interface{}, rh RowHandler) error {
	excuter, mappedSql, e := t.processMapped(id, param)
	if e != nil {
		return e
	}
	return selectWithRowHandler(mappedSql, excuter, rh)
}

func (t *TxOp) Select(id string, param interface{}) ([]map[string]interface{}, error) {
	rh := &defRowHandler{vs: make([]map[string]interface{}, 0, 10)}
	e := t.SelectWithRowHandler(id, param, rh)
	return rh.vs, e
}

func (t *TxOp) InsertDeleteUpdate(id string, param interface{}) (int64, error) {
	excuter, mappedSql, e := t.processMapped(id, param)
	if e != nil {
		return -1, e
	}
	return insertDeleteUpdate(mappedSql, excuter)
}

func (t *TxOp) CommitTrans() error {
	return t.tx.Commit()
}

func (t *TxOp) RollbackTrans() error {
	return t.tx.Rollback()
}
