package norm

import (
	"context"
	"errors"

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

		lastError   error
		isCommitted bool
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

func (s *TxNorm) Query() *TxNorm {
	if !s.isCommitted {
		if s.lastError == nil {
			if s.lastError = s.op.QueryContext(s.c, s.id, s.param, s.mrh); s.lastError != nil {
				if s.autoRollback {
					s.Rollback()
				}
			}
		}
	}
	return s
}

func (s *TxNorm) Exec() *TxNorm {
	if !s.isCommitted {
		if s.lastError == nil {
			if _, s.lastError = s.op.ExecContext(s.c, s.id, s.param); s.lastError != nil {
				if s.autoRollback {
					s.Rollback()
				}
			}
		}
	}
	return s
}

func (s *TxNorm) ExecRtn() *TxNorm {
	if !s.isCommitted {
		if s.lastError == nil {
			if s.lastError = s.op.ExecRtnContext(s.c, s.id, s.param, s.mrh); s.lastError != nil {
				if s.autoRollback {
					s.Rollback()
				}
			}
		}
	}
	return s
}

func (s *TxNorm) Rollback() error {
	s.isCommitted = true
	return s.op.Rollback()
}

func (s *TxNorm) Commit() error {	
	if s.lastError != nil {
		return s.lastError
	}
	if s.isCommitted {
		return errors.New("TxNorm is already committed!")
	}	
	s.c = nil
	s.isCommitted = true
	return s.op.Commit()
}
