package notorm

import "github.com/it512/sqlt"

type (
	SyncOp struct {
		op *sqlt.DbOp
		ctx
	}
)

func (s *SyncOp) WithId(id string) *SyncOp {
	s.id = id
	return s
}

func (s *SyncOp) WithHandler(mrh sqlt.MultiRowsHandler) *SyncOp {
	s.mrh = mrh
	return s
}

func (s *SyncOp) AddParam(k string, v interface{}) *SyncOp {
	if k != "" && v != nil {
		s.param[k] = v
	}
	return s
}

func (s *SyncOp) ResetAll() *SyncOp {
	s.ctx = ctx{param: make(map[string]interface{})}
	return s
}

func (s *SyncOp) Reset() *SyncOp {
	s.ctx = ctx{param: s.param}
	return s
}

func (s *SyncOp) Query() *SyncOp {
	e := s.op.Query(s.id, s.param, s.mrh)
	if e != nil {
		panic(e)
	}

	return s.Reset()
}

func (s *SyncOp) Exec() *SyncOp {
	_, e := s.op.Exec(s.id, s.param)
	if e != nil {
		panic(e)
	}
	return s.Reset()
}

func (s *SyncOp) ExecReturning() *SyncOp {
	e := s.op.ExecReturning(s.id, s.param, s.mrh)
	if e != nil {
		panic(e)
	}

	return s.Reset()
}

func (s *SyncOp) Finish() Collator {
	return Collator{}
}

func NewSyncOp(op *sqlt.DbOp) *SyncOp {
	return &SyncOp{op: op, ctx: ctx{param: make(map[string]interface{})}}
}
