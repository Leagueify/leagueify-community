package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/leagueify/leagueify/internal/lib/date"
	"github.com/leagueify/leagueify/internal/lib/error"
	"github.com/leagueify/leagueify/internal/lib/token"
	"github.com/leagueify/leagueify/internal/model"
)

func (h *handler) CreatePlayers(accountID string, payload model.Players) error {
	if len(payload.Players) <= 0 {
		return &errors.LeagueifyError{Message: "no players in payload"}
	}

	for i, player := range payload.Players {
		if !date.ValidDate(player.DateOfBirth) {
			return &errors.LeagueifyError{
				Message: "invalid dateOfBirth",
			}
		}

		payload.Players[i].ID = token.SignedToken(10)
		payload.Players[i].ParentID = accountID
		payload.Players[i].Hash = generatePlayerHash(
			player.FirstName + player.LastName + player.DateOfBirth,
		)

		if h.db.IsExistingPlayer(payload.Players[i].Hash) {
			return &errors.LeagueifyError{
				Message: fmt.Sprintf(
					"player '%s %s' already exists",
					player.FirstName, player.LastName,
				),
			}
		}

		// validate player payload
		if err := h.Validator.Validate(&player); err != nil {
			return err
		}
	}

	if err := h.db.CreatePlayers(payload); err != nil {
		return err
	}

	return nil
}

func (h *handler) DeletePlayer(accountID, playerID string) error {
	if !token.VerifyToken(playerID) {
		return &errors.LeagueifyError{Message: "invalid playerID"}
	}

	if err := h.db.DeletePlayer(accountID, playerID); err != nil {
		return err
	}

	return nil
}

func (h *handler) GetPlayer(accountID, playerID string) (model.Player, error) {
	if !token.VerifyToken(playerID) {
		return model.Player{}, &errors.LeagueifyError{
			Message: "invalid playerID",
		}
	}

	player, err := h.db.GetPlayerByID(accountID, playerID)
	if err != nil {
		return model.Player{}, err
	}

	return player, nil
}

func (h *handler) ListPlayers(accountID string) ([]model.PlayerList, error) {
	players, err := h.db.ListPlayers(accountID)
	if err != nil {
		return nil, err
	}

	if len(players) == 0 {
		return nil, &errors.LeagueifyError{Message: "no players"}
	}

	return players, nil
}

func generatePlayerHash(inputString string) string {
	hash := sha256.Sum256([]byte(strings.ToLower(inputString)))
	return hex.EncodeToString(hash[:])
}
