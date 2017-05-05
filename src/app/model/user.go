package model

import (
	"app/shared/database"
	"app/shared/ordenate"
	"errors"
	"strconv"

	sq "github.com/Masterminds/squirrel"
)

// constantes para mensagem de erro na validacao
const (
	GtName     = "6"
	GtPassword = "6"
)

// struct do usuario
type User struct {
	Name     string `db:"name" validate:"required,alphanum,gt=6" json:"name"`
	Password string `db:"password" validate:"required,gt=6,excludesall= \n\t" json:"password"`
}

// funcao que adiciona o usuario no banco
// caso ocorra algum erro ele retorna uma mensagem

func UserCreate(name string, password string) error {
	var err error
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query, _, _ := psql.Insert("public.\"User\"").Columns("name", "password").Values(name, password).ToSql()
	_, err = database.SQL.Exec(query, name, password)
	return err
}

// funcao que busca o array de usuarios pela
// recebe parametros de order, page, length e name para filtrar os dados
// a query é feita deinamicamente de acordo com os dados inseridos
// retorna uma lista de usuarios e uma mensagem de erro caso tenha

func UserGetAll(orders []ordenate.Ordenate, page string, length string, name string) ([]User, error) {
	users := []User{}
	err := errors.New("")
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Select("name").From("public.\"User\"")
	var queryStr string
	if len(orders) > 0 {
		for _, order := range orders {
			query = query.OrderBy(order.Column + " " + order.Order)
		}
	}
	offset, err1 := strconv.ParseUint(page, 10, 64)
	limit, err2 := strconv.ParseUint(length, 10, 64)
	if err1 != nil {
		offset = 0
	}
	if err2 != nil {
		query = query.Offset(offset)
	} else {
		query = query.Limit(limit).Offset(offset * limit)
	}
	if len(name) > 0 {
		queryStr, _, _ = query.Where("name ILIKE ? ", name).ToSql()
		err = database.SQL.Select(&users, queryStr, "%"+name+"%")
	} else {
		queryStr, _, _ = query.ToSql()
		err = database.SQL.Select(&users, queryStr)
	}
	if err != nil {
		return []User{}, err
	} else {
		return users, nil
	}
}

// funcao que realiza a busca do usuário pelo name
// usada para a validacao e autenticacao

func UserByName(name string) (User, error) {
	user := User{}
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, _, _ := psql.Select("*").From("public.\"User\"").Where(sq.Eq{"name": name}).ToSql()
	err := database.SQL.Get(&user, query, name)
	if err != nil {
		return User{}, err
	} else {
		return user, nil
	}

}
