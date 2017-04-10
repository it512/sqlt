package sqlt

import (
	"context"
	"database/sql"

	"github.com/it512/dsds"
	"github.com/jmoiron/sqlx"
)

type (
	DbOp struct {
		assembler SqlAssembler
		manager   dsds.DbManager
	}
)

func (c *DbOp) Query(id string, data interface{}, mrh MultiRowsHandler) error {
	p, d, e := processSqlDescriberWithDbManager(c.manager, c.assembler, id, data)
	if e != nil {
		return e
	}
	return query(p, d, data, mrh)
}

func (c *DbOp) Exec(id string, data interface{}) (int64, error) {
	p, d, e := processSqlDescriberWithDbManager(c.manager, c.assembler, id, data)
	if e != nil {
		return -1, e
	}
	return exec(p, d, data)
}

func (c *DbOp) ExecReturning(id string, data interface{}, mrh MultiRowsHandler) error {
	p, d, e := processSqlDescriberWithDbManager(c.manager, c.assembler, id, data)
	if e != nil {
		return e
	}
	return execRetruning(p, d, data, mrh)
}

func (c *DbOp) BeginWithDb(i interface{}, opt *sql.TxOptions) (*TxOp, error) {
	db, e := c.manager.GetDb(i)
	if e != nil {
		return nil, e
	}

	tx, e := db.BeginTxx(context.Background(), opt)
	if e != nil {
		return nil, e
	}

	return &TxOp{tx: tx, assembler: c.assembler}, nil
}

func (c *DbOp) Begin() (*TxOp, error) {
	return c.BeginWithDb(nil, nil)
}

type (
	TxOp struct {
		tx        *sqlx.Tx
		assembler SqlAssembler
	}
)

func (c *TxOp) Query(id string, data interface{}, mrh MultiRowsHandler) error {
	p, d, e := processSqlDescriberWithTx(c.tx, c.assembler, id, data)
	if e != nil {
		return e
	}
	return query(p, d, data, mrh)
}

func (c *TxOp) Exec(id string, data interface{}) (int64, error) {
	p, d, e := processSqlDescriberWithTx(c.tx, c.assembler, id, data)
	if e != nil {
		return -1, e
	}
	return exec(p, d, data)
}

func (c *TxOp) ExecReturning(id string, data interface{}, mrh MultiRowsHandler) error {
	p, d, e := processSqlDescriberWithTx(c.tx, c.assembler, id, data)
	if e != nil {
		return e
	}
	return execRetruning(p, d, data, mrh)
}

func (c *TxOp) Commit() error {
	return c.tx.Commit()
}

func (c *TxOp) Rollback() error {
	return c.tx.Rollback()
}

func New(dbset dsds.DbManager, loader SqlAssembler) *DbOp {
	return &DbOp{manager: dbset, assembler: loader}
}
