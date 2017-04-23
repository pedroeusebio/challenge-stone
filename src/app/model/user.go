package model

import (
	"app/shared/database"
	sq "github.com/Masterminds/squirrel"
)

const (
	GtName = "6"
	GtPassword = "6"
)

type User struct {
	Name     string `db:"name" validate:"required,alphanum,gt=6"`
	Password string `db:"password"validate:"required,gt=6,excludesall= \n\t"`
}


func UserCreate(name string, password string) error {
	var err error
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query, _, _ := psql.Insert("public.\"User\"").Columns("name", "password").Values(name, password).ToSql()
	_, err = database.SQL.Exec(query, name, password)
	return err
}

