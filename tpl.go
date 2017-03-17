package sqlt

import (
	"bytes"
	"io"
	"text/template"
)

type (
	Renderer interface {
		Render(w io.Writer, id string, param interface{}) error
	}

	SqlRender struct {
		t Renderer
	}

	defualtRender *template.Template
)

//defualtRender
func (r *defualtRender) Render(w io.Writer, id string, param interface{}) error {
	tpl := r.(*template.Template)
	return tpl.ExecuteTemplate(w, id, param)
}

func (l *SqlRender) RenderSql(id string, param interface{}) (string, int, error) {
	buffer := new(bytes.Buffer)
	e := l.t.Render(buffer, id, param)
	return buffer.String(), 0, e
}

func NewDefaultSqlRender(pattern string) *SqlRender {
	tpl := template.Must(template.ParseGlob(pattern))
	return &SqlRender{t: defualtRender(tpl)}
}

func NewSqlRender(r Renderer) *SqlRender {
	return &SqlRender{t: r}
}
