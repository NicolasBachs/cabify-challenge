package mocks

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func MockDatabase(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()

	if err != nil {
		assert.Fail(t, "Cannot mock sqlx database")
	}

	dbx := sqlx.NewDb(db, "sqlmock")

	return dbx, mock
}
