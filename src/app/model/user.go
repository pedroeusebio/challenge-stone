package model

import (
	"app/shared/database"
	sq "github.com/Masterminds/squirrel"
)

type User struct {
	name     string `db:"name"`
	password string `db:"password"`
}


func UserCreate(name string, password string) error {
	var err error
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query, _, _ := psql.Insert("public.\"User\"").Columns("name", "password").Values(name, password).ToSql()
	_, err = database.SQL.Exec(query, name, password)
	return err
}

