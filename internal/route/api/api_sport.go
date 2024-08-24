package api

import (
	"net/http"

	"github.com/leagueify/leagueify/internal/lib/response"
	"github.com/leagueify/leagueify/internal/route"

	"github.com/labstack/echo/v4"
)

func sportEndpoints(e *echo.Group) {
	e.GET("/sports", route.AuthRequired(listSports, "api"))
}

func listSports(ctx echo.Context) error {
	sports, err := h.ListSports()
	if err != nil {
		return response.JSON(ctx, http.StatusNotFound, nil)
	}

	return response.JSON(ctx, http.StatusOK, sports)
}
