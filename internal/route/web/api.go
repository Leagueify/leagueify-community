package web

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func API(e *echo.Echo) {
	e.GET("/api", sendAPIDocs)
}

func sendAPIDocs(c echo.Context) error {
	_ = c.Render(http.StatusOK, "api-docs", nil)
	return nil
}
