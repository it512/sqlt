package norm

import (
	"context"

	"github.com/it512/sqlt"
)

type (
	TxNorm struct {
		id           string
		param        map[string]interface{}
		mrh          sqlt.MultiRowsHandler
		c            context.Context
		op           *sqlt.TxOp
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

func (s *TxNorm) WithHandler(mrh sqlt.MultiRowsHandler) *TxNorm {
	s.mrh = mrh
	return s
}

func (s *TxNorm) AddAll(m map[string]interface{}) *TxNorm {
	for k, v := range m {
		s.param[k] = v
	}
	return s
}

func (s *TxNorm) AddParam(k string, v interface{}) *TxNorm {
	if k != "" && v != nil {
		s.param[k] = v
	}
	return s
}

func (s *TxNorm) RemoveParam(k string) *TxNorm {
	delete(s.param, k)
	return s
}

func (s *TxNorm) Reset() *TxNorm {
	s.param = make(map[string]interface{})
	s.id = ""
	s.mrh = nil
	return s
}

func (s *TxNorm) Query() *TxNorm {
	e := s.op.QueryContext(s.c, s.id, s.param, s.mrh)
	if e != nil {
		if s.autoRollback {
			s.Rollback()
		}
		panic(e)
	}

	return s.Reset()
}

func (s *TxNorm) Exec() *TxNorm {
	_, e := s.op.ExecContext(s.c, s.id, s.param)
	if e != nil {
		if s.autoRollback {
			s.Rollback()
		}
		panic(e)
	}
	return s.Reset()
}

func (s *TxNorm) ExecRtn() *TxNorm {
	e := s.op.ExecRtnContext(s.c, s.id, s.param, s.mrh)
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
