package main

import (
	"context"
	"fmt"
	"log"

	"github.com/it512/dsds/simple"
	_ "github.com/it512/slf4go/simplelog"
	"github.com/it512/sqlt"
	"github.com/it512/sqlt/funcs"
	_ "github.com/lib/pq"
)

func main() {
	dbop := sqlt.NewSqlt(simple.NewSimpleDbSet("postgres", "dbname=test sslmode=disable"), sqlt.NewStdSqlAssemblerDefault("template/*.tpl"))
	smr := sqlt.NewSliceMapRowsHandler(funcs.Camal)

	param := make(map[string]interface{})
	param["name"] = "mike"
	e := dbop.QueryContext(context.Background(), "select.student", param, smr)

	if e != nil {
		log.Fatal(e)
	}

	for i := 0; i < smr.Count(); i++ {
		c := smr.ResuleSet(i)
		for _, r := range c {
			fmt.Printf("%s, %s\n", r["id"], r["name"])
		}
	}
}
