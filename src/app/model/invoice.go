package model

import (
	"app/shared/database"
	"app/shared/ordenate"
	"strconv"

	sq "github.com/Masterminds/squirrel"
)

// constantes para retornar mensagem de erro na validacao
const (
	GteAmount   = "0"
	GteMonth    = "1"
	GteDocument = "11"
	GtYear      = "0"
	LteMonth    = "12"
	LteDocument = "14"
)

//struct do invoice
type Invoice struct {
	Id        string  `db:"id" json:"id"`
	Amount    float64 `db:"amount" validate:"gte=0" json:"amount"`
	Document  string  `db:"document" validate:"required,numeric,len=11|len=14" json:"document"`
	Month     int     `db:"month" validate:"gte=1,lte=12" json:"month"`
	Year      int     `db:"year" validate:"gt=1" json:"year"`
	Is_active bool    `db:"is_active" validate:"required" json:"is_active"`
}

// funcao para adicionar o invoice no banco de dados
// retorna erro ou nil caso nao exista nenhum erro

func InvoiceCreate(amount float64, document string, month int, year int) error {
	var err error
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query, _, _ := psql.Insert("public.\"Invoice\"").Columns("amount", "document", "month", "year", "is_active").Values(amount, document, month, year, true).ToSql()
	_, err = database.SQL.Exec(query, amount, document, month, year, true)
	return err
}

// funcao que retornar um array de invoice e erro caso ocorra.
// a query é feita dinamicamente de acordo com os parametros passados

func InvoiceGetAll(orders []ordenate.Ordenate, page string, length string) ([]Invoice, error) {
	invoices := []Invoice{}
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Select("*").From("public.\"Invoice\"")
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
		queryStr, _, _ = query.Offset(offset).ToSql()
	} else {
		queryStr, _, _ = query.Limit(limit).Offset(offset * limit).ToSql()
	}
	err := database.SQL.Select(&invoices, queryStr)
	if err != nil {
		return []Invoice{}, err
	} else {
		return invoices, nil
	}
}

// funcao de remocao logica dos invoice
// verifica a existencia do dado no banco e retorna erro caso nao exista
// sucesso : retorna o invoice deletado e nil

func InvoiceDelete(id string) (Invoice, error) {
	var err error
	var invoice Invoice
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	qGet, _, _ := psql.Select("*").From("public.\"Invoice\"").Where(sq.Eq{"id": id}).ToSql()
	err1 := database.SQL.Get(&invoice, qGet, id)
	if err1 != nil {
		return invoice, err1
	}
	query, _, _ := psql.Update("public.\"Invoice\"").Set("is_active", false).Where(sq.Eq{"id": id}).ToSql()
	_, err = database.SQL.Exec(query, false, id)
	return invoice, err
}
