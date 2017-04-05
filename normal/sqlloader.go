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
		R     SqlRender
		L     log.Logger
		Debug bool
	}
)

func (l *NormalSqlLoader) LoadSql(id string, param interface{}) (sqlt.MappedSql, error) {
	mappedSql := new(NormalMappedSql)
	mappedSql.Id = id
	mappedSql.Param = param

	if l.Debug {
		if r, ok := l.R.(reloader); ok {
			r.Reload()
		}
	}

	e := l.R.Render(mappedSql, id, param)

	if l.L.IsDebugEnable() && e == nil {
		if sql, err := mappedSql.GetSql(); err == nil {
			l.L.Debugln(sql, param)
		}
	}

	return mappedSql, e
}

func NewNormalSqlLoader(pattern string) *NormalSqlLoader {
	r := NewStandardTemplateRender(pattern)
	logger := log.GetLogger("sqlt-default-loader")
	return &NormalSqlLoader{R: r, L: logger, Debug: false}
}
