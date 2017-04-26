package model

import (
	"strconv"
	"app/shared/database"
	"app/shared/ordenate"
	sq "github.com/Masterminds/squirrel"
)

const (
	GtName = "6"
	GtPassword = "6"
)

type User struct {
	Name     string `db:"name" validate:"required,alphanum,gt=6"`
	Password string `db:"password" validate:"required,gt=6,excludesall= \n\t"`
}


func UserCreate(name string, password string) error {
	var err error
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query, _, _ := psql.Insert("public.\"User\"").Columns("name", "password").Values(name, password).ToSql()
	_, err = database.SQL.Exec(query, name, password)
	return err
}

func UserGetAll( orders []ordenate.Ordenate, page string, length string) ([]User, error) {
	users := []User{}
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query:= psql.Select("name").From("public.\"User\"")
	var queryStr string
	if len(orders) > 0 {
		for _, order := range orders {
			query = query.OrderBy(order.Column + " " +  order.Order)
		}
	}
	offset, err1 := strconv.ParseUint(page, 10, 64)
	limit, err2 := strconv.ParseUint(length, 10, 64)
	if err1 != nil {
		offset = 0
	}
	if err2 != nil {
		queryStr, _, _ = query.Offset(offset).ToSql()
	} else {
		queryStr, _, _ = query.Limit(limit).Offset(offset * limit).ToSql()
	}
	err := database.SQL.Select(&users, queryStr)
	if err != nil {
		return []User{}, err
	} else {
		return users, nil
	}
}
