package handler

import (
	"github.com/leagueify/leagueify/internal/lib/error"
	"github.com/leagueify/leagueify/internal/lib/token"
	"github.com/leagueify/leagueify/internal/model"
)

func (h *handler) CreatePositions(payload model.PositionCreation) error {
	if h.db.HasExistingPositions() {
		return &errors.LeagueifyError{
			Message: "positions already exist",
		}
	}

	var positionPayload []model.Position
	for _, position := range payload.Positions {
		positionPayload = append(positionPayload, model.Position{
			ID:   token.SignedToken(3),
			Name: position,
		})
	}

	if err := h.db.CreatePositions(positionPayload); err != nil {
		return err
	}

	return nil
}

func (h *handler) ListPositions() ([]model.Position, error) {
	positions, err := h.db.ListPositions()
	if err != nil {
		return nil, err
	}

	if len(positions) == 0 {
		return nil, &errors.LeagueifyError{Message: "no positions"}
	}

	return positions, nil
}
