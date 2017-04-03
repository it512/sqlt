package sqlt

import (
	"bytes"
	"io"
	"strings"
	"text/template"

	log "github.com/it512/slf4go"
	"github.com/it512/sqlt"
)

type (
	NormalMappedSql struct {
		Id    string
		Param interface{}
		bytes.Buffer
	}
)

var (
	selectBytes []byte = []byte("select")
)

func (s *NormalMappedSql) GetSql() (string, error) {
	return s.String(), nil
}

func (s *NormalMappedSql) GetId() string {
	return s.Id
}

func (s *NormalMappedSql) GetParam() interface{} {
	return s.Param
}

func (s *NormalMappedSql) IsReadOnly() bool {
	return strings.HasSuffix(s.Id, "-r") || bytes.Contains(s.Bytes(), selectBytes)
}

type (
	TemplateRender interface {
		Render(w io.Writer, id string, param interface{}) error
	}

	StandardTemplateRender struct {
		t *template.Template
	}
)

func (str *StandardTemplateRender) Render(w io.Writer, id string, param interface{}) error {
	return str.t.ExecuteTemplate(w, id, param)
}

func NewStandardTemplateRender(pattern string, funcMap template.FuncMap) *StandardTemplateRender {
	tpl := template.Must(template.ParseGlob(pattern))
	tpl.Funcs(funcMap)

	return &StandardTemplateRender{t: tpl}
}

type (
	NormalSqlLoader struct {
		T TemplateRender
		L log.Logger
	}
)

func (l *NormalSqlLoader) LoadSql(id string, param interface{}) (sqlt.MappedSql, error) {
	mappedSql := new(NormalMappedSql)
	mappedSql.Id = id
	mappedSql.Param = param

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
	tr := NewStandardTemplateRender(pattern, make(template.FuncMap))
	logger := log.GetLogger("sqlt-default-loader")

	return &NormalSqlLoader{T: tr, L: logger}
}
