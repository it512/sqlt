package notorm

import (
	"sync"

	"github.com/it512/sqlt"
)

type (
	AsyncOp struct {
		op *sqlt.DbOp
		ctx
		wg sync.WaitGroup
	}
)

func (s *AsyncOp) WithId(id string) *AsyncOp {
	s.id = id
	return s
}

func (s *AsyncOp) WithHandler(mrh sqlt.MultiRowsHandler) *AsyncOp {
	s.mrh = mrh
	return s
}

func (s *AsyncOp) AddParam(k string, v interface{}) *AsyncOp {
	if k != "" && v != nil {
		s.param[k] = v
	}
	return s
}

func (s *AsyncOp) ResetAll() *AsyncOp {
	s.ctx = ctx{param: make(map[string]interface{})}
	return s
}

func (s *AsyncOp) Reset() *AsyncOp {
	s.ctx = ctx{param: s.param}
	return s
}

func (s *AsyncOp) Query() *AsyncOp {
	s.wg.Add(1)
	go func(id string, param map[string]interface{}, mrh sqlt.MultiRowsHandler) {
		defer s.wg.Done()
		e := s.op.Query(id, param, mrh)
		if e != nil {
			panic(e)
		}
	}(s.id, s.param, s.mrh)

	return s.Reset()
}

func (s *AsyncOp) Exec() *AsyncOp {
	s.wg.Add(1)
	go func(id string, param map[string]interface{}) {
		defer s.wg.Done()
		_, e := s.op.Exec(id, param)
		if e != nil {
			panic(e)
		}
	}(s.id, s.param)

	return s.Reset()
}

func (s *AsyncOp) ExecReturning() *AsyncOp {
	s.wg.Add(1)
	go func(id string, param map[string]interface{}, mrh sqlt.MultiRowsHandler) {
		defer s.wg.Done()
		e := s.op.ExecReturning(id, param, mrh)
		if e != nil {
			panic(e)
		}
	}(s.id, s.param, s.mrh)

	return s.Reset()
}

func (s *AsyncOp) Wait() Collator {
	s.wg.Wait()
	return Collator{}
}

func NewAsyncOp(op *sqlt.DbOp) *AsyncOp {
	return &AsyncOp{op: op, ctx: ctx{param: make(map[string]interface{})}}
}
