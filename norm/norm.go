package norm

import (
	"context"
	"database/sql"
	"errors"

	"github.com/it512/sqlt"
)

type (
	CollateFunc func() error

	Collator struct{}
)

func (c Collator) CollateWithFunc(cf CollateFunc) {
	if cf != nil {
		if e := cf(); e != nil {
			panic(e)
		}
	}
}

type (
	Norm struct {
		dbop *sqlt.DbOp
	}
)

func (n *Norm) NewSyncNorm(c context.Context) *SyncNorm {
	return &SyncNorm{op: n.dbop, param: make(map[string]interface{}), c: c}
}

func (n *Norm) NewSyncNormDefault() *SyncNorm {
	return n.NewSyncNorm(context.Background())
}

func (n *Norm) NewTxNorm(c context.Context, i interface{}, opt *sql.TxOptions) *TxNorm {
	tx, e := n.dbop.BeginTxWithDb(c, i, opt)
	if e != nil {
		panic(e)
	}
	return &TxNorm{op: tx, param: make(map[string]interface{}), c: c, autoRollback: true}
}

func (n *Norm) NewTxNormDefault(opt *sql.TxOptions) *TxNorm {
	return n.NewTxNorm(context.Background(), nil, opt)
}

func (n *Norm) NewTxNormWithContext(c context.Context, opt *sql.TxOptions) *TxNorm {
	return n.NewTxNorm(c, nil, opt)
}

func NewNorm(dbop *sqlt.DbOp) *Norm {
	return &Norm{dbop: dbop}
}

var (
	errorStatus = errors.New("context is nil!")
)

func mustCheckContext(c context.Context) {
	if c == nil {
		panic(errorStatus)
	}
}
