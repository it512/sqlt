package normal

import (
	"bytes"
	"strings"
)

type (
	NormalMappedSql struct {
		Id    string
		Param interface{}
		bytes.Buffer
	}
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
	return strings.HasSuffix(s.Id, "-r") || strings.Contains(s.Id, "select")
}
