{{define "select.student"}}
	select
		*
	from
		student
	where
		{{if .name}}
		name = :name
		{{end}}
{{end}}


{{define "insert.student.returning"}}
	insert into student
	(
		id,
		name,
		age,
		sex
	)
	values
	(
		:id,
		:name,
		:age,
		:sex
	)
	returning
		id,
		name
{{end}}
