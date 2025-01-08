package sentry

import (
	"github.com/getsentry/sentry-go"
	"time"
)

func InitSentry(systemName string, dsn string, env string, debug bool) error {
	sentrySyncTransport := sentry.NewHTTPSyncTransport()
	sentrySyncTransport.Timeout = time.Second * 2

	clientOptions := sentry.ClientOptions{
		ServerName:       systemName,
		Dsn:              dsn,
		Debug:            debug,
		Environment:      env,
		Transport:        sentrySyncTransport,
		AttachStacktrace: true,
		BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			if event.Level == sentry.LevelWarning || event.Level == sentry.LevelError || event.Level == sentry.LevelFatal {
				return event
			}
			return nil
		},
	}

	err := sentry.Init(clientOptions)
	if err != nil {
		return err
	}

	return nil
}
