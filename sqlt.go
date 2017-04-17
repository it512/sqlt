package sqlt

import (
	"context"
	"database/sql"

	"github.com/it512/dsds"
	"github.com/jmoiron/sqlx"
)

var (
	DefaultTxOptions *sql.TxOptions = new(sql.TxOptions)
)

func CreateTxOptions(level sql.IsolationLevel, readonly bool) *sql.TxOptions {
	return &sql.TxOptions{Isolation: level, ReadOnly: readonly}
}

type (
	DbOp struct {
		assembler SqlAssembler
		manager   dsds.DbManager
	}
)

func (c *DbOp) QueryContext(ctx context.Context, id string, data interface{}, mrh MultiRowsHandler) error {
	p, d, e := processSqlDescriberWithDbManager(c.manager, c.assembler, id, data)
	if e != nil {
		return e
	}
	defer d.Release()
	return query(ctx, p, d, data, mrh)
}

func (c *DbOp) ExecContext(ctx context.Context, id string, data interface{}) (int64, error) {
	p, d, e := processSqlDescriberWithDbManager(c.manager, c.assembler, id, data)
	if e != nil {
		return -1, e
	}
	defer d.Release()
	return exec(ctx, p, d, data)
}

func (c *DbOp) ExecRtnContext(ctx context.Context, id string, data interface{}, mrh MultiRowsHandler) error {
	return c.QueryContext(ctx, id, data, mrh)
}

func (c *DbOp) Query(id string, data interface{}, mrh MultiRowsHandler) error {
	return c.QueryContext(context.Background(), id, data, mrh)
}

func (c *DbOp) Exec(id string, data interface{}) (int64, error) {
	return c.ExecContext(context.Background(), id, data)
}

func (c *DbOp) ExecRtn(id string, data interface{}, mrh MultiRowsHandler) error {
	return c.ExecRtnContext(context.Background(), id, data, mrh)
}

func (c *DbOp) BeginTxWithDb(ctx context.Context, i interface{}, opt *sql.TxOptions) (*TxOp, error) {
	db, e := c.manager.GetDb(i)
	if e != nil {
		return nil, e
	}

	tx, e := db.BeginTxx(ctx, opt)
	if e != nil {
		return nil, e
	}

	return &TxOp{tx: tx, assembler: c.assembler}, nil
}

func (c *DbOp) BeginTx(ctx context.Context, opt *sql.TxOptions) (*TxOp, error) {
	return c.BeginTxWithDb(ctx, nil, opt)
}

type (
	TxOp struct {
		tx        *sqlx.Tx
		assembler SqlAssembler
	}
)

func (c *TxOp) QueryContext(ctx context.Context, id string, data interface{}, mrh MultiRowsHandler) error {
	p, d, e := processSqlDescriberWithTx(c.tx, c.assembler, id, data)
	if e != nil {
		return e
	}
	defer d.Release()
	return query(ctx, p, d, data, mrh)
}

func (c *TxOp) ExecContext(ctx context.Context, id string, data interface{}) (int64, error) {
	p, d, e := processSqlDescriberWithTx(c.tx, c.assembler, id, data)
	if e != nil {
		return -1, e
	}
	defer d.Release()
	return exec(ctx, p, d, data)
}

func (c *TxOp) ExecRtnContext(ctx context.Context, id string, data interface{}, mrh MultiRowsHandler) error {
	return c.QueryContext(ctx, id, data, mrh)
}

func (c *TxOp) Query(id string, data interface{}, mrh MultiRowsHandler) error {
	return c.QueryContext(context.Background(), id, data, mrh)
}

func (c *TxOp) Exec(id string, data interface{}) (int64, error) {
	return c.ExecContext(context.Background(), id, data)
}

func (c *TxOp) ExecRtn(id string, data interface{}, mrh MultiRowsHandler) error {
	return c.ExecRtnContext(context.Background(), id, data, mrh)
}

func (c *TxOp) Commit() error {
	return c.tx.Commit()
}

func (c *TxOp) Rollback() error {
	return c.tx.Rollback()
}

func NewSqlt(dbset dsds.DbManager, loader SqlAssembler) *DbOp {
	return &DbOp{manager: dbset, assembler: loader}
}
