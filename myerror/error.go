package myerror

import (
	"fmt"
	"runtime"
	"strings"

	pg "github.com/lib/pq"
)

type MyError struct {
	Inner            error
	Message          string
	SingleStacktrace string
}

func (e MyError) Status() int {
	// Default to 500 if no specific status code is set
	return 500
}

func WrapError(err error, messagef string, msgArgs ...any) MyError {
	_, currentFile, currentLine, _ := runtime.Caller(1)
	stackTrace := fmt.Sprintf(" >>: %s:%d ", currentFile, currentLine)

	switch err2 := err.(type) {
	case MyError:
		// If it's already a MyError, append the new message and stack trace
		return MyError{
			Inner:            err2.Inner, // Keep the original inner error
			Message:          fmt.Sprintf("%s >> %s", err2.Message, fmt.Sprintf(messagef, msgArgs...)),
			SingleStacktrace: fmt.Sprintf("%s%s", err2.SingleStacktrace, stackTrace),
		}

	case *pg.Error:
		// Handle PostgreSQL errors, extracting relevant information
		pgerr := err.(*pg.Error)
		return MyError{
			Inner:            err2,
			Message:          pgerr.Message, // Use the PostgreSQL error message
			SingleStacktrace: stackTrace,
		}

	default:
		// For other error types, create a new MyError
		return MyError{
			Inner:            err,
			Message:          fmt.Sprintf(messagef, msgArgs...),
			SingleStacktrace: stackTrace,
		}
	}
}

func (err MyError) Error() string {
	// Build a comprehensive error string, including wrapped messages and stack traces
	var sb strings.Builder
	sb.WriteString(err.Message)
	if err.Inner != nil {
		sb.WriteString(", Inner Error: ")
		sb.WriteString(err.Inner.Error())
	}
	if err.SingleStacktrace != "" {
		sb.WriteString(" Stack Trace:")
		sb.WriteString(err.SingleStacktrace)
	}
	return sb.String()
}

// New allows creating custom errors with specific status codes, but it's simplified
func New(text string, any ...interface{}) error {
	return fmt.Errorf(text, any...)
}

// Is checks if the underlying error matches
func Is(incomingError, matchError error) bool {
	ierr, ok := incomingError.(MyError)
	return ok && ierr.Inner == matchError
}
