package sqlt

type (
	Param struct {
		m map[string]interface{}
	}
)

func (p *Param) Add(k string, v interface{}) *Param {
	if k == "" || v == nil {
		return p
	}

	p.m[k] = v

	return p
}

func (p *Param) Remove(k string) *Param {
	if k != "" {
		delete(p.m, k)
	}
	return p
}

func (p *Param) Get(k string) interface{} {
	if k != "" {
		return p.m[k]
	}

	return nil
}

func (p *Param) Map() map[string]interface{} {
	return p.m
}

func NewParam() *Param {
	return &Param{m: make(map[string]interface{})}
}
