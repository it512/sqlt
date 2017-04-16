package nested

import (
	"github.com/it512/sqlt"
)

type (
	NestableSqlAssmbler interface {
		sqlt.SqlAssembler
		HasId(id string) bool
	}

	SqlAssemblerSet struct {
		def        sqlt.SqlAssembler
		assemblers []NestableSqlAssmbler
	}
)

func (n *SqlAssemblerSet) AssembleSql(id string, data interface{}) (sqlt.SqlDescriber, error) {
	for _, l := range n.assemblers {
		if l.HasId(id) {
			return l.AssembleSql(id, data)
		}
	}

	return n.def.AssembleSql(id, data)
}

func NewSqlAssemblerSet(def sqlt.SqlAssembler, assemblers ...NestableSqlAssmbler) *SqlAssemblerSet {
	return &SqlAssemblerSet{def: def, assemblers: assemblers}
}
