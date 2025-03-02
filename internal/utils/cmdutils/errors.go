package cmdutils

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/smithy-go"
	"github.com/charmbracelet/huh"
)

// some of this is copied from GithubCLI

// Error unrelated to usage
// FlagErrorf returns a new FlagError that wraps an error produced by
// fmt.Errorf(format, args...).
func FlagErrorf(format string, args ...interface{}) error {
	return FlagErrorWrap(fmt.Errorf(format, args...))
}

// FlagErrorWrap returns a new FlagError that wraps the specified error.
func FlagErrorWrap(err error) error { return &FlagError{err} }

// A *FlagError indicates an error processing command-line flags or other arguments.
// Such errors cause the application to display the usage message.
type FlagError struct {
	// Note: not struct{error}: only *FlagError should satisfy error.
	err error
}

func (fe *FlagError) Error() string {
	return fe.err.Error()
}

func (fe *FlagError) Unwrap() error {
	return fe.err
}

// SilentError is an error that triggers exit code 1 without any error messaging
var SilentError = errors.New("SilentError")

// CancelError signals user-initiated cancellation
var CancelError = errors.New("CancelError")

// PendingError signals nothing failed but something is pending
var PendingError = errors.New("PendingError")

func IsUserCancellation(err error) bool {
	return errors.Is(err, CancelError) || errors.Is(err, context.Canceled) || errors.Is(err, huh.ErrUserAborted)
}

type NonUsageError struct {
	Message string
}

func (e *NonUsageError) Error() string {
	return e.Message
}

var _ error = (*NonUsageError)(nil)

func NewNonUsageError(text string) *NonUsageError {
	return &NonUsageError{
		Message: text,
	}
}

func NewNonUsageErrorf(fmtstr string, v ...any) *NonUsageError {
	return &NonUsageError{
		Message: fmt.Sprintf(fmtstr, v...),
	}
}

func IsNonUsageError(err error) bool {
	var nue *NonUsageError
	return errors.As(err, &nue)
}

type awsError interface {
	ErrorFault() smithy.ErrorFault
}

func IsAWSError(err error) bool {
	_, ok := err.(awsError)
	return ok
}
