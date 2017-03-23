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

var (
	sqlTypeMap = make(map[string]int)
)

func (l *DefaultSqlLoader) LoadSql(id string, param interface{}) (string, int, error) {
	buffer := new(bytes.Buffer)
	e := l.t.ExecuteTemplate(buffer, id, param)

	var sql string
	var sqlType int = SQL_TYPE_NORMAL

	if e == nil {
		sql = buffer.String()
		sqlType = determineSqlType(id)
		if l.l.IsDebugEnable() {
			l.l.Debugln(sql, sqlType, param)
		}
	}

	return sql, sqlType, e
}

func determineSqlType(id string) int {
	if i, ok := sqlTypeMap[id]; ok {
		return i
	}

	var t int = SQL_TYPE_NORMAL

	if strings.ContainsAny(id, "readonly") {
		t = SQL_TYPE_READONLY
	}

	sqlTypeMap[id] = t
	return t
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
