package errors

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/getsentry/sentry-go"
	"github.com/lib/pq"
)

type LeagueifyError struct {
	Message string
}

func (e *LeagueifyError) Error() string {
	return e.Message
}

func HandleError(err error) string {
	switch errType := err.(type) {
	case *LeagueifyError:
		return err.Error()
	case *pq.Error:
		return postgresErrors(errType)
	default:
		msg := fmt.Sprintf(
			"unknown error type: %v from %v",
			reflect.TypeOf(err), err.Error(),
		)
		sentry.CaptureMessage(msg)
		return msg
	}
}

func postgresErrors(err *pq.Error) string {
	if err.Code.Name() == "unique_violation" {
		key := strings.Split(err.Constraint, "_")[1]
		return fmt.Sprintf("%v already in use", key)
	}
	msg := fmt.Sprintf(
		"unandled postgres error: %v from %v",
		err.Code.Name(), err.Error(),
	)
	sentry.CaptureMessage(msg)
	return err.Code.Name()
}
