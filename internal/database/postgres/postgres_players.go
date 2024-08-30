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
		CREATE TABLE IF NOT EXISTS players (
			id TEXT PRIMARY KEY,
			parent_id TEXT NOT NULL,
			first_name TEXT NOT NULL,
			last_name TEXT NOT NULL,
			date_of_birth TEXT NOT NULL,
			division_id TEXT,
			team_id TEXT,
			is_registered BOOLEAN DEFAULT false,
			hash TEXT UNIQUE NOT NULL
		)
	`); err != nil {
		panic("Error creating table 'players'")
	}
}

func (p Postgres) CreatePlayers(payload model.Players) error {
	tx, err := p.BeginTransaction()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, player := range payload.Players {
		if _, err := tx.Exec(`
			INSERT INTO players (
				id, parent_id, first_name, last_name,
				date_of_birth, hash
			)
			VALUES ($1, $2, $3, $4, $5, $6)
		`,
			player.ID[:len(player.ID)-1],
			player.ParentID[:len(player.ParentID)-1],
			player.FirstName, player.LastName, player.DateOfBirth,
			player.Hash,
		); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (p Postgres) DeletePlayer(accountID, playerID string) error {
	if _, err := p.DB.Exec(`
		DELETE FROM players WHERE id = $1 AND parent_id = $2
	`, playerID[:len(playerID)-1], accountID[:len(accountID)-1],
	); err != nil {
		return err
	}

	return nil
}

func (p Postgres) GetPlayerByID(accountID, playerID string) (model.Player, error) {
	var player model.Player

	if err := p.DB.QueryRow(`
		SELECT * FROM players WHERE id = $1 AND parent_id = $2
	`, playerID[:len(playerID)-1], accountID[:len(accountID)-1]).Scan(
		&player.ID, &player.ParentID, &player.FirstName,
		&player.LastName, &player.DateOfBirth, &player.DivisionID,
		&player.TeamID, &player.IsRegistered, &player.Hash,
	); err != nil {
		return model.Player{}, err
	}

	player.ID = token.ReturnSignedToken(player.ID)
	player.ParentID = token.ReturnSignedToken(player.ParentID)
	return player, nil
}

func (p Postgres) IsExistingPlayer(hash string) bool {
	var players int

	if err := p.DB.QueryRow(`
		SELECT COUNT(*) FROM players WHERE hash = $1
	`, hash).Scan(&players); err != nil {
		return false
	}

	if players != 0 {
		return true
	}

	return false
}

func (p Postgres) ListPlayers(accountID string) ([]model.PlayerList, error) {
	players := []model.PlayerList{}

	rows, err := p.DB.Query(`
		SELECT id, first_name, last_name
		FROM players WHERE parent_id = $1
	`, accountID[:len(accountID)-1])
	if err != nil {
		return players, err
	}
	defer rows.Close()
	for rows.Next() {
		var player model.PlayerList
		if err := rows.Scan(
			&player.ID,
			&player.FirstName,
			&player.LastName,
		); err != nil {
			return players, err
		}
		player.ID = token.ReturnSignedToken(player.ID)
		players = append(players, player)
	}

	return players, nil
}
