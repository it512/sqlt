package sqlt

import (
	"bytes"
	"text/template"
)

type (
	DefaultSqlLoader struct {
		t *template.Template
	}
)

func (l *DefaultSqlLoader) LoadSql(id string, param interface{}) (string, int, error) {
	buffer := new(bytes.Buffer)
	e := l.t.ExecuteTemplate(buffer, id, param)
	return buffer.String(), 0, e
}

func NewDefaultSqlLoader(pattern string) *DefaultSqlLoader {
	tpl := template.Must(template.ParseGlob(pattern))
	return &DefaultSqlLoader{t: tpl}
}
