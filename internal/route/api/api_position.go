package api

import (
	"net/http"

	"github.com/leagueify/leagueify/internal/lib/error"
	"github.com/leagueify/leagueify/internal/lib/response"
	"github.com/leagueify/leagueify/internal/model"
	"github.com/leagueify/leagueify/internal/route"

	"github.com/labstack/echo/v4"
)

func positionEndpoints(e *echo.Group) {
	e.GET("/positions", listPositions)
	e.POST("/positions", route.AdminRequired(createPositions, "api"))
}

func createPositions(ctx echo.Context) error {
	payload := model.PositionCreation{}

	if err := h.BindAndValidatePayload(ctx, &payload); err != nil {
		return response.JSON(
			ctx, http.StatusBadRequest, errors.HandleError(err),
		)
	}

	if err := h.CreatePositions(payload); err != nil {
		return response.JSON(
			ctx, http.StatusBadRequest, errors.HandleError(err),
		)
	}

	return response.JSON(ctx, http.StatusCreated, nil)
}

func listPositions(ctx echo.Context) error {
	positions, err := h.ListPositions()
	if err != nil {
		return response.JSON(
			ctx, http.StatusBadRequest, errors.HandleError(err),
		)
	}

	return response.JSON(ctx, http.StatusOK, positions)
}
