package postgres

import (
	"fmt"

	"github.com/leagueify/leagueify/internal/config"
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
		CREATE TABLE IF NOT EXISTS leagues (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			sport_id TEXT NOT NULL,
			master_admin TEXT NOT NULL
		)
	`); err != nil {
		panic("Error creating table 'accounts'")
	}

}

func (p Postgres) CreateLeague(league model.LeagueCreation) error {
	if _, err := p.DB.Exec(`
		INSERT INTO leagues (
			id, name, sport_id, master_admin
		)
		VALUES (
			$1, $2, $3, $4
		)`,
		league.ID[:len(league.ID)-1], league.Name,
		league.SportID[:len(league.SportID)-1],
		league.MasterAdmin[:len(league.MasterAdmin)-1],
	); err != nil {
		return err
	}

	return nil
}

func (p Postgres) HasExistingLeague() bool {
	var leagues int

	if err := p.DB.QueryRow(`
		SELECT COUNT(*) FROM leagues
	`).Scan(&leagues); err != nil {
		return false
	}

	if leagues != 0 {
		return true
	}

	return false
}
