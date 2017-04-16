package norm

import (
	"context"

	"github.com/it512/sqlt"
)

type (
	SyncNorm struct {
		op    *sqlt.DbOp
		id    string
		param map[string]interface{}
		mrh   sqlt.MultiRowsHandler
		c     context.Context
	}
)

func (s *SyncNorm) WithId(id string) *SyncNorm {
	s.id = id
	return s
}

func (s *SyncNorm) WithHandler(mrh sqlt.MultiRowsHandler) *SyncNorm {
	s.mrh = mrh
	return s
}

func (s *SyncNorm) AddAll(m map[string]interface{}) *SyncNorm {
	for k, v := range m {
		s.param[k] = v
	}
	return s
}

func (s *SyncNorm) AddParam(k string, v interface{}) *SyncNorm {
	if k != "" && v != nil {
		s.param[k] = v
	}
	return s
}

func (s *SyncNorm) RemoveParam(k string) *SyncNorm {
	delete(s.param, k)
	return s
}

func (s *SyncNorm) Reset() *SyncNorm {
	s.param = make(map[string]interface{})
	s.id = ""
	s.mrh = nil
	return s
}

func (s *SyncNorm) Query() *SyncNorm {
	e := s.op.QueryContext(s.c, s.id, s.param, s.mrh)
	if e != nil {
		panic(e)
	}
	return s
}

func (s *SyncNorm) Exec() *SyncNorm {
	_, e := s.op.ExecContext(s.c, s.id, s.param)
	if e != nil {
		panic(e)
	}
	return s
}

func (s *SyncNorm) ExecRtn() *SyncNorm {
	e := s.op.ExecRtnContext(s.c, s.id, s.param, s.mrh)
	if e != nil {
		panic(e)
	}
	return s
}

func (s *SyncNorm) Finish() Collator {
	s.c = nil
	return Collator{}
}

func (s *SyncNorm) Cancel() {
	s.c = nil
}
