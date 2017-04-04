package normal

import (
	log "github.com/it512/slf4go"
	"github.com/it512/sqlt"
)

type (
	reloader interface {
		Reload()
	}

	NormalSqlLoader struct {
		T     SqlRender
		L     log.Logger
		Debug bool
	}
)

func (l *NormalSqlLoader) LoadSql(id string, param interface{}) (sqlt.MappedSql, error) {
	mappedSql := new(NormalMappedSql)
	mappedSql.Id = id
	mappedSql.Param = param

	if l.Debug {
		if r, ok := l.T.(reloader); ok {
			r.Reload()
		}
	}

	e := l.T.Render(mappedSql, id, param)

	if e == nil {
		if l.L.IsDebugEnable() {
			if sql, err := mappedSql.GetSql(); err == nil {
				l.L.Debugln(sql, param)
			}
		}
	}

	return mappedSql, e
}

func NewNormalSqlLoader(pattern string) *NormalSqlLoader {
	tr := NewStandardTemplateRender(pattern)
	logger := log.GetLogger("sqlt-default-loader")

	return &NormalSqlLoader{T: tr, L: logger, Debug: false}
}
