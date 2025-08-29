package database

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewPSQLConnection(source string) *sql.DB {
	db, err := sql.Open("pgx", source)
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}
	return db
}
