package errors

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// ValidationError is a custom error type for input validation errors
type ValidationError struct {
	Field   string
	Message string
}

// Error implements the error interface
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed for %s: %s", e.Field, e.Message)
}

// ValidationErrors is a collection of validation errors
type ValidationErrors []*ValidationError

// Error implements the error interface for a collection of errors
func (e ValidationErrors) Error() string {
	if len(e) == 0 {
		return "no validation errors"
	}

	var sb strings.Builder
	sb.WriteString("multiple validation errors:\n")
	for _, err := range e {
		sb.WriteString("- " + err.Error() + "\n")
	}
	return sb.String()
}

// UserError represents an error caused by user input
type UserError struct {
	Err error
	Msg string
}

// Error implements the error interface
func (e *UserError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Msg, e.Err)
	}
	return e.Msg
}

// Unwrap returns the underlying error
func (e *UserError) Unwrap() error {
	return e.Err
}

// SystemError represents an internal system error
type SystemError struct {
	Err      error
	Code     string
	Severity string
}

// Error implements the error interface
func (e *SystemError) Error() string {
	return fmt.Sprintf("[%s][%s] %v", e.Severity, e.Code, e.Err)
}

// Unwrap returns the underlying error
func (e *SystemError) Unwrap() error {
	return e.Err
}

// ValidateInput demonstrates using custom error types
func ValidateInput(input string) error {
	if input == "" {
		return &ValidationError{
			Field:   "input",
			Message: "cannot be empty",
		}
	}
	return nil
}

// ProcessWithWrapping demonstrates error wrapping
func ProcessWithWrapping(data string) error {
	// Simulate a chained process with potential errors
	err := step1(data)
	if err != nil {
		// Wrap the error with context
		return errors.Wrap(err, "processing failed at step 1")
	}
	
	err = step2(data)
	if err != nil {
		// Wrap the error with context
		return errors.Wrap(err, "processing failed at step 2")
	}
	
	return nil
}

// step1 is a helper function that might return an error
func step1(data string) error {
	if len(data) < 10 {
		baseErr := fmt.Errorf("data is too short (length: %d)", len(data))
		return &UserError{
			Err: baseErr,
			Msg: "input validation failed",
		}
	}
	return nil
}

// step2 is a helper function that might return an error
func step2(data string) error {
	// This would be a system error, not user's fault
	baseErr := fmt.Errorf("database connection failed")
	return &SystemError{
		Err:      baseErr,
		Code:     "DB_ERROR",
		Severity: "CRITICAL",
	}
}

// PrintErrorChain prints the entire error chain for wrapped errors
func PrintErrorChain(err error, log *logrus.Logger) {
	// Use errors.Cause to get the root cause
	rootCause := errors.Cause(err)
	log.WithError(rootCause).Error("Root cause")

	// Iterate through the error chain
	currentErr := err
	for currentErr != nil {
		log.Error(currentErr.Error())
		// Get the next error in the chain
		unwrapped := errors.Unwrap(currentErr)
		if unwrapped == currentErr || unwrapped == nil {
			break
		}
		currentErr = unwrapped
	}
}

// IsUserError checks if the error is a user error
func IsUserError(err error) bool {
	var ue *UserError
	return errors.As(err, &ue)
}

// IsSystemError checks if the error is a system error
func IsSystemError(err error) bool {
	var se *SystemError
	return errors.As(err, &se)
}

// IsValidationError checks if the error is a validation error
func IsValidationError(err error) bool {
	var ve *ValidationError
	return errors.As(err, &ve)
}
