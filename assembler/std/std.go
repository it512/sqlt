package std

import (
	"bytes"
	"io"
	"sync"
	"text/template"

	log "github.com/it512/slf4go"
	"github.com/it512/sqlt"
)

type (
	SqlRender interface {
		Render(w io.Writer, id string, param interface{}) error
	}

	StdTemplateRender struct {
		pattern string
		funcMap template.FuncMap
		t       *template.Template
		lock    sync.RWMutex
	}
)

func (st *StdTemplateRender) Render(w io.Writer, id string, param interface{}) error {
	st.lock.RLock()
	defer st.lock.RUnlock()
	return st.t.ExecuteTemplate(w, id, param)
}

func (st *StdTemplateRender) Reload() {
	tpl := template.New("sqlt-std-template").Funcs(st.funcMap)
	tpl = template.Must(tpl.ParseGlob(st.pattern))
	st.lock.Lock()
	st.t = tpl
	st.lock.Unlock()
}

func NewStdTemplateRenderWithFuncs(pattern string, funcMap template.FuncMap) *StdTemplateRender {
	tpl := template.New("sqlt-std-template").Funcs(funcMap)
	tpl = template.Must(tpl.ParseGlob(pattern))
	return &StdTemplateRender{pattern: pattern, funcMap: funcMap, t: tpl}
}

func NewStdTemplateRender(pattern string) *StdTemplateRender {
	return NewStdTemplateRenderWithFuncs(pattern, make(template.FuncMap))
}

type (
	StdSqlDescriber struct {
		Id   string
		Data interface{}
		bytes.Buffer
	}
)

func (s *StdSqlDescriber) GetSql() (string, error) {
	return s.String(), nil
}

type (
	reloader interface {
		Reload()
	}

	StdSqlAssembler struct {
		Render SqlRender
		Logger log.Logger
		Debug  bool
	}
)

func (l *StdSqlAssembler) AssembleSql(id string, data interface{}) (sqlt.SqlDescriber, error) {
	desc := new(StdSqlDescriber)
	desc.Id = id
	desc.Data = data

	if l.Debug {
		if r, ok := l.Render.(reloader); ok {
			r.Reload()
		}
	}

	e := l.Render.Render(desc, id, data)

	if l.Logger.IsDebugEnable() && e == nil {
		if sql, err := desc.GetSql(); err == nil {
			l.Logger.Debugln(sql, data)
		}
	}

	return desc, e
}

func NewStdSqlAssemblerWithDefault(pattern string) *StdSqlAssembler {
	r := NewStdTemplateRender(pattern)
	logger := log.GetLogger("sqlt-default-loader")
	return &StdSqlAssembler{Render: r, Logger: logger, Debug: false}
}

func NewStdSqlAssembler(r SqlRender, debug bool) *StdSqlAssembler {
	logger := log.GetLogger("sqlt-default-loader")
	return &StdSqlAssembler{Render: r, Logger: logger, Debug: debug}
}
