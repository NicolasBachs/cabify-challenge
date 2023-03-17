package db

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresDBHandler struct {
	db *sqlx.DB
}

func NewPostgresDBHandler(connStr string) *PostgresDBHandler {
	h := &PostgresDBHandler{}

	err := h.Connect(connStr)

	if err != nil {
		log.Fatal("Error connecting to postgres, err: ", err)
	}

	return h
}

func (h *PostgresDBHandler) Connect(connStr string) error {
	db, err := sqlx.Connect("postgres", connStr)

	if err != nil {
		return err
	}

	h.db = db

	return nil
}

func (h *PostgresDBHandler) GetDB() *sqlx.DB {
	return h.db
}
