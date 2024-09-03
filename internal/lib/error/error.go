package errors

import (
	"fmt"
	"net/textproto"
	"reflect"
	"strings"

	"github.com/getsentry/sentry-go"
	"github.com/go-playground/validator/v10"
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
	case *textproto.Error:
		return textprotoErrors(errType)
	case validator.ValidationErrors:
		return validationErrors(errType)
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

func textprotoErrors(err *textproto.Error) string {
	switch err.Code {
	case 535:
		return "credential authentication failure"
	default:
		msg := fmt.Sprintf("unknown textproto.Error: '%v'", err.Code)
		sentry.CaptureMessage(msg)
		return msg
	}
}

func validationErrors(validationErrors validator.ValidationErrors) string {
	var missingFields []string
	for _, err := range validationErrors {
		if err.Tag() == "required" {
			missingFields = append(missingFields, err.Field())
		}
		if err.Tag() == "e164" {
			return "phone must use the E.164 international standard"
		}
		if err.Tag() == "email" {
			return "invalid email"
		}
		if err.Tag() == "min" {
			return fmt.Sprintf(
				"'%s' must have a minimum length of '%v' characters",
				err.Field(), err.Param(),
			)
		}
	}
	if len(missingFields) != 0 {
		return fmt.Sprintf("missing required field(s): %v", missingFields)
	}
	msg := fmt.Sprintf("unknown validation error: %v", validationErrors)
	sentry.CaptureMessage(msg)
	return msg
}
