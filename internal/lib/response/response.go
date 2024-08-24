package response

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func JSON(ctx echo.Context, status int, detail interface{}) error {
	responseData := map[string]interface{}{
		"status": strings.ToLower(http.StatusText(status)),
	}

	if detail != nil {
		responseData["detail"] = detail
	}

	return ctx.JSON(status, responseData)
}
