package api

import (
	"net/http"

	"github.com/leagueify/leagueify/internal/lib/error"
	"github.com/leagueify/leagueify/internal/lib/response"
	"github.com/leagueify/leagueify/internal/model"
	"github.com/leagueify/leagueify/internal/route"

	"github.com/labstack/echo/v4"
)

func playerEndpoints(e *echo.Group) {
	e.GET("/players", route.AuthRequired(listPlayers, "api"))
	e.POST("/players", route.AuthRequired(createPlayers, "api"))
	e.DELETE("/players/:id", route.AuthRequired(deletePlayer, "api"))
	e.GET("/players/:id", route.AuthRequired(getPlayer, "api"))
}

func createPlayers(ctx echo.Context) error {
	payload := model.Players{}

	if err := h.BindAndValidatePayload(ctx, &payload); err != nil {
		return response.JSON(
			ctx, http.StatusBadRequest, errors.HandleError(err),
		)
	}

	if err := h.CreatePlayers(ctx.Get("user").(string), payload); err != nil {
		return response.JSON(
			ctx, http.StatusBadRequest, errors.HandleError(err),
		)
	}

	return response.JSON(ctx, http.StatusCreated, nil)
}

func deletePlayer(ctx echo.Context) error {
	if err := h.DeletePlayer(ctx.Get("user").(string), ctx.Param("id")); err != nil {
		return response.JSON(
			ctx, http.StatusBadRequest, errors.HandleError(err),
		)
	}

	return response.JSON(ctx, http.StatusNoContent, nil)
}

func getPlayer(ctx echo.Context) error {
	player, err := h.GetPlayer(ctx.Get("user").(string), ctx.Param("id"))
	if err != nil {
		return response.JSON(
			ctx, http.StatusBadRequest, errors.HandleError(err),
		)
	}

	return response.JSON(ctx, http.StatusOK, player)
}

func listPlayers(ctx echo.Context) error {
	players, err := h.ListPlayers(ctx.Get("user").(string))
	if err != nil {
		return response.JSON(
			ctx, http.StatusBadRequest, errors.HandleError(err),
		)
	}

	return response.JSON(ctx, http.StatusOK, players)
}
