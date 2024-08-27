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
	// league function
	CreateLeague(league model.LeagueCreation) error
	HasExistingLeague() bool
	// positin functions
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
