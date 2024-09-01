package handler

import (
	"fmt"

	"github.com/leagueify/leagueify/internal/lib/date"
	"github.com/leagueify/leagueify/internal/lib/error"
	"github.com/leagueify/leagueify/internal/lib/token"
	"github.com/leagueify/leagueify/internal/model"
)

func (h *handler) CreateSeason(payload model.Season) error {
	if !h.db.HasExistingLeague() {
		return &errors.LeagueifyError{
			Message: "active league required",
		}
	}

	var dateErrors []string
	if !date.ValidDateRange(payload.Season) {
		dateErrors = append(dateErrors, "season")
	}
	if !date.ValidDateRange(payload.Registration) {
		dateErrors = append(dateErrors, "registration")
	}
	if len(dateErrors) != 0 {
		return &errors.LeagueifyError{
			Message: fmt.Sprintf(
				"invalid date range(s): %v", dateErrors,
			),
		}
	}

	payload.ID = token.SignedToken(6)
	if err := h.db.CreateSeason(payload); err != nil {
		return err
	}

	return nil
}

func (h *handler) GetSeasonByID(seasonID string) (model.Season, error) {
	if !token.VerifyToken(seasonID) {
		return model.Season{}, &errors.LeagueifyError{
			Message: "invalid seasonID",
		}
	}

	season, err := h.db.GetSeasonByID(seasonID)
	if err != nil {
		return model.Season{}, err
	}

	return season, nil
}

func (h *handler) ListSeasons() ([]model.SeasonList, error) {
	seasons, err := h.db.ListSeasons()
	if err != nil {
		return nil, err
	}

	if len(seasons) == 0 {
		return nil, &errors.LeagueifyError{
			Message: "no seasons",
		}
	}

	return seasons, nil
}

func (h *handler) UpdateSeason(seasonID string, payload model.SeasonUpdate) error {
	if !token.VerifyToken(seasonID) {
		return &errors.LeagueifyError{Message: "invalid seasonID"}
	}

	season, err := h.db.GetSeasonByID(seasonID)
	if err != nil {
		return err
	}

	if payload.Name != "" {
		season.Name = payload.Name
	}
	seasonDates := model.SeasonDates{}
	var dateErrors []string
	if payload.Season != seasonDates {
		if !date.ValidDateRange(payload.Season) {
			dateErrors = append(dateErrors, "season")
		}
		season.Season = payload.Season
	}
	if payload.Registration != seasonDates {
		if !date.ValidDateRange(payload.Registration) {
			dateErrors = append(dateErrors, "registration")
		}
		season.Registration = payload.Registration
	}
	if len(dateErrors) != 0 {
		return &errors.LeagueifyError{
			Message: fmt.Sprintf(
				"invalid date range(s): %v", dateErrors,
			),
		}
	}

	if err := h.db.UpdateSeason(seasonID, season); err != nil {
		return err
	}

	return nil
}
