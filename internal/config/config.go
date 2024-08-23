package config

import (
	"os"
	"strconv"
	"strings"
)

type config struct {
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
	// Sentry
	c.Sentry = true
	c.SentryDSN = "https://502e62292f21d9aa5841d46babe95fbf@o4507651956932608.ingest.us.sentry.io/4507827637780480"
	c.SentryENV = "production"
	c.SentryTSR = 1.0
}
