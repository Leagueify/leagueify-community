package api

import (
	"net/http"

	"github.com/leagueify/leagueify/internal/lib/error"
	"github.com/leagueify/leagueify/internal/lib/response"
	"github.com/leagueify/leagueify/internal/model"
	"github.com/leagueify/leagueify/internal/route"

	"github.com/labstack/echo/v4"
)

func seasonEndpoints(e *echo.Group) {
	e.GET("/seasons", listSeasons)
	e.POST("/seasons", route.AdminRequired(createSeason, "api"))
	e.GET("/seasons/:id", getSeason)
	e.PATCH("/seasons/:id", route.AdminRequired(updateSeason, "api"))
}

func createSeason(ctx echo.Context) error {
	payload := model.Season{}

	if err := h.BindAndValidatePayload(ctx, &payload); err != nil {
		return response.JSON(
			ctx, http.StatusBadRequest, errors.HandleError(err),
		)
	}

	if err := h.CreateSeason(payload); err != nil {
		return response.JSON(
			ctx, http.StatusBadRequest, errors.HandleError(err),
		)
	}

	return response.JSON(ctx, http.StatusCreated, nil)
}

func getSeason(ctx echo.Context) error {
	season, err := h.GetSeasonByID(ctx.Param("id"))
	if err != nil {
		return response.JSON(
			ctx, http.StatusBadRequest, errors.HandleError(err),
		)
	}

	return response.JSON(ctx, http.StatusOK, season)
}

func listSeasons(ctx echo.Context) error {
	seasons, err := h.ListSeasons()
	if err != nil {
		return response.JSON(
			ctx, http.StatusBadRequest, errors.HandleError(err),
		)
	}

	return response.JSON(ctx, http.StatusOK, seasons)
}

func updateSeason(ctx echo.Context) error {
	payload := model.SeasonUpdate{}

	if err := h.BindAndValidatePayload(ctx, &payload); err != nil {
		return response.JSON(
			ctx, http.StatusBadRequest, errors.HandleError(err),
		)
	}

	if err := h.UpdateSeason(ctx.Param("id"), payload); err != nil {
		return response.JSON(
			ctx, http.StatusBadRequest, errors.HandleError(err),
		)
	}

	return response.JSON(ctx, http.StatusOK, nil)
}
