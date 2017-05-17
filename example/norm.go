package main

import (
	"fmt"

	"github.com/it512/dsds/simple"
	"github.com/it512/sqlt"
	"github.com/it512/sqlt/funcs"
	"github.com/it512/sqlt/norm"
	_ "github.com/lib/pq"
)

func main() {
	dbop := sqlt.NewSqlt(simple.NewSimpleDbSet("postgres", "dbname=test sslmode=disable"), sqlt.NewStdSqlAssemblerDefault("template/*.tpl"))
	n := norm.NewNorm(dbop)
	op := n.NewSimpleNormDefault()

	smrh := sqlt.NewSliceMapRowsHandler(funcs.Camal)

	_, err := op.
		WithId("select.student").
		AddParam("name", "mike").
		WithHandler(smrh).
		Query()

	if err == nil {
		for i := 0; i < smrh.Count(); i++ {
			c := smrh.ResuleSet(i)
			for _, r := range c {
				fmt.Printf("%s, %s\n", r["id"], r["name"])
			}
		}
	}
}
