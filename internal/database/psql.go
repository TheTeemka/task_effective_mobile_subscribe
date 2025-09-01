package database

import (
	"database/sql"
	"log/slog"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewPSQLConnection(source string) *sql.DB {
	slog.Info("Connecting to PostgreSQL database")

	db, err := sql.Open("pgx", source)
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	slog.Info("Successfully connected to PostgreSQL database")
	return db
}
