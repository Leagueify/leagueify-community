package postgres

import (
	"fmt"

	"github.com/leagueify/leagueify/internal/config"
	"github.com/leagueify/leagueify/internal/lib/error"
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
		CREATE TABLE IF NOT EXISTS accounts (
			id TEXT PRIMARY KEY,
			first_name TEXT NOT NULL,
			last_name TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			phone TEXT NOT NULL UNIQUE,
			date_of_birth TEXT NOT NULL,
			player_ids TEXT[] NOT NULL,
			coach BOOLEAN DEFAULT false,
			volunteer BOOLEAN DEFAULT false,
			is_active BOOLEAN DEFAULT false,
			is_admin BOOLEAN DEFAULT false
		)
	`); err != nil {
		panic("Error creating table 'accounts'")
	}
}

func (p Postgres) ActivateAccount(accountID string) error {
	results, err := p.DB.Exec(`
		UPDATE accounts SET is_active = true
		WHERE id = $1 AND is_active = false
	`, accountID[:len(accountID)-1])
	if err != nil {
		return err
	}

	rows, err := results.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return &errors.LeagueifyError{Message: "account update failed"}
	}

	return nil
}

func (p Postgres) CreateAccount(account model.AccountCreation) error {
	if _, err := p.DB.Exec(`
		INSERT INTO accounts (
			id, first_name, last_name, email, password, phone,
			date_of_birth, player_ids, coach, volunteer, is_active,
			is_admin
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
		)`,
		account.ID[:len(account.ID)-1], account.FirstName,
		account.LastName, account.Email, account.Password,
		account.Phone, account.DateOfBirth, "{}", account.Coach,
		account.Volunteer, account.IsActive, account.IsAdmin,
	); err != nil {
		return err
	}

	return nil
}

func (p Postgres) GetAccountByEmail(email string) (model.Account, error) {
	var account model.Account

	if err := p.DB.QueryRow(`
		SELECT * FROM accounts WHERE email = $1
	`, email).Scan(
		&account.ID, &account.FirstName, &account.LastName,
		&account.Email, &account.Password, &account.Phone,
		&account.DateOfBirth, &account.Players, &account.Coach,
		&account.Volunteer, &account.IsActive, &account.IsAdmin,
	); err != nil {
		return account, err
	}

	return account, nil
}

func (p Postgres) HasActiveAccounts() bool {
	var totalAccounts int

	row := p.DB.QueryRow(`
			SELECT COUNT(*) FROM accounts WHERE is_active = true
		`,
	)
	if err := row.Scan(&totalAccounts); err != nil {
		return false
	}

	if totalAccounts != 0 {
		return true
	}

	return false
}
