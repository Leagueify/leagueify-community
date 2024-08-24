package handler

import (
	"github.com/leagueify/leagueify/internal/lib/error"
	"github.com/leagueify/leagueify/internal/model"
)

func (h *handler) ListSports() ([]model.Sport, error) {
	sports, err := h.db.ListSports()
	if err != nil {
		return nil, err
	}

	if len(sports) == 0 {
		return nil, &errors.LeagueifyError{Message: "no sports"}
	}

	return sports, nil
}
