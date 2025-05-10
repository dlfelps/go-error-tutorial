package errors

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// Re-export functions from github.com/pkg/errors
var (
	Wrap   = errors.Wrap
	Wrapf  = errors.Wrapf
	Cause  = errors.Cause
	Is     = errors.Is
	As     = errors.As
	New    = errors.New
	Errorf = fmt.Errorf
)

// ValidationError represents a validation error for a specific field
type ValidationError struct {
	Field   string
	Message string
}

// Error implements the error interface
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error for field '%s': %s", e.Field, e.Message)
}

// NetworkError represents an error occurring during network operations
type NetworkError struct {
	URL       string
	Op        string
	Cause     error
	Retriable bool
}

// Error implements the error interface
func (e *NetworkError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s operation failed for URL %s: %v", e.Op, e.URL, e.Cause)
	}
	return fmt.Sprintf("%s operation failed for URL %s", e.Op, e.URL)
}

// Unwrap returns the underlying cause of the error
func (e *NetworkError) Unwrap() error {
	return e.Cause
}

// IsRetriable returns whether the error is retriable
func (e *NetworkError) IsRetriable() bool {
	return e.Retriable
}

// NewNetworkError creates a new NetworkError
func NewNetworkError(url, op string, cause error, retriable bool) *NetworkError {
	return &NetworkError{
		URL:       url,
		Op:        op,
		Cause:     cause,
		Retriable: retriable,
	}
}

// DatabaseError represents an error occurring during database operations
type DatabaseError struct {
	Operation string
	Table     string
	Cause     error
}

// Error implements the error interface
func (e *DatabaseError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("database operation '%s' on table '%s' failed: %v", e.Operation, e.Table, e.Cause)
	}
	return fmt.Sprintf("database operation '%s' on table '%s' failed", e.Operation, e.Table)
}

// Unwrap returns the underlying cause of the error
func (e *DatabaseError) Unwrap() error {
	return e.Cause
}

// NewDatabaseError creates a new DatabaseError
func NewDatabaseError(operation, table string, cause error) *DatabaseError {
	return &DatabaseError{
		Operation: operation,
		Table:     table,
		Cause:     cause,
	}
}

// MultiError is an error type that combines multiple errors
type MultiError struct {
	Errors []error
}

// Error implements the error interface
func (e *MultiError) Error() string {
	if len(e.Errors) == 0 {
		return "no errors"
	}

	errorMessages := make([]string, len(e.Errors))
	for i, err := range e.Errors {
		errorMessages[i] = err.Error()
	}

	return fmt.Sprintf("multiple errors occurred: [%s]", strings.Join(errorMessages, "; "))
}

// Add adds an error to the MultiError
func (e *MultiError) Add(err error) {
	if err != nil {
		e.Errors = append(e.Errors, err)
	}
}

// HasErrors returns true if the MultiError contains any errors
func (e *MultiError) HasErrors() bool {
	return len(e.Errors) > 0
}

// NewMultiError creates a new MultiError
func NewMultiError() *MultiError {
	return &MultiError{Errors: []error{}}
}
