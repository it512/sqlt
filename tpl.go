package sqlt

import (
	"bytes"
	"text/template"
)

type (
	tplSqlLoader struct {
		t *template.Template
	}
)

func (l *tplSqlLoader) LoadSql(id string, param interface{}) (string, int, error) {
	buffer := new(bytes.Buffer)
	if l.t.ExecuteTemplate(buffer, id, param) != nil {
		return "", 0, nil
	}

	//s := buffer.String()
	//strings.Contains(s, "select")

	return buffer.String(), 0, nil

}

func newTemplateSqlLoader(pattern string) *tplSqlLoader {
	tpl := template.Must(template.ParseGlob(pattern))
	return &tplSqlLoader{t: tpl}
}
