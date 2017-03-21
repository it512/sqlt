package sqlt

import (
	"bytes"
	"text/template"

	log "github.com/it512/slf4go"
)

type (
	DefaultSqlLoader struct {
		t *template.Template
		l log.Logger
	}
)

func (l *DefaultSqlLoader) LoadSql(id string, param interface{}) (string, int, error) {
	buffer := new(bytes.Buffer)
	e := l.t.ExecuteTemplate(buffer, id, param)

	var sql string
	if e != nil {
		sql = buffer.String()
		if l.l.IsDebugEnable() {
			l.l.Debugln(sql, param)
		}
	}

	return sql, 0, e
}

func NewDefaultSqlLoader(pattern string) *DefaultSqlLoader {
	return NewSqlLoader(pattern, make(template.FuncMap))
}

func NewSqlLoader(pattern string, funcMap template.FuncMap) *DefaultSqlLoader {
	tpl := template.Must(template.ParseGlob(pattern))
	tpl.Funcs(funcMap)

	logger := log.GetLogger("sqlt-default-loader")

	return &DefaultSqlLoader{t: tpl, l: logger}
}
