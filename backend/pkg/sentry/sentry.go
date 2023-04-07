package sentry

import (
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/hiennguyen9874/stockk-go/config"
)

func Init(cfg *config.Config) error {
	err := sentry.Init(sentry.ClientOptions{
		Dsn: cfg.Sentry.Dsn,
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for performance monitoring.
		// We recommend adjusting this value in production,
		TracesSampleRate: 1.0,
		Environment:      cfg.Sentry.Environment,
	})
	if err != nil {
		return err
	}

	return nil
}

func Flush() {
	// Flush buffered events before the program terminates.
	sentry.Flush(2 * time.Second)
}
