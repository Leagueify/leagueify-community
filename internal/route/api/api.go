package api

import (
	"github.com/leagueify/leagueify/internal/handler"

	"github.com/labstack/echo/v4"
)

var (
	h = handler.GetHandler()
)

func Routes(e *echo.Echo) {
	// initialize validator
	e.Validator = h.Validator
	// api group
	routes := e.Group("/api")
	// api endpoints
	accountEndpoints(routes)
	divisionEndpoints(routes)
	emailEndpoints(routes)
	leagueEndpoints(routes)
	playerEndpoints(routes)
	positionEndpoints(routes)
	seasonEndpoints(routes)
	sportEndpoints(routes)
}
