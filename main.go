package main

import (
        "context"
        "fmt"
        "os"
        "time"

        contextpkg "error-handling-demo/context"
        "error-handling-demo/db"
        "error-handling-demo/errors"
        "error-handling-demo/file"
        "error-handling-demo/logger"
        "error-handling-demo/network"
        panicpkg "error-handling-demo/panic"

        "github.com/sirupsen/logrus"
)

func main() {
        // Initialize structured logger
        log := logger.NewLogger()
        log.Info("Starting error handling demonstration application")

        // Create a base context with timeout
        ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        defer cancel() // Ensure resources are released when done

        // Run different error handling demos
        runBasicErrorHandling(log)
        runCustomErrorDemo(log)
        runFileOperationsDemo(log)
        runNetworkOperationsDemo(ctx, log)
        runDatabaseOperationsDemo(ctx, log)
        runPanicRecoveryDemo(log)
        runContextCancellationDemo(log)

        log.Info("Error handling demonstration completed")
}

// runBasicErrorHandling demonstrates the most basic error handling in Go
func runBasicErrorHandling(log *logrus.Logger) {
        log.Info("=== Basic Error Handling ===")

        // Example 1: Simple error checking
        if _, err := os.Open("non-existent-file.txt"); err != nil {
                log.WithError(err).Error("Failed to open file")
                // Note: We don't panic or exit, we log and continue - graceful degradation
        }

        // Example 2: Error variables
        _, err := os.ReadFile("another-non-existent-file.txt")
        if err != nil {
                // Check specific error types
                if os.IsNotExist(err) {
                        log.Warn("The file does not exist")
                } else if os.IsPermission(err) {
                        log.Warn("Permission denied")
                } else {
                        log.WithError(err).Error("Unknown error occurred")
                }
        }

        // Example 3: Return error from function
        if err := functionThatReturnsError(); err != nil {
                log.WithError(err).Error("Error from function")
        }

        log.Info("Basic error handling demonstration completed")
}

// functionThatReturnsError is a simple function that returns an error
func functionThatReturnsError() error {
        return fmt.Errorf("this is a sample error")
}

// runCustomErrorDemo demonstrates custom error types and error wrapping
func runCustomErrorDemo(log *logrus.Logger) {
        log.Info("=== Custom Error Types and Error Wrapping ===")

        // Example 1: Using custom error types
        if err := errors.ValidateInput(""); err != nil {
                log.WithError(err).Error("Validation error")

                // Type assertion to check for specific error types
                if valErr, ok := err.(*errors.ValidationError); ok {
                        log.WithFields(logrus.Fields{
                                "field":   valErr.Field,
                                "message": valErr.Message,
                        }).Error("Validation details")
                }
        }

        // Example 2: Error wrapping
        err := errors.ProcessWithWrapping("sample data")
        if err != nil {
                log.WithError(err).Error("Process error with wrapping")
                
                // Unwrap to get the original error
                log.Error("Unwrapped error chain:")
                errors.PrintErrorChain(err, log)
        }

        log.Info("Custom error types demonstration completed")
}

// runFileOperationsDemo demonstrates file operations with proper error handling
func runFileOperationsDemo(log *logrus.Logger) {
        log.Info("=== File Operations with Error Handling ===")

        // Create a temporary file for testing
        tempFile, err := os.CreateTemp("", "error-handling-demo-*.txt")
        if err != nil {
                log.WithError(err).Error("Failed to create temporary file")
                return
        }
        tempFileName := tempFile.Name()
        defer func() {
                // Clean up: close and remove the temporary file
                tempFile.Close()
                os.Remove(tempFileName)
                log.Info("Temporary file cleaned up")
        }()

        // Write to file with error handling
        if err := file.WriteToFile(tempFileName, "Hello, error handling world!"); err != nil {
                log.WithError(err).Error("Failed to write to file")
        } else {
                log.Info("Successfully wrote to file")
        }

        // Read from file with error handling
        content, err := file.ReadFromFile(tempFileName)
        if err != nil {
                log.WithError(err).Error("Failed to read from file")
        } else {
                log.WithField("content", content).Info("Successfully read from file")
        }

        // Try to read a non-existent file
        _, err = file.ReadFromFile("this-file-does-not-exist.txt")
        if err != nil {
                log.WithError(err).Error("Expected error: reading non-existent file")
        }

        // Demonstrate safe file copying
        if err := file.SafeCopyFile(tempFileName, "copy-"+tempFileName); err != nil {
                log.WithError(err).Error("Failed to copy file")
        } else {
                log.Info("Successfully copied file")
                // Clean up the copied file
                os.Remove("copy-" + tempFileName)
        }

        log.Info("File operations demonstration completed")
}

// runNetworkOperationsDemo demonstrates network operations with error handling
func runNetworkOperationsDemo(ctx context.Context, log *logrus.Logger) {
        log.Info("=== Network Operations with Error Handling ===")

        // Make a simple HTTP request with retries and timeout
        response, err := network.FetchURL(ctx, "https://httpbin.org/get", 3)
        if err != nil {
                log.WithError(err).Error("Failed to fetch URL after retries")
        } else {
                log.WithField("status", response.Status).Info("Successfully fetched URL")
        }

        // Try an invalid URL to demonstrate error handling
        _, err = network.FetchURL(ctx, "https://invalid-url-that-doesnt-exist.xyz", 2)
        if err != nil {
                log.WithError(err).Error("Expected error: invalid URL")
        }

        // Demonstrate timeout handling
        timeoutCtx, cancel := context.WithTimeout(ctx, 1*time.Millisecond)
        defer cancel()
        
        _, err = network.FetchURL(timeoutCtx, "https://httpbin.org/delay/3", 1)
        if err != nil {
                log.WithError(err).Error("Expected error: request timeout")
        }

        log.Info("Network operations demonstration completed")
}

// runDatabaseOperationsDemo demonstrates database operations with error handling
func runDatabaseOperationsDemo(ctx context.Context, log *logrus.Logger) {
        log.Info("=== Database Operations with Error Handling ===")

        // Initialize database
        dbConn, err := db.OpenDatabase(ctx, ":memory:")
        if err != nil {
                log.WithError(err).Error("Failed to open database")
                return
        }
        defer dbConn.Close()

        // Create schema
        if err := db.CreateSchema(ctx, dbConn); err != nil {
                log.WithError(err).Error("Failed to create database schema")
                return
        }

        // Insert data
        id, err := db.InsertUser(ctx, dbConn, "John Doe", "john@example.com")
        if err != nil {
                log.WithError(err).Error("Failed to insert user")
        } else {
                log.WithField("user_id", id).Info("Successfully inserted user")
        }

        // Query data
        user, err := db.GetUser(ctx, dbConn, id)
        if err != nil {
                log.WithError(err).Error("Failed to get user")
        } else {
                log.WithFields(logrus.Fields{
                        "id":    user.ID,
                        "name":  user.Name,
                        "email": user.Email,
                }).Info("Successfully retrieved user")
        }

        // Try to get a non-existent user
        _, err = db.GetUser(ctx, dbConn, 999)
        if err != nil {
                log.WithError(err).Error("Expected error: user not found")
        }

        // Demonstrate transaction with error handling
        if err := db.ExecuteTransaction(ctx, dbConn); err != nil {
                log.WithError(err).Error("Transaction failed")
        } else {
                log.Info("Transaction completed successfully")
        }

        log.Info("Database operations demonstration completed")
}

// runPanicRecoveryDemo demonstrates panic and recovery mechanisms
func runPanicRecoveryDemo(log *logrus.Logger) {
        log.Info("=== Panic and Recovery Mechanisms ===")

        // Demonstrate panic recovery
        result, err := panicpkg.ExecuteWithRecover(func() (string, error) {
                // This function will panic
                panicpkg.SomethingThatPanics()
                return "This will never be reached", nil
        })

        if err != nil {
                log.WithError(err).Error("Function panicked but was recovered")
        } else {
                log.WithField("result", result).Info("Function executed successfully")
        }

        // Demonstrate safe array access
        values := []int{1, 2, 3}
        
        // Safe access
        if val, err := panicpkg.GetValueSafely(values, 1); err != nil {
                log.WithError(err).Error("Failed to access array")
        } else {
                log.WithField("value", val).Info("Safely accessed array")
        }
        
        // Out of bounds access (would normally panic)
        if val, err := panicpkg.GetValueSafely(values, 10); err != nil {
                log.WithError(err).Error("Expected error: array index out of bounds")
        } else {
                log.WithField("value", val).Info("Safely accessed array")
        }

        log.Info("Panic and recovery demonstration completed")
}

// runContextCancellationDemo demonstrates context usage for cancellation
func runContextCancellationDemo(log *logrus.Logger) {
        log.Info("=== Context Cancellation ===")

        // Create a cancellable context
        ctx, cancel := context.WithCancel(context.Background())
        
        // Start a long-running operation
        go func() {
                time.Sleep(500 * time.Millisecond)
                log.Info("Cancelling the operation")
                cancel()
        }()

        // Execute operation with context
        err := contextpkg.ExecuteWithContext(ctx)
        if err != nil {
                log.WithError(err).Error("Operation was cancelled, as expected")
        }

        // Demonstrate timeout
        timeoutCtx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
        defer cancel()

        err = contextpkg.SlowOperation(timeoutCtx)
        if err != nil {
                log.WithError(err).Error("Operation timed out, as expected")
        }

        log.Info("Context cancellation demonstration completed")
}
