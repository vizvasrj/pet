package myerror

import (
	"fmt"
	"runtime"
	"strings"

	"go.uber.org/zap"
	"google.golang.org/grpc/status"
)

// MyError represents a structured error with the source line information.
type MyError struct {
	Inner      error  // The original, underlying error
	Message    string // User-friendly message
	SourceLine string // File and line where the error originated
}

// Error returns a string representation of the error.
func (e MyError) Error() string {
	var sb strings.Builder
	sb.WriteString(e.Message)
	if e.Inner != nil {
		sb.WriteString(", Inner Error: ")
		sb.WriteString(e.Inner.Error())
	}
	if e.SourceLine != "" {
		sb.WriteString(fmt.Sprintf(" (Source: %s)", e.SourceLine))
	}
	return sb.String()
}

// WrapError creates a new MyError, logs it, and returns the error.
func WrapError(logger *zap.Logger, err error, messagef string, msgArgs ...any) error {
	if _, ok := status.FromError(err); ok {
		return err
	}
	_, file, line, _ := runtime.Caller(1) // Get the caller's file and line
	sourceLine := fmt.Sprintf("%s:%d", file, line)

	entry := logger.With(zap.String("source_line", sourceLine), zap.Error(err))

	newMessage := fmt.Sprintf(messagef, msgArgs...)
	entry.Error(newMessage)
	return MyError{
		Inner:      err,
		Message:    newMessage,
		SourceLine: sourceLine, // Use the new source line
	}

}

// New creates a new MyError, logs it, and returns the error.
func New(logger *zap.Logger, message string) MyError {
	_, file, line, _ := runtime.Caller(1)
	sourceLine := fmt.Sprintf("%s:%d", file, line)

	entry := logger.With(zap.String("source_line", sourceLine))
	entry.Error(message)

	return MyError{
		Message:    message,
		SourceLine: sourceLine,
	}
}

// Is checks if target error is the same as the wrapped error or its inner error.
func Is(err, target error) bool {
	if err == target {
		return true
	}
	if myErr, ok := err.(MyError); ok && myErr.Inner != nil {
		return myErr.Inner == target
	}
	return false
}
