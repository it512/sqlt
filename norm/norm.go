package norm

import "github.com/it512/sqlt"

type (
	CollateFunc func() error

	ctx struct {
		id    string
		param map[string]interface{}
		mrh   sqlt.MultiRowsHandler
	}

	Collator struct{}
)

func (c Collator) CollateWithFunc(cf CollateFunc) {
	if cf != nil {
		if e := cf(); e != nil {
			panic(e)
		}
	}
}
