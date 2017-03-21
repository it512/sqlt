package sqlt

import (
	"bytes"
	"text/template"

	"strings"

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
	var sqlType int = SQL_TYPE_NORMAL

	if e != nil {
		sql = buffer.String()
		sqlType = determineSqlType(&sql)
		if l.l.IsDebugEnable() {
			l.l.Debugln(sql, sqlType, param)
		}
	}

	return sql, sqlType, e
}

func determineSqlType(sql *string) int {
	if strings.ContainsAny(*sql, "readonly") {
		return SQL_TYPE_READONLY
	}

	return SQL_TYPE_NORMAL
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
