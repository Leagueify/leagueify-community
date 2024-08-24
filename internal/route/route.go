package route

import (
	"net/http"
	"strings"

	"github.com/leagueify/leagueify/internal/lib/auth"
	"github.com/leagueify/leagueify/internal/lib/error"
	"github.com/leagueify/leagueify/internal/lib/response"

	"github.com/labstack/echo/v4"
)

func AuthRequired(f func(echo.Context) error, aud string) echo.HandlerFunc {
	return func(c echo.Context) error {
		var authToken string

		if err := getAuthToken(c, &authToken); err != nil {
			return response.JSON(c, http.StatusUnauthorized, nil)
		}

		claims, err := auth.VerifyJWT(authToken)
		if err != nil {
			return response.JSON(c, http.StatusUnauthorized, nil)
		}

		audience, err := claims.GetAudience()
		if err != nil {
			return response.JSON(c, http.StatusUnauthorized, nil)
		}

		result := false
		for _, a := range audience {
			if a == aud {
				result = true
				break
			}
		}

		if !result {
			return response.JSON(c, http.StatusUnauthorized, nil)
		}

		c.Set("user", claims.Subject)

		return f(c)
	}
}

func getAuthToken(c echo.Context, authToken *string) error {
	authHeader := c.Request().Header.Get("Authorization")

	if strings.HasPrefix(authHeader, "Bearer ") {
		*authToken = strings.TrimSpace(
			strings.TrimPrefix(authHeader, "Bearer "),
		)
	} else {
		cookie, err := c.Cookie("access")
		if err == nil {
			*authToken = cookie.Value
		}
	}

	if *authToken == "" {
		return &errors.LeagueifyError{
			Message: "no authorization token provided",
		}
	}

	return nil
}
