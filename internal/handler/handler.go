package handler

import (
	"github.com/leagueify/leagueify/internal/database"
	"github.com/leagueify/leagueify/internal/lib/error"

	"github.com/getsentry/sentry-go"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type handler struct {
	db        database.Database
	Validator *customValidator
}

type customValidator struct {
	validator *validator.Validate
}

func (cv *customValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return err
	}
	return nil
}

func GetHandler() *handler {
	db, err := database.GetDatabase()
	if err != nil {
		sentry.CaptureException(err)
	}
	return &handler{
		db: db, Validator: &customValidator{validator: validator.New()},
	}
}

func (h *handler) BindAndValidatePayload(ctx echo.Context, payload interface{}) error {
	if err := ctx.Bind(&payload); err != nil {
		return &errors.LeagueifyError{Message: "invalid json payload"}
	}

	if err := ctx.Validate(payload); err != nil {
		return err
	}

	return nil
}
