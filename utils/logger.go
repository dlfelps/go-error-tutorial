package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

// NewLogger creates a new configured logrus logger
func NewLogger() *logrus.Logger {
	logger := logrus.New()

	// Set the default log level to Info
	logger.SetLevel(logrus.InfoLevel)

	// Configure formatter with caller information
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			// Extract just the filename and line number
			fileName := filepath.Base(f.File)
			return "", fmt.Sprintf("%s:%d", fileName, f.Line)
		},
	})

	// Enable caller information
	logger.SetReportCaller(true)

	// Set output to stderr by default
	logger.SetOutput(os.Stderr)

	return logger
}

// FileLogger creates a logger that writes to a file
func FileLogger(filePath string) (*logrus.Logger, error) {
	// Create the logger
	logger := NewLogger()

	// Open the log file
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("could not open log file: %v", err)
	}

	// Set output to the file
	logger.SetOutput(file)

	return logger, nil
}

// LogLevelFromString converts a string to a logrus log level
func LogLevelFromString(level string) (logrus.Level, error) {
	switch strings.ToLower(level) {
	case "debug":
		return logrus.DebugLevel, nil
	case "info":
		return logrus.InfoLevel, nil
	case "warn", "warning":
		return logrus.WarnLevel, nil
	case "error":
		return logrus.ErrorLevel, nil
	case "fatal":
		return logrus.FatalLevel, nil
	case "panic":
		return logrus.PanicLevel, nil
	default:
		return logrus.InfoLevel, fmt.Errorf("invalid log level: %s", level)
	}
}

// ContextLogger adds common fields to a logger for consistent context
func ContextLogger(baseLogger *logrus.Logger, fields logrus.Fields) *logrus.Entry {
	return baseLogger.WithFields(fields)
}

// ErrorWithContext creates an error log entry with contextual information
func ErrorWithContext(logger *logrus.Logger, err error, message string, fields logrus.Fields) {
	if fields == nil {
		fields = logrus.Fields{}
	}

	// Extract stack trace if available from github.com/pkg/errors
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}

	var stackTrace string
	if err, ok := err.(stackTracer); ok {
		stackTrace = fmt.Sprintf("%+v", err.StackTrace())
		fields["stack_trace"] = stackTrace
	}

	logger.WithFields(fields).WithError(err).Error(message)
}
