package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"
)

// Step 1: Basic Error Handling
// ===========================

func basicErrorHandling() {
	fmt.Println("\n=== BASIC ERROR HANDLING ===")
	fmt.Println("In Go, functions typically return an error as the last return value.")
	fmt.Println("The caller checks if the error is nil using 'if err != nil' pattern.")

	// Example: File operations
	fmt.Println("\nExample: Opening a file")
	fmt.Println("Code: file, err := os.Open(\"non-existent-file.txt\")")

	file, err := os.Open("non-existent-file.txt")
	if err != nil {
		fmt.Printf("Error occurred: %v\n", err)
		fmt.Println("This demonstrates the basic error handling pattern!")
	} else {
		defer file.Close()
		fmt.Println("File opened successfully!")
	}

	// Example: String conversion
	fmt.Println("\nExample: String to int conversion")
	fmt.Println("Code: value, err := strconv.Atoi(\"not-a-number\")")

	str := "not-a-number"
	_, err = fmt.Sscanf(str, "%d", new(int))
	if err != nil {
		fmt.Printf("Error converting '%s' to integer: %v\n", str, err)
	}

	// Common pattern explanation
	fmt.Println("\nCommon error handling pattern in Go:")
	fmt.Println("  result, err := someFunction()")
	fmt.Println("  if err != nil {")
	fmt.Println("      // Handle error")
	fmt.Println("      return err  // Or take other appropriate action")
	fmt.Println("  }")
	fmt.Println("  // Use result if no error occurred")

	waitForEnter()
}

// Step 2: Creating Errors
// =====================

func createErrors() {
	fmt.Println("\n=== CREATING ERRORS ===")
	fmt.Println("Go provides several ways to create errors:")

	// Using errors.New
	fmt.Println("\n1. Using errors.New:")
	fmt.Println("   Code: err := errors.New(\"something went wrong\")")
	err1 := errors.New("something went wrong")
	fmt.Printf("   Result: %v\n", err1)

	// Using fmt.Errorf
	fmt.Println("\n2. Using fmt.Errorf (allows formatting):")
	fmt.Println("   Code: err := fmt.Errorf(\"error with value: %d\", 42)")
	err2 := fmt.Errorf("error with value: %d", 42)
	fmt.Printf("   Result: %v\n", err2)

	// In practice
	fmt.Println("\nIn practice with a divide function:")
	fmt.Println("func divide(a, b int) (int, error) {")
	fmt.Println("    if b == 0 {")
	fmt.Println("        return 0, errors.New(\"division by zero\")")
	fmt.Println("    }")
	fmt.Println("    return a / b, nil")
	fmt.Println("}")

	waitForEnter()
}

// Step 3: Custom Error Types
// ========================

// ValidationError is a custom error type
type ValidationError struct {
	Field   string
	Message string
}

// Error implements the error interface
func (e ValidationError) Error() string {
	return fmt.Sprintf("validation failed for field '%s': %s", e.Field, e.Message)
}

func validateAge(age int) error {
	if age < 0 {
		return ValidationError{Field: "age", Message: "cannot be negative"}
	}
	if age > 120 {
		return ValidationError{Field: "age", Message: "too high, nobody is that old"}
	}
	return nil
}

func customErrorTypes() {
	fmt.Println("\n=== CUSTOM ERROR TYPES ===")
	fmt.Println("You can create custom error types by implementing the error interface:")
	fmt.Println("type error interface {")
	fmt.Println("    Error() string")
	fmt.Println("}")

	fmt.Println("\nExample custom error type for validation:")
	fmt.Println("type ValidationError struct {")
	fmt.Println("    Field   string")
	fmt.Println("    Message string")
	fmt.Println("}")
	fmt.Println("func (e ValidationError) Error() string {")
	fmt.Println("    return fmt.Sprintf(\"validation failed for field '%s': %s\", e.Field, e.Message)")
	fmt.Println("}")

	fmt.Println("\nUsing the custom error type:")
	for _, age := range []int{-5, 150, 30} {
		err := validateAge(age)
		if err != nil {
			fmt.Printf("For age %d: %v\n", age, err)

			// Type assertion to access the structure fields
			if valErr, ok := err.(ValidationError); ok {
				fmt.Printf("  Field: %s, Message: %s\n", valErr.Field, valErr.Message)
			}
		} else {
			fmt.Printf("Age %d is valid\n", age)
		}
	}

	fmt.Println("\nAdvantages of custom error types:")
	fmt.Println("1. Can include structured data (like which field failed)")
	fmt.Println("2. Allows for type assertion to check for specific error types")
	fmt.Println("3. Enables more specific error handling")

	waitForEnter()
}

// Step 4: Error Wrapping (Go 1.13+)
// ===============================

func wrapError(id string) error {
	if len(id) == 0 {
		return fmt.Errorf("invalid ID: empty string")
	}

	err := queryDatabase(id)
	if err != nil {
		// Using the %w verb wraps the error
		return fmt.Errorf("lookup failed: %w", err)
	}

	return nil
}

func queryDatabase(id string) error {
	// Simulate a database error
	return fmt.Errorf("database connection timeout")
}

func errorWrapping() {
	fmt.Println("\n=== ERROR WRAPPING (Go 1.13+) ===")
	fmt.Println("Error wrapping allows you to add context while preserving the original error.")
	fmt.Println("This helps in creating error chains and tracing error origins.")

	// Demonstrate wrapping
	fmt.Println("\nDemonstrating error wrapping:")
	err := wrapError("user123")
	if err != nil {
		fmt.Printf("Top-level error: %v\n", err)

		// Unwrap once
		fmt.Println("\nUnwrapping the error:")
		unwrapped := errors.Unwrap(err)
		if unwrapped != nil {
			fmt.Printf("Unwrapped error: %v\n", unwrapped)
		}

		// Check if the error contains a specific error
		if errors.Is(err, errors.New("database connection timeout")) {
			fmt.Println("\nThe error chain contains 'database connection timeout'")
		}
	}

	fmt.Println("\nError wrapping pattern:")
	fmt.Println("err := someFunction()")
	fmt.Println("if err != nil {")
	fmt.Println("    return fmt.Errorf(\"context info: %w\", err)")
	fmt.Println("}")

	fmt.Println("\nChecking for specific errors in the chain:")
	fmt.Println("if errors.Is(err, specificError) {")
	fmt.Println("    // Handle the specific error")
	fmt.Println("}")

	waitForEnter()
}

// Step 5: Panic and Recovery
// ========================

func demoFunction() {
	fmt.Println("Step 1: Before the panic")
	panic("simulated panic situation")
	fmt.Println("This line will never execute")
}

func safeFunction() (err error) {
	// Set up the recover deferred function
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
			// Convert panic to an error
			err = fmt.Errorf("something went wrong: %v", r)
		}
	}()

	// Call the function that might panic
	demoFunction()

	return nil
}

func panicAndRecovery() {
	fmt.Println("\n=== PANIC AND RECOVERY ===")
	fmt.Println("Panic is for exceptional cases where normal error handling isn't appropriate.")
	fmt.Println("Recovery allows you to catch a panic and convert it to an error.")

	fmt.Println("\nPanics are for:")
	fmt.Println("1. Programming errors (nil pointer, index out of bounds)")
	fmt.Println("2. Truly exceptional situations")
	fmt.Println("3. When continuing would be dangerous")

	fmt.Println("\nRecovery pattern:")
	fmt.Println("defer func() {")
	fmt.Println("    if r := recover(); r != nil {")
	fmt.Println("        // Handle panic")
	fmt.Println("        err = fmt.Errorf(\"panic occurred: %v\", r)")
	fmt.Println("    }")
	fmt.Println("}()")

	fmt.Println("\nDemonstrating panic and recovery:")
	err := safeFunction()
	if err != nil {
		fmt.Printf("Function returned error: %v\n", err)
	}

	fmt.Println("\nNote: The program continues running after recovery")

	waitForEnter()
}

// Step 6: Context Package for Cancellation
// =====================================

func longRunningOperation(ctx context.Context) error {
	// Create a channel to signal completion
	done := make(chan struct{})

	// Run the slow operation in a goroutine
	go func() {
		// Simulate work
		time.Sleep(2 * time.Second)
		close(done)
	}()

	// Wait for completion or cancellation
	select {
	case <-done:
		fmt.Println("Operation completed successfully")
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func contextCancellation() {
	fmt.Println("\n=== CONTEXT-BASED CANCELLATION ===")
	fmt.Println("Go's context package provides a standard way to handle cancellation.")

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel() // Always call cancel to avoid resource leaks

	fmt.Println("Starting a long operation with a 1-second timeout...")
	err := longRunningOperation(ctx)
	if err != nil {
		fmt.Printf("Operation error: %v\n", err)

		// Check specific context errors
		if errors.Is(err, context.DeadlineExceeded) {
			fmt.Println("The operation timed out!")
		} else if errors.Is(err, context.Canceled) {
			fmt.Println("The operation was canceled!")
		}
	}

	fmt.Println("\nContext cancellation pattern:")
	fmt.Println("ctx, cancel := context.WithTimeout(context.Background(), timeout)")
	fmt.Println("defer cancel()")
	fmt.Println("select {")
	fmt.Println("case <-ctx.Done():")
	fmt.Println("    return ctx.Err()")
	fmt.Println("case result := <-resultChan:")
	fmt.Println("    return result")
	fmt.Println("}")

	waitForEnter()
}

// Step 7: Error Handling in Concurrent Operations
// ===========================================

func concurrentTask(id int) error {
	// Simulate work with random success/failure
	time.Sleep(time.Duration(500+id*100) * time.Millisecond)

	// Randomly fail some tasks
	if id%2 == 1 {
		return fmt.Errorf("task %d failed", id)
	}

	return nil
}

func concurrentErrorHandling() {
	fmt.Println("\n=== ERROR HANDLING IN CONCURRENT OPERATIONS ===")
	fmt.Println("Managing errors across multiple goroutines requires careful coordination.")

	// Approach 1: Collect all errors
	fmt.Println("\n1. Collecting all errors from multiple goroutines:")
	var wg sync.WaitGroup
	var mu sync.Mutex
	var errs []error

	for i := 1; i <= 4; i++ {
		wg.Add(1)
		id := i

		go func() {
			defer wg.Done()

			err := concurrentTask(id)
			if err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
			} else {
				fmt.Printf("Task %d completed successfully\n", id)
			}
		}()
	}

	wg.Wait()

	// Report errors
	if len(errs) > 0 {
		fmt.Println("\nErrors encountered:")
		for _, err := range errs {
			fmt.Printf("- %v\n", err)
		}
	}

	// Approach 2: Fail fast (first error cancels others)
	fmt.Println("\n2. Fail-fast approach (first error cancels other operations):")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var firstErr error
	var firstErrMu sync.Mutex

	for i := 5; i <= 8; i++ {
		wg.Add(1)
		id := i

		go func() {
			defer wg.Done()

			// Check if already cancelled
			select {
			case <-ctx.Done():
				fmt.Printf("Task %d cancelled\n", id)
				return
			default:
				// Continue with the task
			}

			err := concurrentTask(id)
			if err != nil {
				// Save the first error and cancel others
				firstErrMu.Lock()
				if firstErr == nil {
					firstErr = err
					cancel() // Cancel all other operations
				}
				firstErrMu.Unlock()
			} else {
				fmt.Printf("Task %d completed successfully\n", id)
			}
		}()
	}

	wg.Wait()

	if firstErr != nil {
		fmt.Printf("\nOperation failed with error: %v\n", firstErr)
		fmt.Println("Other operations were cancelled")
	}

	waitForEnter()
}

// Step 8: Sentinel Errors
// =====================

// Define sentinel errors (predefined error values)
var (
	ErrNotFound     = errors.New("not found")
	ErrUnauthorized = errors.New("unauthorized access")
	ErrInvalidInput = errors.New("invalid input")
)

func lookupUser(id string) (string, error) {
	// Check for empty ID
	if id == "" {
		return "", ErrInvalidInput
	}

	// Simulate user lookup
	if id == "admin" {
		return "Admin User", nil
	}

	// User not found
	return "", ErrNotFound
}

func sentinelErrors() {
	fmt.Println("\n=== SENTINEL ERRORS ===")
	fmt.Println("Sentinel errors are predefined error values that can be checked with ==")

	fmt.Println("\nDefining sentinel errors:")
	fmt.Println("var (")
	fmt.Println("    ErrNotFound = errors.New(\"not found\")")
	fmt.Println("    ErrUnauthorized = errors.New(\"unauthorized access\")")
	fmt.Println(")")

	fmt.Println("\nUsing sentinel errors in practice:")
	testIDs := []string{"", "admin", "user123"}
	for _, id := range testIDs {
		user, err := lookupUser(id)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				fmt.Printf("User with ID '%s' not found in the system\n", id)
			} else if errors.Is(err, ErrInvalidInput) {
				fmt.Println("Please provide a valid user ID")
			} else {
				fmt.Printf("Error looking up user: %v\n", err)
			}
		} else {
			fmt.Printf("Found user: %s\n", user)
		}
	}

	fmt.Println("\nAdvantages of sentinel errors:")
	fmt.Println("1. Allows for specific error checking using errors.Is()")
	fmt.Println("2. Makes expected errors part of your API contract")
	fmt.Println("3. Provides clearer error handling semantics")

	waitForEnter()
}

// Step 9: Best Practices
// ====================

func bestPractices() {
	fmt.Println("\n=== GO ERROR HANDLING BEST PRACTICES ===")

	fmt.Println("\n1. Be explicit about error handling")
	fmt.Println("   ✓ Always check errors returned by functions")
	fmt.Println("   ✓ Don't use _ to ignore errors unless you have a good reason")

	fmt.Println("\n2. Add context to errors")
	fmt.Println("   ✓ Wrap errors with additional information when returning them")
	fmt.Println("   ✓ Include relevant details about what operation failed")

	fmt.Println("\n3. Return errors, don't just log them")
	fmt.Println("   ✓ Let the caller decide how to handle the error")
	fmt.Println("   ✓ Don't force specific handling by logging in low-level functions")

	fmt.Println("\n4. Keep error handling close to error checking")
	fmt.Println("   ✓ Handle errors immediately after checking for them")
	fmt.Println("   ✓ Avoid long if-else chains with error checks")

	fmt.Println("\n5. Use custom error types for domain-specific errors")
	fmt.Println("   ✓ Create error types that carry contextual information")
	fmt.Println("   ✓ Implement the error interface for custom types")

	fmt.Println("\n6. Use errors.Is() and errors.As() for error checking (Go 1.13+)")
	fmt.Println("   ✓ Check for specific errors in a chain with errors.Is()")
	fmt.Println("   ✓ Extract error types with errors.As()")

	fmt.Println("\n7. Use panic only for unrecoverable situations")
	fmt.Println("   ✓ Use errors for expected failure cases")
	fmt.Println("   ✓ Reserve panic for programming errors and exceptional conditions")

	fmt.Println("\n8. Always include context cancellation in concurrent operations")
	fmt.Println("   ✓ Make long-running operations cancellable")
	fmt.Println("   ✓ Check for cancellation regularly in loops")

	waitForEnter()
}

// Helper function to pause briefly between sections
func waitForEnter() {
	fmt.Println("\n--- Proceeding to next section ---")
	time.Sleep(1 * time.Second)
}

func main() {
	fmt.Println("=== GO ERROR HANDLING PATTERNS TUTORIAL ===")
	fmt.Println("This program demonstrates different error handling patterns in Go.")
	fmt.Println("We'll go through several examples step by step.")
	fmt.Println("\nStarting in 2 seconds...")
	time.Sleep(2 * time.Second)

	// Run each demonstration in sequence
	basicErrorHandling()
	createErrors()
	customErrorTypes()
	errorWrapping()
	panicAndRecovery()
	contextCancellation()
	concurrentErrorHandling()
	sentinelErrors()
	bestPractices()

	// Final message
	fmt.Println("\n=== TUTORIAL COMPLETE ===")
	fmt.Println("You've learned about the following error handling patterns in Go:")
	fmt.Println("1. Basic error handling with 'if err != nil'")
	fmt.Println("2. Creating errors with errors.New and fmt.Errorf")
	fmt.Println("3. Custom error types")
	fmt.Println("4. Error wrapping and unwrapping")
	fmt.Println("5. Panic and recovery")
	fmt.Println("6. Context-based cancellation")
	fmt.Println("7. Error handling in concurrent operations")
	fmt.Println("8. Sentinel errors")
	fmt.Println("9. Best practices for error handling")
	fmt.Println("\nHappy coding!")
}
