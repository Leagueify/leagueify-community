package handler

import (
	"github.com/leagueify/leagueify/internal/lib/error"
	"github.com/leagueify/leagueify/internal/lib/token"
	"github.com/leagueify/leagueify/internal/model"
)

func (h *handler) CreateLeague(payload model.LeagueCreation) error {
	if !token.VerifyToken(payload.SportID) {
		return &errors.LeagueifyError{Message: "invalid sportID"}
	}

	if _, err := h.db.GetSportByID(payload.SportID); err != nil {
		return &errors.LeagueifyError{Message: "invalid sportID"}
	}

	if h.db.HasExistingLeague() {
		return &errors.LeagueifyError{Message: "league already exists"}
	}

	payload.ID = token.SignedToken(6)
	if err := h.db.CreateLeague(payload); err != nil {
		return err
	}

	return nil
}
