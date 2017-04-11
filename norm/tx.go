package norm

import "github.com/it512/sqlt"

type (
	TxNorm struct {
		op *sqlt.TxOp
		ctx
		autoRollback bool
	}
)

func (s *TxNorm) AutoRollback(b bool) *TxNorm {
	s.autoRollback = b
	return s
}

func (s *TxNorm) WithId(id string) *TxNorm {
	s.id = id
	return s
}

func (s *TxNorm) AddParam(k string, v interface{}) *TxNorm {
	if k != "" && v != nil {
		s.param[k] = v
	}
	return s
}

func (s *TxNorm) WithHandler(mrh sqlt.MultiRowsHandler) *TxNorm {
	s.mrh = mrh
	return s
}

func (s *TxNorm) ResetAll() *TxNorm {
	s.ctx = ctx{param: make(map[string]interface{})}
	return s
}

func (s *TxNorm) Reset() *TxNorm {
	s.ctx = ctx{param: s.param}
	return s
}

func (s *TxNorm) Query() *TxNorm {
	e := s.op.Query(s.id, s.param, s.mrh)
	if e != nil {
		if s.autoRollback {
			s.Rollback()
		}
		panic(e)
	}

	return s.Reset()
}

func (s *TxNorm) Exec() *TxNorm {
	_, e := s.op.Exec(s.id, s.param)
	if e != nil {
		if s.autoRollback {
			s.Rollback()
		}
		panic(e)
	}
	return s.Reset()
}

func (s *TxNorm) ExecReturning() *TxNorm {
	e := s.op.ExecReturning(s.id, s.param, s.mrh)
	if e != nil {
		if s.autoRollback {
			s.Rollback()
		}
		panic(e)
	}

	return s.Reset()
}

func (s *TxNorm) Rollback() {
	if err := s.op.Rollback(); err != nil {
		panic(err)
	}
}

func (s *TxNorm) Commit() Collator {
	e := s.op.Commit()
	if e != nil {
		if s.autoRollback {
			s.Rollback()
		}
		panic(e)
	}
	return Collator{}
}

func NewTxNorm(op *sqlt.DbOp) *TxNorm {
	return &TxNorm{autoRollback: true, ctx: ctx{param: make(map[string]interface{})}}
}
