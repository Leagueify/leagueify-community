package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// custom banner
	fmt.Print(getBanner())
	// echo initialization and middleware config
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} -- ${remote_ip} -- ${status}:${method}:${uri}\n",
	}))
	e.Use(middleware.Recover())
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!\n")
	})
	// start server
	e.Logger.Fatal(e.Start(":8888"))
}

func getBanner() string {
	// generated using
	// https://patorjk.com/software/taag/#p=display&f=Kban&t=LEAGUEIFY
	return `
'||'      '||''''|      |      ..|'''.|  '||'  '|' '||''''|  '||' '||''''| '||' '|'
 ||        ||  .       |||    .|'     '   ||    |   ||  .     ||   ||  .     || |
 ||        ||''|      |  ||   ||    ....  ||    |   ||''|     ||   ||''|      ||
 ||        ||        .''''|.  '|.    ||   ||    |   ||        ||   ||         ||
.||.....| .||.....| .|.  .||.  ''|...'|    '|..'   .||.....| .||. .||.       .||.
`
}
