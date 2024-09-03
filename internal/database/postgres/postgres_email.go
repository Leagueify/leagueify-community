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
		CREATE TABLE IF NOT EXISTS email (
			id TEXT PRIMARY KEY,
			outbound_email TEXT NOT NULL,
			smtp_host TEXT NOT NULL,
			smtp_port TEXT NOT NULL,
			smtp_user TEXT NOT NULL,
			smtp_pass TEXT NOT NULL,
			has_error BOOLEAN DEFAULT false,
			is_active BOOLEAN DEFAULT false
		)
	`); err != nil {
		panic("Error creating table 'email'")
	}
}

func (p Postgres) CreateEmailConfig(payload model.EmailConfig) error {
	if _, err := p.DB.Exec(`
		INSERT INTO email (
			id, outbound_email, smtp_host, smtp_port, smtp_user,
			smtp_pass, is_active
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`,
		payload.ID[:len(payload.ID)-1], payload.OutboundEmail,
		payload.SMTPHost, payload.SMTPPort, payload.SMTPUser,
		payload.SMTPPass, payload.IsActive,
	); err != nil {
		return err
	}

	return nil
}

func (p Postgres) GetEmailConfig() (model.EmailConfig, error) {
	emailConfig := model.EmailConfig{}

	if err := p.DB.QueryRow(`
		SELECT * FROM email LIMIT(1)
	`).Scan(
		&emailConfig.ID, &emailConfig.OutboundEmail,
		&emailConfig.SMTPHost, &emailConfig.SMTPPort,
		&emailConfig.SMTPUser, &emailConfig.SMTPPass,
		&emailConfig.HasError, &emailConfig.IsActive,
	); err != nil {
		return emailConfig, err
	}
	emailConfig.ID = token.ReturnSignedToken(emailConfig.ID)

	return emailConfig, nil
}

func (p Postgres) HasActiveEmailConfig() bool {
	var emailConfig int

	if err := p.DB.QueryRow(`
		SELECT COUNT(*) FROM email WHERE is_active = true
	`).Scan(&emailConfig); err != nil {
		return false
	}

	if emailConfig != 0 {
		return true
	}

	return false
}
func (p Postgres) HasExistingEmailConfig() bool {
	var emailConfig int

	if err := p.DB.QueryRow(`
		SELECT COUNT(*) FROM email
	`).Scan(&emailConfig); err != nil {
		return false
	}

	if emailConfig != 0 {
		return true
	}

	return false
}

func (p Postgres) UpdateEmailConfig(emailConfigID string, payload model.EmailConfig) error {
	if _, err := p.DB.Exec(`
		UPDATE email
		SET outbound_email = $1, smtp_host = $2, smtp_port = $3,
			smtp_user = $4, smtp_pass = $5, is_active = $6
		WHERE id = $7
	`,
		payload.OutboundEmail, payload.SMTPHost, payload.SMTPPort,
		payload.SMTPUser, payload.SMTPPass, payload.IsActive,
		emailConfigID[:len(emailConfigID)-1],
	); err != nil {
		return err
	}

	return nil
}
