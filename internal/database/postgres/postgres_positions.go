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
		CREATE TABLE IF NOT EXISTS positions (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL
		)
	`); err != nil {
		panic("Error creating table 'positions'")
	}
}

func (p Postgres) CreatePositions(payload []model.Position) error {
	tx, err := p.BeginTransaction()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, position := range payload {
		if _, err := tx.Exec(`
			INSERT INTO positions (id, name) VALUES ($1, $2)
		`, position.ID[:len(position.ID)-1], position.Name,
		); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (p Postgres) HasExistingPositions() bool {
	var totalPositions int

	if err := p.DB.QueryRow(`
		SELECT COUNT(*) FROM positions
	`).Scan(&totalPositions); err != nil {
		return false
	}

	if totalPositions != 0 {
		return true
	}

	return false
}

func (p Postgres) ListPositions() ([]model.Position, error) {
	positions := []model.Position{}

	rows, err := p.DB.Query(`SELECT * FROM positions`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var position model.Position
		if err := rows.Scan(
			&position.ID,
			&position.Name,
		); err != nil {
			return nil, err
		}
		position.ID = token.ReturnSignedToken(position.ID)
		positions = append(positions, position)
	}

	return positions, nil
}
