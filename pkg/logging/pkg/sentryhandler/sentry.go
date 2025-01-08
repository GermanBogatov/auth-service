package sentryhandler

import (
	"errors"
	"github.com/getsentry/sentry-go"
)

func CaptureMessage(arg string, level sentry.Level) {
	localHub := sentry.CurrentHub().Clone()
	localHub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetLevel(level)
	})
	localHub.CaptureMessage(arg)
}

func CaptureException(arg string, level sentry.Level) {
	localHub := sentry.CurrentHub().Clone()
	localHub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetLevel(level)
	})
	localHub.CaptureException(errors.New(arg))
}
