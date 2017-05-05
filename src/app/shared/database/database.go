package database

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	database Database
	SQL      *sqlx.DB
)

// struct de configuracao do banco de dados
type Database struct {
	Username  string
	Password  string
	Dbname    string
	Hostname  string
	Port      string
	Parameter string
}

// funcao para a formacao da string de conexao do banco

func stringConnect(db Database) string {
	return "postgres://" + db.Username + ":" + db.Password + "@" + db.Hostname + ":" + db.Port + "/" + db.Dbname + db.Parameter
}

// funcao para conectar ao banco de dados

func Connect(db Database) {
	var err error

	database = db

	if SQL, err = sqlx.Open("postgres", stringConnect(db)); err != nil {
		log.Println("Postgres DRIVER ERROR", err)
	}
	if err = SQL.Ping(); err != nil {
		log.Println("Database Error", err)
	}
}

// funcao para retornar os dados usados no banco de dados

func readConfig() Database {
	return database
}
