package api

import (
	"net/http"

	"github.com/leagueify/leagueify/internal/lib/error"
	"github.com/leagueify/leagueify/internal/lib/response"
	"github.com/leagueify/leagueify/internal/model"
	"github.com/leagueify/leagueify/internal/route"

	"github.com/labstack/echo/v4"
)

func divisionEndpoints(e *echo.Group) {
	e.GET("/divisions", listDivisions)
	e.POST("/divisions", route.AdminRequired(createDivisions, "api"))
}

func createDivisions(ctx echo.Context) error {
	payload := model.DivisionCreation{}

	if err := h.BindAndValidatePayload(ctx, &payload); err != nil {
		return response.JSON(
			ctx, http.StatusBadRequest, errors.HandleError(err),
		)
	}

	if err := h.CreateDivisions(payload); err != nil {
		return response.JSON(
			ctx, http.StatusBadRequest, errors.HandleError(err),
		)
	}

	return response.JSON(ctx, http.StatusCreated, nil)
}

func listDivisions(ctx echo.Context) error {
	divisions, err := h.ListDivisions()
	if err != nil {
		return response.JSON(
			ctx, http.StatusBadRequest, errors.HandleError(err),
		)
	}
	return response.JSON(ctx, http.StatusOK, divisions)
}
