package database

import (
	"github.com/jmoiron/sqlx"
	"log"
	_ "github.com/lib/pq"
)

var (
	database Database
	SQL *sqlx.DB
)

type Database struct {
	Username  string
	Password  string
	Dbname    string
	Hostname  string
	Port      string
	Parameter string
}

func stringConnect(db Database) string {
	return "postgres://" + db.Username + ":" + db.Password + "@" + db.Hostname + ":" + db.Port + "/" + db.Dbname + db.Parameter
}

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

func readConfig() Database {
	return database
}
