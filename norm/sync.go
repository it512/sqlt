package norm

import "github.com/it512/sqlt"

type (
	SyncNorm struct {
		op *sqlt.DbOp
		ctx
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

func (s *SyncNorm) AddParam(k string, v interface{}) *SyncNorm {
	if k != "" && v != nil {
		s.param[k] = v
	}
	return s
}

func (s *SyncNorm) ResetAll() *SyncNorm {
	s.ctx = ctx{param: make(map[string]interface{})}
	return s
}

func (s *SyncNorm) Reset() *SyncNorm {
	s.ctx = ctx{param: s.param}
	return s
}

func (s *SyncNorm) Query() *SyncNorm {
	e := s.op.Query(s.id, s.param, s.mrh)
	if e != nil {
		panic(e)
	}

	return s.Reset()
}

func (s *SyncNorm) Exec() *SyncNorm {
	_, e := s.op.Exec(s.id, s.param)
	if e != nil {
		panic(e)
	}
	return s.Reset()
}

func (s *SyncNorm) ExecReturning() *SyncNorm {
	e := s.op.ExecReturning(s.id, s.param, s.mrh)
	if e != nil {
		panic(e)
	}

	return s.Reset()
}

func (s *SyncNorm) Finish() Collator {
	return Collator{}
}

func (s *SyncNorm) Cancel() {}

func NewSyncNorm(op *sqlt.DbOp) *SyncNorm {
	return &SyncNorm{op: op, ctx: ctx{param: make(map[string]interface{})}}
}
