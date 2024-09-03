package api

import (
	"net/http"

	"github.com/leagueify/leagueify/internal/lib/error"
	"github.com/leagueify/leagueify/internal/lib/response"
	"github.com/leagueify/leagueify/internal/model"
	"github.com/leagueify/leagueify/internal/route"

	"github.com/labstack/echo/v4"
)

func emailEndpoints(e *echo.Group) {
	e.GET("/email", route.AdminRequired(getEmailConfig, "api"))
	e.PATCH("/email", route.AdminRequired(updateEmailConfig, "api"))
	e.POST("/email", route.AdminRequired(createEmailConfig, "api"))
}

func createEmailConfig(ctx echo.Context) error {
	payload := model.EmailConfig{}

	if err := h.BindAndValidatePayload(ctx, &payload); err != nil {
		return response.JSON(
			ctx, http.StatusBadRequest, errors.HandleError(err),
		)
	}

	if err := h.CreateEmailConfig(payload); err != nil {
		return response.JSON(
			ctx, http.StatusBadRequest, errors.HandleError(err),
		)
	}

	return response.JSON(ctx, http.StatusCreated, nil)
}

func getEmailConfig(ctx echo.Context) error {
	emailConfig, err := h.GetEmailConfig()

	if err != nil {
		return response.JSON(
			ctx, http.StatusBadRequest, errors.HandleError(err),
		)
	}

	return response.JSON(ctx, http.StatusOK, emailConfig)
}

func updateEmailConfig(ctx echo.Context) error {
	payload := model.UpdateEmailConfig{}

	if err := h.BindAndValidatePayload(ctx, &payload); err != nil {
		return response.JSON(
			ctx, http.StatusBadRequest, errors.HandleError(err),
		)
	}

	emailConfig, err := h.GetEmailConfig()

	if err != nil {
		return response.JSON(
			ctx, http.StatusBadRequest, errors.HandleError(err),
		)
	}

	if err := h.UpdateEmailConfig(emailConfig.ID, payload); err != nil {
		return response.JSON(
			ctx, http.StatusBadRequest, errors.HandleError(err),
		)
	}

	return response.JSON(ctx, http.StatusOK, nil)
}
