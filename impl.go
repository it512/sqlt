package sqlt

import "github.com/jmoiron/sqlx"

type (
	DbOp struct {
		db *sqlx.DB
		l  SqlLoader
	}

	TxOp struct {
		tx *sqlx.Tx
		l  SqlLoader
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

func (c *DbOp) SelectWithRowHandler(id string, param interface{}, rh RowHandler) error {
	return selectWithRowHandler(c.l, c.db, id, param, rh)
}

func (c *DbOp) Select(id string, param interface{}) ([]map[string]interface{}, error) {
	rh := &defRowHandler{vs: make([]map[string]interface{}, 0, 10)}
	e := c.SelectWithRowHandler(id, param, rh)
	return rh.vs, e
}

func (c *DbOp) InsertDeleteUpdate(id string, param interface{}) (int64, error) {
	return insertDeleteUpdate(c.l, c.db, id, param)
}

func (c *DbOp) Begin() (*TxOp, error) {
	tx, e := c.db.Beginx()
	if e != nil {
		return nil, e
	}
	return &TxOp{tx: tx, l: c.l}, nil
}

func (t *TxOp) SelectWithRowHandler(id string, param interface{}, rh RowHandler) error {
	return selectWithRowHandler(t.l, t.tx, id, param, rh)
}

func (t *TxOp) Select(id string, param interface{}) ([]map[string]interface{}, error) {
	rh := &defRowHandler{vs: make([]map[string]interface{}, 0, 10)}
	e := t.SelectWithRowHandler(id, param, rh)
	return rh.vs, e
}

func (t *TxOp) InsertDeleteUpdate(id string, param interface{}) (int64, error) {
	return insertDeleteUpdate(t.l, t.tx, id, param)
}

func (t *TxOp) Commit() error {
	return t.tx.Commit()
}

func (t *TxOp) Rollback() error {
	return t.tx.Rollback()
}
