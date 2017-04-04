package normal

import (
	"io"
	"sync"
	"text/template"
)

type (
	SqlRender interface {
		Render(w io.Writer, id string, param interface{}) error
	}

	StandardTemplateRender struct {
		pattern string
		funcMap template.FuncMap
		t       *template.Template
		lock    *sync.RWMutex
	}
)

func (st *StandardTemplateRender) Render(w io.Writer, id string, param interface{}) error {
	st.lock.RLock()
	defer st.lock.RUnlock()
	return st.t.ExecuteTemplate(w, id, param)
}

func (st *StandardTemplateRender) Reload() {
	tpl := template.Must(template.ParseGlob(st.pattern))
	tpl.Funcs(st.funcMap)
	st.lock.Lock()
	st.t = tpl
	st.lock.Unlock()
}

func NewStandardTemplateRenderWithFunc(pattern string, funcMap template.FuncMap) *StandardTemplateRender {
	tpl := template.Must(template.ParseGlob(pattern))
	tpl.Funcs(funcMap)

	return &StandardTemplateRender{pattern: pattern, funcMap: funcMap, t: tpl, lock: new(sync.RWMutex)}
}

func NewStandardTemplateRender(pattern string) *StandardTemplateRender {
	return NewStandardTemplateRenderWithFunc(pattern, make(template.FuncMap))
}
