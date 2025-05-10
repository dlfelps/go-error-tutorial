package logger

import (
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// ErrorWithContext wraps an error with additional context information
type ErrorWithContext struct {
	Err     error
	Context map[string]interface{}
}

// Error implements the error interface
func (e *ErrorWithContext) Error() string {
	if e.Err == nil {
		return "no error"
	}
	return e.Err.Error()
}

// Unwrap returns the underlying error
func (e *ErrorWithContext) Unwrap() error {
	return e.Err
}

// NewLogger creates a new structured logger
func NewLogger() *logrus.Logger {
	log := logrus.New()

	// Set the output to stdout
	log.SetOutput(os.Stdout)

	// Set the log level
	log.SetLevel(logrus.InfoLevel)

	// Set the formatter to JSON
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
		PrettyPrint:     false,
	})

	return log
}

// NewFileLogger creates a logger that writes to a file
func NewFileLogger(logPath string) (*logrus.Logger, error) {
	log := logrus.New()

	// Ensure the log directory exists
	logDir := filepath.Dir(logPath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, errors.Wrap(err, "failed to create log directory")
	}

	// Open the log file
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open log file")
	}

	// Set the output to the file
	log.SetOutput(file)

	// Set the log level
	log.SetLevel(logrus.InfoLevel)

	// Set the formatter
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})

	return log, nil
}

// LogError logs an error with context
func LogError(log *logrus.Logger, err error, message string, fields logrus.Fields) {
	// Create a new entry with fields
	entry := log.WithFields(fields)

	// Add error to fields
	if err != nil {
		entry = entry.WithError(err)

		// Check for wrapped error context
		var errWithContext *ErrorWithContext
		if errors.As(err, &errWithContext) && errWithContext.Context != nil {
			// Add context fields from error
			for k, v := range errWithContext.Context {
				entry = entry.WithField(k, v)
			}
		}
	}

	// Log the error
	entry.Error(message)
}

// WrapErrorWithContext adds context to an error
func WrapErrorWithContext(err error, message string, context map[string]interface{}) error {
	// Wrap the error with the message
	wrappedErr := errors.Wrap(err, message)

	// Return a new ErrorWithContext
	return &ErrorWithContext{
		Err:     wrappedErr,
		Context: context,
	}
}

// LogFatalError logs an error and exits
func LogFatalError(log *logrus.Logger, err error, message string, fields logrus.Fields) {
	// Create a new entry with fields
	entry := log.WithFields(fields)

	// Add error to fields
	if err != nil {
		entry = entry.WithError(err)
	}

	// Log the error and exit
	entry.Fatal(message)
}

// LogWithOperation adds operation context to the log entry
func LogWithOperation(log *logrus.Logger, operation string) *logrus.Entry {
	return log.WithField("operation", operation)
}

// LogWithUserContext adds user context to the log entry
func LogWithUserContext(log *logrus.Logger, userID string) *logrus.Entry {
	return log.WithField("user_id", userID)
}

// LogWithRequestContext adds request context to the log entry
func LogWithRequestContext(log *logrus.Logger, requestID string) *logrus.Entry {
	return log.WithField("request_id", requestID)
}
