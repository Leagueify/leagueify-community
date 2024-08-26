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
		CREATE TABLE IF NOT EXISTS seasons (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			season_start TEXT NOT NULL,
			season_end TEXT NOT NULL,
			registration_start TEXT NOT NULL,
			registration_end TEXT NOT NULL
		)
	`); err != nil {
		panic("Error creating table 'seasons'")
	}
}

func (p Postgres) CreateSeason(payload model.Season) error {
	if _, err := p.DB.Exec(`
		INSERT INTO seasons (
			id, name, season_start, season_end, registration_start,
			registration_end
		)
		VALUES (
			$1, $2, $3, $4, $5, $6
		)`,
		payload.ID[:len(payload.ID)-1], payload.Name,
		payload.Season.Start, payload.Season.End,
		payload.Registration.Start, payload.Registration.End,
	); err != nil {
		return err
	}

	return nil
}

func (p Postgres) GetSeasonByID(seasonID string) (model.Season, error) {
	var season model.Season

	if err := p.DB.QueryRow(`
		SELECT * FROM seasons WHERE id = $1
	`, seasonID[:len(seasonID)-1]).Scan(
		&season.ID, &season.Name, &season.Season.Start,
		&season.Season.End, &season.Registration.Start,
		&season.Registration.End,
	); err != nil {
		return season, err
	}

	season.ID = token.ReturnSignedToken(season.ID)
	return season, nil
}

func (p Postgres) ListSeasons() ([]model.SeasonList, error) {
	seasons := []model.SeasonList{}

	rows, err := p.DB.Query(`SELECT id, name FROM seasons`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var season model.SeasonList
		if err := rows.Scan(
			&season.ID,
			&season.Name,
		); err != nil {
			return nil, err
		}
		season.ID = token.ReturnSignedToken(season.ID)
		seasons = append(seasons, season)
	}

	return seasons, nil

}

func (p Postgres) UpdateSeason(seasonID string, payload model.Season) error {
	if _, err := p.DB.Exec(`
		UPDATE SEASONS
		SET name = $1, season_start = $2, season_end = $3,
			registration_start = $4, registration_end = $5
		WHERE id = $6
	`,
		payload.Name, payload.Season.Start, payload.Season.End,
		payload.Registration.Start, payload.Registration.End,
		seasonID[:len(seasonID)-1],
	); err != nil {
		return err
	}

	return nil
}
