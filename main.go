package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"

	"./config"
	"./dbops"
	"./errors"
	"./fileops"
	"./netops"
	"./utils"
)

func main() {
	// Initialize logger
	log := utils.NewLogger()
	log.Info("Starting error handling demonstration application")

	// Load configuration with error handling
	cfg, err := config.Load("config.json")
	if err != nil {
		log.WithError(err).Fatal("Failed to load configuration")
	}

	// Create a cancellable context that will be used across operations
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Create a wait group to ensure all goroutines finish before exiting
	var wg sync.WaitGroup

	// Start a goroutine to handle cancellation on signal
	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case sig := <-sigChan:
			log.WithField("signal", sig.String()).Info("Received termination signal")
			cancel()
		case <-ctx.Done():
			return
		}
	}()

	// Demonstrate basic error handling
	demoBasicErrorHandling(log)

	// Demonstrate custom error types
	demoCustomErrorTypes(log)

	// Demonstrate error wrapping
	demoErrorWrapping(log)

	// Demonstrate panic and recover
	demoPanicAndRecover(log)

	// Demonstrate file operations with proper error handling
	wg.Add(1)
	go func() {
		defer wg.Done()
		demoFileOperations(ctx, log)
	}()

	// Demonstrate network operations with error handling
	wg.Add(1)
	go func() {
		defer wg.Done()
		demoNetworkOperations(ctx, log)
	}()

	// Demonstrate database operations with error handling
	wg.Add(1)
	go func() {
		defer wg.Done()
		demoDatabaseOperations(ctx, log, cfg.DatabasePath)
	}()

	// Wait for all operations to complete
	wg.Wait()
	log.Info("Application completed successfully")
}

// demoBasicErrorHandling demonstrates the most basic form of error handling in Go
func demoBasicErrorHandling(log *logrus.Logger) {
	log.Info("Demonstrating basic error handling")

	// Basic if err != nil pattern
	value, err := divide(10, 2)
	if err != nil {
		log.WithError(err).Error("Division failed")
	} else {
		log.WithField("result", value).Info("Division successful")
	}

	// Demonstrating the error case
	value, err = divide(10, 0)
	if err != nil {
		log.WithError(err).Error("Division failed")
	} else {
		log.WithField("result", value).Info("Division successful")
	}
}

// divide demonstrates the basic error return pattern
func divide(a, b int) (int, error) {
	if b == 0 {
		return 0, fmt.Errorf("cannot divide by zero")
	}
	return a / b, nil
}

// demoCustomErrorTypes demonstrates the use of custom error types
func demoCustomErrorTypes(log *logrus.Logger) {
	log.Info("Demonstrating custom error types")

	// Using a validation function that returns custom errors
	user := struct {
		Username string
		Email    string
		Age      int
	}{
		Username: "johndoe",
		Email:    "invalid-email",
		Age:      17,
	}

	err := validateUser(user)
	if err != nil {
		switch e := err.(type) {
		case *errors.ValidationError:
			log.WithFields(logrus.Fields{
				"field":   e.Field,
				"message": e.Message,
			}).Error("Validation error")
		default:
			log.WithError(err).Error("Unknown error occurred during validation")
		}
	}
}

// validateUser demonstrates returning custom error types
func validateUser(user struct {
	Username string
	Email    string
	Age      int
}) error {
	if user.Username == "" {
		return &errors.ValidationError{
			Field:   "username",
			Message: "username cannot be empty",
		}
	}

	if !isValidEmail(user.Email) {
		return &errors.ValidationError{
			Field:   "email",
			Message: "email is not valid",
		}
	}

	if user.Age < 18 {
		return &errors.ValidationError{
			Field:   "age",
			Message: "user must be at least 18 years old",
		}
	}

	return nil
}

// isValidEmail is a simple email validation function
func isValidEmail(email string) bool {
	// This is a very simplified check - in real code use a proper validation
	return len(email) > 5 && (email[len(email)-4:] == ".com" || email[len(email)-4:] == ".org")
}

// demoErrorWrapping demonstrates error wrapping using github.com/pkg/errors
func demoErrorWrapping(log *logrus.Logger) {
	log.Info("Demonstrating error wrapping")

	// Simulate a chain of function calls with error wrapping
	err := simulateDeepOperation()
	if err != nil {
		// Log the full error chain
		log.WithError(err).Error("Operation failed")

		// Extract the original error
		original := errors.Cause(err)
		log.WithField("original_error", original.Error()).Info("Original error")
	}
}

// simulateDeepOperation demonstrates a chain of functions that wrap errors
func simulateDeepOperation() error {
	err := levelThree()
	if err != nil {
		return errors.Wrap(err, "deep operation failed")
	}
	return nil
}

func levelThree() error {
	err := levelTwo()
	if err != nil {
		return errors.Wrap(err, "level three failed")
	}
	return nil
}

func levelTwo() error {
	err := levelOne()
	if err != nil {
		return errors.Wrap(err, "level two failed")
	}
	return nil
}

func levelOne() error {
	// Simulate a low-level error
	return fmt.Errorf("simulated low-level error")
}

// demoPanicAndRecover demonstrates panic handling with recover
func demoPanicAndRecover(log *logrus.Logger) {
	log.Info("Demonstrating panic and recover")

	// Use deferred recover function to handle panics
	defer func() {
		if r := recover(); r != nil {
			log.WithField("panic", r).Error("Recovered from panic")
		}
	}()

	// Demonstrate explicit panic
	log.Info("About to call function that will panic...")
	panicFunction()
	log.Info("This line will not be reached")
}

// panicFunction demonstrates a function that panics
func panicFunction() {
	panic("intentional panic for demonstration")
}

// demoFileOperations demonstrates file operations with proper error handling
func demoFileOperations(ctx context.Context, log *logrus.Logger) {
	log.Info("Demonstrating file operations with error handling")

	// Create a test file
	fileName := "test_file.txt"
	err := fileops.WriteFile(fileName, "Hello, World!")
	if err != nil {
		log.WithError(err).Error("Failed to write to file")
		return
	}
	log.Info("Successfully wrote to file")

	// Read the file with context
	content, err := fileops.ReadFileWithContext(ctx, fileName)
	if err != nil {
		log.WithError(err).Error("Failed to read file")
	} else {
		log.WithField("content", content).Info("Successfully read file")
	}

	// Clean up the file with defer
	defer func() {
		err := os.Remove(fileName)
		if err != nil {
			log.WithError(err).Error("Failed to remove test file")
		} else {
			log.Info("Successfully removed test file")
		}
	}()

	// Try to read a non-existent file to demonstrate error handling
	_, err = fileops.ReadFileWithContext(ctx, "non_existent_file.txt")
	if err != nil {
		log.WithError(err).Error("Failed to read non-existent file (expected)")
	}
}

// demoNetworkOperations demonstrates network operations with error handling
func demoNetworkOperations(ctx context.Context, log *logrus.Logger) {
	log.Info("Demonstrating network operations with error handling")

	// Set up a timeout context
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Fetch data with timeout and retry
	url := "https://jsonplaceholder.typicode.com/posts/1"
	data, err := netops.FetchWithRetry(timeoutCtx, url, 3)
	if err != nil {
		log.WithError(err).Error("Failed to fetch data after retries")
	} else {
		log.WithField("data_length", len(data)).Info("Successfully fetched data")
	}

	// Try with an invalid URL to demonstrate error handling
	invalidURL := "https://non-existent-domain-12345.com"
	_, err = netops.FetchWithRetry(timeoutCtx, invalidURL, 2)
	if err != nil {
		log.WithError(err).Error("Failed to fetch from invalid URL (expected)")
	}
}

// demoDatabaseOperations demonstrates database operations with error handling
func demoDatabaseOperations(ctx context.Context, log *logrus.Logger, dbPath string) {
	log.Info("Demonstrating database operations with error handling")

	// Initialize database with defer for cleanup
	db, err := dbops.InitDatabase(dbPath)
	if err != nil {
		log.WithError(err).Error("Failed to initialize database")
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.WithError(err).Error("Failed to close database connection")
		} else {
			log.Info("Database connection closed successfully")
		}
	}()

	// Create a user
	user := struct {
		ID       int
		Username string
		Email    string
	}{
		ID:       1,
		Username: "johndoe",
		Email:    "john@example.com",
	}

	// Insert with transaction handling
	err = dbops.InsertUser(ctx, db, user)
	if err != nil {
		log.WithError(err).Error("Failed to insert user")
	} else {
		log.WithField("user_id", user.ID).Info("User inserted successfully")
	}

	// Query the user
	retrievedUser, err := dbops.GetUser(ctx, db, user.ID)
	if err != nil {
		log.WithError(err).Error("Failed to retrieve user")
	} else {
		log.WithFields(logrus.Fields{
			"user_id":   retrievedUser.ID,
			"username":  retrievedUser.Username,
			"email":     retrievedUser.Email,
		}).Info("User retrieved successfully")
	}

	// Try to get a non-existent user
	_, err = dbops.GetUser(ctx, db, 999)
	if err != nil {
		if errors.Is(err, dbops.ErrUserNotFound) {
			log.Info("User not found error handled correctly")
		} else {
			log.WithError(err).Error("Unexpected error when retrieving non-existent user")
		}
	}
}
