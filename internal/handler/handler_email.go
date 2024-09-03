package handler

import (
	"crypto/tls"
	"fmt"
	"net/smtp"

	"github.com/leagueify/leagueify/internal/lib/error"
	"github.com/leagueify/leagueify/internal/lib/token"
	"github.com/leagueify/leagueify/internal/model"
)

func (h *handler) CreateEmailConfig(payload model.EmailConfig) error {
	if !h.db.HasExistingLeague() {
		return &errors.LeagueifyError{
			Message: "active league required",
		}
	}

	if h.db.HasExistingEmailConfig() {
		return &errors.LeagueifyError{Message: "existing email config"}
	}

	if err := h.ValidateCredentials(payload); err != nil {
		return err
	}

	payload.ID = token.SignedToken(6)
	payload.IsActive = true
	if err := h.db.CreateEmailConfig(payload); err != nil {
		return err
	}

	return nil
}

func (h *handler) GetEmailConfig() (model.EmailConfig, error) {
	emailConfig, err := h.db.GetEmailConfig()
	if err != nil {
		return model.EmailConfig{}, err
	}

	return emailConfig, nil
}

func (h *handler) UpdateEmailConfig(emailConfigID string, payload model.UpdateEmailConfig) error {
	emailConfig, err := h.GetEmailConfig()
	if err != nil {
		return err
	}

	if payload.OutboundEmail != "" {
		emailConfig.OutboundEmail = payload.OutboundEmail
	}

	if payload.SMTPHost != "" {
		emailConfig.SMTPHost = payload.SMTPHost
	}

	if payload.SMTPPort != 0 {
		emailConfig.SMTPPort = payload.SMTPPort
	}

	if payload.SMTPUser != "" {
		emailConfig.SMTPUser = payload.SMTPUser
	}

	if payload.SMTPPass != "" {
		emailConfig.SMTPPass = payload.SMTPPass
	}

	if payload.IsActive != emailConfig.IsActive {
		emailConfig.IsActive = payload.IsActive
	}

	if payload.IsActive {
		if err := h.ValidateCredentials(emailConfig); err != nil {
			return err
		}
	}

	if err := h.db.UpdateEmailConfig(emailConfigID, emailConfig); err != nil {
		return err
	}

	return nil
}

func (h *handler) ValidateCredentials(payload model.EmailConfig) error {
	tlsCfg := &tls.Config{
		InsecureSkipVerify: false, ServerName: payload.SMTPHost,
	}

	conn, err := tls.Dial(
		"tcp", fmt.Sprintf(
			"%s:%v", payload.SMTPHost, payload.SMTPPort,
		), tlsCfg,
	)
	if err != nil {
		return err
	}

	client, err := smtp.NewClient(conn, payload.SMTPHost)
	if err != nil {
		return err
	}

	auth := smtp.PlainAuth(
		"", payload.SMTPUser, payload.SMTPPass, payload.SMTPHost,
	)

	if err := client.Auth(auth); err != nil {
		return err
	}

	if err := client.Quit(); err != nil {
		return err
	}

	return nil
}
