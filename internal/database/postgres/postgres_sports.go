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
		CREATE TABLE IF NOT EXISTS sports (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL UNIQUE
		)
	`); err != nil {
		panic("Error creating table 'sports'")
	}

	// add sports to table
	sports := []string{
		"baseball", "basketball", "football", "hockey", "quidditch",
		"rugby", "soccer", "softball", "volleyball",
	}
	for _, sport := range sports {
		sportID := token.SignedToken(4)
		if _, err = p.DB.Exec(`
			INSERT INTO sports (id, name)
			VALUES ($1, $2)
			ON CONFLICT (name)
			DO NOTHING
		`, sportID[:len(sportID)-1], sport); err != nil {
			fmt.Println("Error populating table 'sports'")
		}
	}
}

func (p Postgres) ListSports() ([]model.Sport, error) {
	sports := []model.Sport{}

	rows, err := p.DB.Query(`SELECT * FROM sports`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var sport model.Sport
		if err := rows.Scan(
			&sport.ID,
			&sport.Name,
		); err != nil {
			return nil, err
		}
		sport.ID = token.ReturnSignedToken(sport.ID)
		sports = append(sports, sport)
	}

	return sports, nil
}
