package sqlt

import (
	"bytes"
	"log"
	"text/template"
)

type (
	DefaultSqlLoader struct {
		t *template.Template
		l *log.Logger
	}
)

func (l *DefaultSqlLoader) LoadSql(id string, param interface{}) (string, int, error) {
	buffer := new(bytes.Buffer)
	e := l.t.ExecuteTemplate(buffer, id, param)
	sql := buffer.String()

	if l.l != nil {
		if e != nil {
			l.l.Panicln(e)
		}
		l.l.Println(sql, param)
	}
	return sql, 0, e
}

func NewDefaultSqlLoader(pattern string) *DefaultSqlLoader {
	return NewSqlLoader(pattern, nil, make(template.FuncMap))
}

func NewSqlLoader(pattern string, l *log.Logger, funcMap template.FuncMap) *DefaultSqlLoader {
	tpl := template.Must(template.ParseGlob(pattern))
	tpl.Funcs(funcMap)
	return &DefaultSqlLoader{t: tpl, l: l}
}
