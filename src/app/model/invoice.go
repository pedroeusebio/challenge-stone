package model

import (
	"fmt"
	"strconv"
	"app/shared/database"
	"app/shared/ordenate"
	sq "github.com/Masterminds/squirrel"
)

const (
	GteAmount = "0"
	GtYear = "0"
	GteMonth = "1"
	LteMonth = "12"
)

type Invoice struct {
	Id string `db:"id"`
	Amount float64 `db:"amount" validate:"gte=0"`
	Document string `db:"document" validate:"required"`
	Month int `db:"month" validate:"gte=1,lte=12"`
	Year int `db:"year" validate:"gt=1"`
	Is_active bool `db:"is_active" validate:"required"`
}


func InvoiceCreate(amount float64, document string, month int, year int) error {
	var err error
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query, _, _ := psql.Insert("public.\"Invoice\"").Columns("amount", "document", "month", "year","is_active").Values(amount, document, month, year, true).ToSql()
	_, err = database.SQL.Exec(query, amount, document, month, year, true)
	return err
}

func InvoiceGetAll( orders []ordenate.Ordenate, page string, length string) ([]Invoice, error) {
	invoices := []Invoice{}
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query:= psql.Select("*").From("public.\"Invoice\"")
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
	err := database.SQL.Select(&invoices, queryStr)
	if err != nil {
		return []Invoice{}, err
	} else {
		return invoices, nil
	}
}

func InvoiceDelete(id string) (Invoice, error) {
	var err error
	var invoice Invoice
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	qGet, _, _ := psql.Select("*").From("public.\"Invoice\"").Where(sq.Eq{"id": id}).ToSql()
	fmt.Println(qGet)
	err1 := database.SQL.Get(&invoice, qGet, id)
	fmt.Println(invoice)
	if err1 != nil {
		fmt.Println(err1)
		return invoice, err1
	}
	query, _, _:= psql.Update("public.\"Invoice\"").Set("is_active", false).Where(sq.Eq{"id": id}).ToSql()
	fmt.Println(query)
	_, err = database.SQL.Exec(query, false, id)
	return invoice, err
}
