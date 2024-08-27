package handler

import (
	"fmt"

	"github.com/leagueify/leagueify/internal/lib/error"
	"github.com/leagueify/leagueify/internal/lib/token"
	"github.com/leagueify/leagueify/internal/model"
)

func (h *handler) CreateDivisions(payload model.DivisionCreation) error {
	if !h.db.HasExistingLeague() {
		return &errors.LeagueifyError{Message: "active league required"}
	}

	var ageErrors []string
	for i, division := range payload.Divisions {
		payload.Divisions[i].ID = token.SignedToken(6)

		if division.Age.Min >= division.Age.Max {
			ageErrors = append(ageErrors, division.Name)
		}
	}
	if len(ageErrors) != 0 {
		return &errors.LeagueifyError{
			Message: fmt.Sprintf(
				"invalid age range(s): %v", ageErrors,
			),
		}
	}

	if err := h.db.CreateDivisions(payload); err != nil {
		return err
	}

	return nil
}

func (h *handler) ListDivisions() ([]model.Division, error) {
	divisions, err := h.db.ListDivisions()
	if err != nil {
		return nil, err
	}

	if len(divisions) == 0 {
		return nil, &errors.LeagueifyError{Message: "no divisions"}
	}

	return divisions, nil
}
