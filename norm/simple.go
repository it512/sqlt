package norm

import (
	"context"

	"github.com/it512/sqlt"
)

type (
	SimpleNorm struct {
		op    *sqlt.DbOp
		id    string
		param map[string]interface{}
		mrh   sqlt.MultiRowsHandler
		c     context.Context

		lastError error
	}
)

func (s *SimpleNorm) With(id string, p map[string]interface{}, mrh sqlt.MultiRowsHandler) *SimpleNorm {
	s.id = id
	s.AddAll(p)
	s.mrh = mrh

	return s
}

func (s *SimpleNorm) WithId(id string) *SimpleNorm {
	s.id = id
	return s
}

func (s *SimpleNorm) WithHandler(mrh sqlt.MultiRowsHandler) *SimpleNorm {
	s.mrh = mrh
	return s
}

func (s *SimpleNorm) AddAll(m map[string]interface{}) *SimpleNorm {
	if m != nil {
		for k, v := range m {
			s.param[k] = v
		}
	}
	return s
}

func (s *SimpleNorm) AddParam(k string, v interface{}) *SimpleNorm {
	if k != "" && v != nil {
		s.param[k] = v
	}
	return s
}

func (s *SimpleNorm) RemoveParam(k string) *SimpleNorm {
	delete(s.param, k)
	return s
}

func (s *SimpleNorm) Query() (sqlt.MultiRowsHandler, error) {
	s.lastError = s.op.QueryContext(s.c, s.id, s.param, s.mrh)
	return s.mrh, s.lastError
}

func (s *SimpleNorm) Exec() (int64, error) {
	return s.op.ExecContext(s.c, s.id, s.param)
}

func (s *SimpleNorm) ExecRtn() (sqlt.MultiRowsHandler, error) {
	s.lastError = s.op.ExecRtnContext(s.c, s.id, s.param, s.mrh)
	return s.mrh, s.lastError
}
