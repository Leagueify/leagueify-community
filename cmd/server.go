package main

import (
	"fmt"
	"net/http"

	"github.com/leagueify/leagueify/internal/config"

	"github.com/getsentry/sentry-go"
	sentryecho "github.com/getsentry/sentry-go/echo"
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

	// initialize config
	cfg := config.LoadConfig()

	// setnry initialization
	if cfg.Sentry {
		if err := sentry.Init(sentry.ClientOptions{
			Dsn:              cfg.SentryDSN,
			Environment:      cfg.SentryENV,
			TracesSampleRate: cfg.SentryTSR,
		}); err != nil {
			fmt.Printf("Sentry config failed: %v\n", err)
		}
		e.Use(sentryecho.New(sentryecho.Options{Repanic: true}))
	}

	// initialize routes
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
