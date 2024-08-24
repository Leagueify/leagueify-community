package api

import (
	"net/http"

	"github.com/leagueify/leagueify/internal/lib/error"
	"github.com/leagueify/leagueify/internal/lib/response"
	"github.com/leagueify/leagueify/internal/model"
	"github.com/leagueify/leagueify/internal/route"

	"github.com/labstack/echo/v4"
)

func accountEndpoints(e *echo.Group) {
	e.POST("/accounts", createAccount)
	e.POST("/accounts/login", loginAccount)
	e.POST("/accounts/logout", route.AuthRequired(logoutAccount, "api"))
	e.POST("/accounts/:id/activate", activateAccount)
}

func activateAccount(ctx echo.Context) error {
	if err := h.ActivateAccount(ctx.Param("id")); err != nil {
		return response.JSON(ctx, http.StatusUnauthorized, nil)
	}

	return response.JSON(ctx, http.StatusOK, nil)
}

func createAccount(ctx echo.Context) error {
	payload := model.AccountCreation{}

	if err := h.BindAndValidatePayload(ctx, &payload); err != nil {
		return response.JSON(
			ctx, http.StatusBadRequest, errors.HandleError(err),
		)
	}

	if err := h.CreateAccount(payload); err != nil {
		return response.JSON(
			ctx, http.StatusBadRequest, errors.HandleError(err),
		)
	}

	return response.JSON(ctx, http.StatusCreated, nil)
}

func loginAccount(ctx echo.Context) error {
	payload := model.AccountCredentials{}

	if err := h.BindAndValidatePayload(ctx, &payload); err != nil {
		return response.JSON(
			ctx, http.StatusBadRequest, errors.HandleError(err),
		)
	}

	authToken, err := h.LoginAccount(payload, "api")
	if err != nil {
		return response.JSON(ctx, http.StatusUnauthorized, nil)
	}

	return response.JSON(ctx, http.StatusOK, authToken)
}

func logoutAccount(ctx echo.Context) error {
	return response.JSON(ctx, http.StatusOK, nil)
}
