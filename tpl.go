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
	if l.t.ExecuteTemplate(buffer, id, param) != nil {
		return "", 0, nil
	}

	return buffer.String(), 0, nil
}

func NewDefaultSqlLoader(pattern string) *DefaultSqlLoader {
	tpl := template.Must(template.ParseGlob(pattern))
	return &DefaultSqlLoader{t: tpl}
}
