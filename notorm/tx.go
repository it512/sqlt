package notorm

import "github.com/it512/sqlt"

type (
	SyncTxOp struct {
		op *sqlt.TxOp
		ctx
		autoRollback bool
	}
)

func (s *SyncTxOp) AutoRollback(b bool) *SyncTxOp {
	s.autoRollback = b
	return s
}

func (s *SyncTxOp) WithId(id string) *SyncTxOp {
	s.id = id
	return s
}

func (s *SyncTxOp) AddParam(k string, v interface{}) *SyncTxOp {
	if k != "" && v != nil {
		s.param[k] = v
	}
	return s
}

func (s *SyncTxOp) WithHandler(mrh sqlt.MultiRowsHandler) *SyncTxOp {
	s.mrh = mrh
	return s
}

func (s *SyncTxOp) ResetAll() *SyncTxOp {
	s.ctx = ctx{param: make(map[string]interface{})}
	return s
}

func (s *SyncTxOp) Reset() *SyncTxOp {
	s.ctx = ctx{param: s.param}
	return s
}

func (s *SyncTxOp) Query() *SyncTxOp {
	e := s.op.Query(s.id, s.param, s.mrh)
	if e != nil {
		if s.autoRollback {
			s.Rollback()
		}
		panic(e)
	}

	return s.Reset()
}

func (s *SyncTxOp) Exec() *SyncTxOp {
	_, e := s.op.Exec(s.id, s.param)
	if e != nil {
		if s.autoRollback {
			s.Rollback()
		}
		panic(e)
	}
	return s.Reset()
}

func (s *SyncTxOp) ExecReturning() *SyncTxOp {
	e := s.op.ExecReturning(s.id, s.param, s.mrh)
	if e != nil {
		if s.autoRollback {
			s.Rollback()
		}
		panic(e)
	}

	return s.Reset()
}

func (s *SyncTxOp) Rollback() {
	if err := s.op.Rollback(); err != nil {
		panic(err)
	}
}

func (s *SyncTxOp) Commit() Collator {
	e := s.op.Commit()
	if e != nil {
		if s.autoRollback {
			s.Rollback()
		}
		panic(e)
	}
	return Collator{}
}

func NewSyncTxOp(op *sqlt.DbOp) *SyncTxOp {
	return &SyncTxOp{autoRollback: true, ctx: ctx{param: make(map[string]interface{})}}
}
