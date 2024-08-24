package postgres

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Postgres struct {
	DB *sql.DB
}

func Connect(DBConnStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", DBConnStr)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (p Postgres) BeginTransaction() (*sql.Tx, error) {
	tx, err := p.DB.Begin()
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (p Postgres) Close() error {
	if err := p.DB.Close(); err != nil {
		return err
	}
	return nil
}
