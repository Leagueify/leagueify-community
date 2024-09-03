package database

import (
	"database/sql"
	"fmt"

	"github.com/leagueify/leagueify/internal/config"
	"github.com/leagueify/leagueify/internal/database/postgres"
	"github.com/leagueify/leagueify/internal/lib/error"
	"github.com/leagueify/leagueify/internal/model"
)

type Database interface {
	// account functions
	ActivateAccount(accountID string) error
	CreateAccount(account model.AccountCreation) error
	GetAccountByEmail(email string) (model.Account, error)
	HasActiveAccounts() bool
	// division functions
	CreateDivisions(payload model.DivisionCreation) error
	ListDivisions() ([]model.Division, error)
	// email functions
	CreateEmailConfig(payload model.EmailConfig) error
	GetEmailConfig() (model.EmailConfig, error)
	HasActiveEmailConfig() bool
	HasExistingEmailConfig() bool
	UpdateEmailConfig(emailConfigID string, payload model.EmailConfig) error
	// league function
	CreateLeague(league model.LeagueCreation) error
	HasExistingLeague() bool
	// player functions
	CreatePlayers(payload model.Players) error
	DeletePlayer(accountID, playerID string) error
	GetPlayerByID(accountID, playerID string) (model.Player, error)
	IsExistingPlayer(hash string) bool
	ListPlayers(accountID string) ([]model.PlayerList, error)
	// position functions
	CreatePositions(payload []model.Position) error
	HasExistingPositions() bool
	ListPositions() ([]model.Position, error)
	// season functions
	CreateSeason(season model.Season) error
	GetSeasonByID(sesonID string) (model.Season, error)
	ListSeasons() ([]model.SeasonList, error)
	UpdateSeason(seasonID string, payload model.Season) error
	// sport function
	GetSportByID(sportID string) (model.Sport, error)
	ListSports() ([]model.Sport, error)
	// database functions
	BeginTransaction() (*sql.Tx, error)
	Close() error
}

func GetDatabase() (Database, error) {
	cfg := config.LoadConfig()
	switch cfg.DB {
	case "postgres":
		db, err := postgres.Connect(cfg.DBConnStr)
		if err != nil {
			return nil, &errors.LeagueifyError{
				Message: fmt.Sprintf(
					"database connection error: '%s'", err,
				),
			}
		}
		return postgres.Postgres{
			DB: db,
		}, nil
	default:
		return nil, &errors.LeagueifyError{
			Message: fmt.Sprintf(
				"unsupported database '%s'", cfg.DB,
			),
		}
	}
}
