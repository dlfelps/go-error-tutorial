package contextdemo

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// ContextError represents an error that occurred due to context cancellation
type ContextError struct {
	Operation string
	Err       error
}

// Error implements the error interface
func (e *ContextError) Error() string {
	return fmt.Sprintf("context error during %s: %v", e.Operation, e.Err)
}

// Unwrap returns the underlying error
func (e *ContextError) Unwrap() error {
	return e.Err
}

// ErrOperationCancelled is returned when an operation is cancelled
var ErrOperationCancelled = errors.New("operation was cancelled")

// ErrOperationTimedOut is returned when an operation times out
var ErrOperationTimedOut = errors.New("operation timed out")

// ExecuteWithContext demonstrates context cancellation
func ExecuteWithContext(ctx context.Context) error {
	// Create a channel to signal operation completion
	done := make(chan struct{})

	// Start a goroutine for the operation
	go func() {
		// Simulate a long-running operation
		time.Sleep(1 * time.Second)
		
		// Signal that the operation is complete
		close(done)
	}()

	// Wait for either the operation to complete or the context to be cancelled
	select {
	case <-done:
		// Operation completed successfully
		fmt.Println("Operation completed successfully")
		return nil
	case <-ctx.Done():
		// Context was cancelled
		return &ContextError{
			Operation: "execute",
			Err:       ctx.Err(),
		}
	}
}

// SlowOperation demonstrates a slow operation that can be cancelled
func SlowOperation(ctx context.Context) error {
	// Check if context is already cancelled before starting
	if ctx.Err() != nil {
		return &ContextError{
			Operation: "slow_operation_start",
			Err:       ctx.Err(),
		}
	}

	fmt.Println("Starting slow operation...")

	// Create a ticker to update progress
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	// Simulate work with progress updates
	for i := 0; i < 10; i++ {
		select {
		case <-ticker.C:
			// Continue with the operation
			fmt.Printf("Operation progress: %d%%\n", (i+1)*10)
		case <-ctx.Done():
			// Context was cancelled
			fmt.Printf("Operation cancelled at %d%% completion\n", (i+1)*10)
			
			// Check the specific context error
			if ctx.Err() == context.DeadlineExceeded {
				return &ContextError{
					Operation: "slow_operation",
					Err:       ErrOperationTimedOut,
				}
			}
			
			return &ContextError{
				Operation: "slow_operation",
				Err:       ErrOperationCancelled,
			}
		}
	}

	fmt.Println("Slow operation completed successfully")
	return nil
}

// ProcessWithTimeout processes data with a timeout
func ProcessWithTimeout(data string, timeout time.Duration) (string, error) {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Call process function with context
	return ProcessWithContext(ctx, data)
}

// ProcessWithContext processes data with context awareness
func ProcessWithContext(ctx context.Context, data string) (string, error) {
	// Check if context is already cancelled
	if ctx.Err() != nil {
		return "", &ContextError{
			Operation: "process_with_context",
			Err:       ctx.Err(),
		}
	}

	// Simulate a time-consuming process
	select {
	case <-time.After(500 * time.Millisecond):
		// Process completed successfully
		return "Processed: " + data, nil
	case <-ctx.Done():
		// Context was cancelled
		return "", &ContextError{
			Operation: "process_with_context",
			Err:       ctx.Err(),
		}
	}
}

// IsContextCancelled checks if the error was caused by context cancellation
func IsContextCancelled(err error) bool {
	// Unwrap the error if it's a ContextError
	var ctxErr *ContextError
	if errors.As(err, &ctxErr) {
		err = ctxErr.Err
	}

	// Check if it's a context.Canceled error
	return errors.Is(err, context.Canceled)
}

// IsContextTimeout checks if the error was caused by context timeout
func IsContextTimeout(err error) bool {
	// Unwrap the error if it's a ContextError
	var ctxErr *ContextError
	if errors.As(err, &ctxErr) {
		err = ctxErr.Err
	}

	// Check if it's a deadline exceeded error
	return errors.Is(err, context.DeadlineExceeded) || errors.Is(err, ErrOperationTimedOut)
}
