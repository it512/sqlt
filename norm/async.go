package norm

import (
	"sync"

	"github.com/it512/sqlt"
)

type (
	AsyncNorm struct {
		op *sqlt.DbOp
		ctx
		wg sync.WaitGroup
	}
)

func (s *AsyncNorm) WithId(id string) *AsyncNorm {
	s.id = id
	return s
}

func (s *AsyncNorm) WithHandler(mrh sqlt.MultiRowsHandler) *AsyncNorm {
	s.mrh = mrh
	return s
}

func (s *AsyncNorm) AddParam(k string, v interface{}) *AsyncNorm {
	if k != "" && v != nil {
		s.param[k] = v
	}
	return s
}

func (s *AsyncNorm) ResetAll() *AsyncNorm {
	s.ctx = ctx{param: make(map[string]interface{})}
	return s
}

func (s *AsyncNorm) Reset() *AsyncNorm {
	s.ctx = ctx{param: s.param}
	return s
}

func (s *AsyncNorm) Query() *AsyncNorm {
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

func (s *AsyncNorm) Exec() *AsyncNorm {
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

func (s *AsyncNorm) ExecReturning() *AsyncNorm {
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

func (s *AsyncNorm) Wait() Collator {
	s.wg.Wait()
	return Collator{}
}

func (s *AsyncNorm) Cancel() {
	s.wg.Wait()
}

func NewAsyncNorm(op *sqlt.DbOp) *AsyncNorm {
	return &AsyncNorm{op: op, ctx: ctx{param: make(map[string]interface{})}}
}
