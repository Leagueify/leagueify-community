package api

import (
	"net/http"

	"github.com/leagueify/leagueify/internal/lib/error"
	"github.com/leagueify/leagueify/internal/lib/response"
	"github.com/leagueify/leagueify/internal/model"
	"github.com/leagueify/leagueify/internal/route"

	"github.com/labstack/echo/v4"
)

func leagueEndpoints(e *echo.Group) {
	e.POST("/leagues", route.AuthRequired(createLeague, "api"))
}

func createLeague(ctx echo.Context) error {
	payload := model.LeagueCreation{}

	if err := h.BindAndValidatePayload(ctx, &payload); err != nil {
		return response.JSON(
			ctx, http.StatusBadRequest, errors.HandleError(err),
		)
	}

	payload.MasterAdmin = ctx.Get("user").(string)
	if err := h.CreateLeague(payload); err != nil {
		return response.JSON(
			ctx, http.StatusBadRequest, errors.HandleError(err),
		)
	}

	return response.JSON(ctx, http.StatusCreated, nil)
}
