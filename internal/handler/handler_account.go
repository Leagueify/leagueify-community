package handler

import (
	"github.com/leagueify/leagueify/internal/lib/auth"
	"github.com/leagueify/leagueify/internal/lib/date"
	"github.com/leagueify/leagueify/internal/lib/error"
	"github.com/leagueify/leagueify/internal/lib/token"
	"github.com/leagueify/leagueify/internal/model"
)

func (h *handler) ActivateAccount(accountID string) error {
	if !token.VerifyToken(accountID) {
		return &errors.LeagueifyError{Message: "invalid account ID"}
	}

	if err := h.db.ActivateAccount(accountID); err != nil {
		return err
	}

	return nil
}

func (h *handler) CreateAccount(payload model.AccountCreation) error {
	if !date.MeetsYearRequirement(18, payload.DateOfBirth, nil) {
		return &errors.LeagueifyError{
			Message: "must be 18 or older to create an account",
		}
	}

	if err := auth.HashPassword(&payload.Password); err != nil {
		return err
	}

	payload.ID = token.SignedToken(8)
	payload.IsAdmin = !h.db.HasActiveAccounts()
	payload.IsActive = !h.db.HasActiveEmailConfig()

	if err := h.db.CreateAccount(payload); err != nil {
		return err
	}

	return nil
}

func (h *handler) LoginAccount(payload model.AccountCredentials, aud string) (string, error) {
	account, err := h.db.GetAccountByEmail(payload.Email)
	if err != nil {
		return "", err
	}

	if !account.IsActive {
		return "", &errors.LeagueifyError{
			Message: "inactive account",
		}
	}

	if !auth.PasswordsMatch(payload.Password, account.Password) {
		return "", &errors.LeagueifyError{
			Message: "invalid account credentials",
		}
	}

	authToken, err := auth.CreateJWT(account, aud, 5)
	if err != nil {
		return "", err
	}

	return authToken, nil
}
