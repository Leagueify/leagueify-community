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
	// sport function
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
