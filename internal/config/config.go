package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/leagueify/leagueify/internal/lib/token"
)

type config struct {
	DB        string
	DBConnStr string
	JWTSecret string
	Sentry    bool
	SentryDSN string
	SentryENV string
	SentryTSR float64
}

func LoadConfig() *config {
	config := &config{}
	config.setDefaults()
	config.loadFromEnv()
	return config
}

func (c *config) loadFromEnv() {
	// database config
	// database service
	if db := os.Getenv("DATABASE"); db != "" {
		c.DB = strings.TrimSpace(db)
	}
	// database connection string
	if dbConnStr := os.Getenv("DB_CONN_STR"); dbConnStr != "" {
		c.DBConnStr = strings.TrimSpace(dbConnStr)
	}

	// jwt config
	// jwt secret
	if jwtSecret := os.Getenv("JWT_SECRET"); jwtSecret != "" {
		c.JWTSecret = strings.TrimSpace(jwtSecret)
	}

	// sentry config
	// sentry
	if sentry := os.Getenv("SENTRY"); sentry != "" {
		enabled, err := strconv.ParseBool(strings.TrimSpace(sentry))
		if err != nil {
			panic("Invalid Environment Variable: SENTRY")
		}
		c.Sentry = enabled
	}
	// sentry dsn
	if sentryDSN := os.Getenv("SENTRY_DSN"); sentryDSN != "" {
		c.SentryDSN = strings.TrimSpace(sentryDSN)
	}
	// sentry environment
	if sentryENV := os.Getenv("SENTRY_ENV"); sentryENV != "" {
		c.SentryENV = strings.TrimSpace(sentryENV)
	}
	// sentry tsr
	if sentryTSR := os.Getenv("SENTRY_TSR"); sentryTSR != "" {
		tsr, err := strconv.ParseFloat(strings.TrimSpace(sentryTSR), 64)
		if err != nil {
			panic("Invalid Environment Variable: SENTRY_TSR")
		}
		c.SentryTSR = tsr
	}
}

func (c *config) setDefaults() {
	// database
	c.DB = "postgres"

	// jwt
	c.JWTSecret = token.UnsignedToken(32)

	// Sentry
	c.Sentry = true
	c.SentryDSN = "https://502e62292f21d9aa5841d46babe95fbf@o4507651956932608.ingest.us.sentry.io/4507827637780480"
	c.SentryENV = "production"
	c.SentryTSR = 1.0
}
