package postgres

import (
	"fmt"

	"github.com/leagueify/leagueify/internal/config"
	"github.com/leagueify/leagueify/internal/lib/token"
	"github.com/leagueify/leagueify/internal/model"
)

func init() {
	cfg := config.LoadConfig()
	database, err := Connect(cfg.DBConnStr)
	if err != nil {
		panic(fmt.Sprintf("Error connecting to database '%s'", err))
	}
	defer database.Close()
	p := Postgres{DB: database}

	// create database table
	if _, err := p.DB.Exec(`
		CREATE TABLE IF NOT EXISTS divisions (
			id TEXT PRIMARY KEY,
			name TEXT UNIQUE NOT NULL,
			min_age INTEGER NOT NULL,
			max_age INTEGER NOT NULL
		)
	`); err != nil {
		panic("Error creating table 'divisions'")
	}
}

func (p Postgres) CreateDivisions(payload model.DivisionCreation) error {
	tx, err := p.BeginTransaction()
	if err != nil {
		return err
	}

	for _, division := range payload.Divisions {
		if _, err := tx.Exec(`
			INSERT INTO divisions (id, name, min_age, max_age)
			VALUES ($1, $2, $3, $4)
		`,
			division.ID[:len(division.ID)-1], division.Name,
			division.Age.Min, division.Age.Max,
		); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (p Postgres) ListDivisions() ([]model.Division, error) {
	divisions := []model.Division{}

	rows, err := p.DB.Query(`SELECT * FROM divisions`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var division model.Division
		if err := rows.Scan(
			&division.ID,
			&division.Name,
			&division.Age.Min,
			&division.Age.Max,
		); err != nil {
			return nil, err
		}
		division.ID = token.ReturnSignedToken(division.ID)
		divisions = append(divisions, division)
	}

	return divisions, nil
}
