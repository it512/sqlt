# sqlt 使用说明
---
sqlt是一个模仿mybatis的go sqlmapping 实现，本质上，sqlt是数据库+模板，模板用于拼接sql（支持判断和自定义函数），底层采用[sqlx](https://github.com/jmoiron/sqlx)实现

## 一个简单的例子

go代码  

	db := sqlx.Open(...)
	loader := sqlt.NewDefaultSqlRender("path/*.tpl")
	dbop := sqlt.New(db, loader)

	param := make(map[string]interface{})
	param["w"] = ...
	param["s"] = ...
	result, e := dbop.Select("defget", param)

default模板

		{{define "defget"}}
			select * from {{.a}} where
			{{if and .w .s}}
				and {{.w}} like '%{{.s}}%'
			{{end}}

			{{if .id}}
				and id = :id
			{{end}}

			{{if .g}}
				group by {{.g}}
			{{end}}

			{{if .h}}
				having by {{.h}}
			{{end}}

			{{if .o}}
				order by {{.o}}
			{{end}}
		{{end}}
