package model

import (
	"app/shared/database"
	// "app/shared/ordenate"
	sq "github.com/Masterminds/squirrel"
)

const (
	GteAmount = "0"
	GtYear = "0"
	GteMonth = "1"
	LteMonth = "12"
)

type Invoice struct {
	Amount float64 `db:"amount" validate:"gte=0"`
	Document string `db:"document" validate:"required"`
	Month int `db:"month" validate:"gte=1,lte=12"`
	Year int `db:"year" validate:"gt=1"`
	Is_active bool `db:"is_active validate:"required"`
}

func InvoiceCreate(amount float64, document string, month int, year int) error {
	var err error
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query, _, _ := psql.Insert("public.\"Invoice\"").Columns("amount", "document", "month", "year","is_active").Values(amount, document, month, year, true).ToSql()
	_, err = database.SQL.Exec(query, amount, document, month, year, true)
	return err
}
