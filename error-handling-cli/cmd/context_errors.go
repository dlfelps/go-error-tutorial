package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// AddContextErrorsCmd adds the context-based error handling command to the root command
func AddContextErrorsCmd(rootCmd *cobra.Command) {
	contextCmd := &cobra.Command{
		Use:   "context",
		Short: "Demonstrates context-based error handling",
		Long: `
CONTEXT-BASED ERROR HANDLING IN GO
---------------------------------
This command demonstrates how to use Go's context package for error handling,
especially for cancellation, timeouts, and deadlines.

The context package provides:
1. A standard way to pass request-scoped values, deadlines, and cancellation signals
2. Context-specific error types for timeouts and cancellations
3. A hierarchical structure that propagates cancellation to child operations

EXAMPLE:
  goerrors context    # Start interactive tutorial
`,
		Run: runContextErrorsDemo,
	}

	rootCmd.AddCommand(contextCmd)
}

func runContextErrorsDemo(cmd *cobra.Command, args []string) {
	// Run the full interactive tutorial
	runContextErrorsTutorial()
}

// simulateSlowOperation demonstrates a slow operation that can be cancelled by a context
func simulateSlowOperation(ctx context.Context) error {
	// Create a channel to signal completion
	done := make(chan struct{})

	// Start the slow operation in a goroutine
	go func() {
		// Simulate work by sleeping
		time.Sleep(2 * time.Second)

		// Only send on the channel if we haven't been cancelled
		select {
		case <-ctx.Done():
			// Context was cancelled, don't send completion signal
			return
		default:
			// Send completion signal
			close(done)
		}
	}()

	// Wait for either completion or cancellation
	select {
	case <-done:
		color.Green("Operation completed successfully\n")
		return nil
	case <-ctx.Done():
		// Return the context's error
		return ctx.Err()
	}
}

// simulateTimeoutOperation demonstrates a timeout using context
func simulateTimeoutOperation() {
	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel() // Always defer the cancel function to avoid resource leaks

	// Try to run the operation with the timeout context
	err := simulateSlowOperation(ctx)
	if err != nil {
		// Check for specific context errors
		if err == context.DeadlineExceeded {
			color.Red("Operation timed out: %v\n", err)
		} else if err == context.Canceled {
			color.Yellow("Operation was cancelled: %v\n", err)
		} else {
			color.Red("Operation failed with error: %v\n", err)
		}
	}
}

// simulateCancellationOperation demonstrates manual cancellation using context
func simulateCancellationOperation() {
	// Create a cancellable context
	ctx, cancel := context.WithCancel(context.Background())

	// Start the operation in a goroutine
	go func() {
		err := simulateSlowOperation(ctx)
		if err != nil {
			fmt.Printf("Operation result: %v\n", err)
		}
	}()

	// Simulate deciding to cancel after a short delay
	time.Sleep(500 * time.Millisecond)
	fmt.Println("Cancelling operation...")
	cancel()

	// Give the operation time to handle the cancellation
	time.Sleep(100 * time.Millisecond)
}

// runContextErrorsTutorial provides a step-by-step tutorial on context-based error handling
func runContextErrorsTutorial() {
	clearScreen()
	printTitle("Context-Based Error Handling in Go")

	fmt.Println("Welcome to the interactive tutorial on context-based error handling in Go!")
	fmt.Println()

	printSection("What is the Context Package?")
	fmt.Println("The context package provides a way to carry deadlines, cancellation signals,")
	fmt.Println("and request-scoped values across API boundaries and between processes.")
	fmt.Println()
	fmt.Println("For error handling, context is especially useful for:")
	fmt.Println("1. Timeouts and deadlines")
	fmt.Println("2. Cancellation of operations")
	fmt.Println("3. Propagating cancellation to multiple goroutines")
	fmt.Println()

	pressEnterToContinue()

	printSection("Context Types")
	fmt.Println("The context package provides several context types:")
	color.Cyan("// The root of all contexts")
	color.Cyan("ctx := context.Background()")
	color.Cyan("")
	color.Cyan("// A context that's already cancelled")
	color.Cyan("ctx := context.TODO()")
	color.Cyan("")
	color.Cyan("// A context with a deadline")
	color.Cyan("ctx, cancel := context.WithDeadline(parentCtx, time.Now().Add(5*time.Second))")
	color.Cyan("defer cancel()  // Always call cancel to avoid resource leaks")
	color.Cyan("")
	color.Cyan("// A context with a timeout (shorthand for WithDeadline)")
	color.Cyan("ctx, cancel := context.WithTimeout(parentCtx, 5*time.Second)")
	color.Cyan("defer cancel()")
	color.Cyan("")
	color.Cyan("// A context that can be cancelled manually")
	color.Cyan("ctx, cancel := context.WithCancel(parentCtx)")
	color.Cyan("defer cancel()")
	fmt.Println()

	pressEnterToContinue()

	printSection("Context Error Types")
	fmt.Println("The context package defines two special error types:")
	color.Cyan("// Returned when a context's deadline passes")
	color.Cyan("context.DeadlineExceeded")
	color.Cyan("")
	color.Cyan("// Returned when a context is cancelled manually")
	color.Cyan("context.Canceled")
	color.Cyan("")
	color.Cyan("// Check for these errors like this:")
	color.Cyan("if err == context.DeadlineExceeded {")
	color.Cyan("    // Handle timeout")
	color.Cyan("} else if err == context.Canceled {")
	color.Cyan("    // Handle cancellation")
	color.Cyan("}")
	fmt.Println()

	pressEnterToContinue()

	printSection("Using Context for Cancellation")
	fmt.Println("A typical pattern for making operations cancellable:")
	color.Cyan("func doOperation(ctx context.Context) error {")
	color.Cyan("    // Check if already cancelled before starting")
	color.Cyan("    select {")
	color.Cyan("    case <-ctx.Done():")
	color.Cyan("        return ctx.Err()")
	color.Cyan("    default:")
	color.Cyan("        // Not cancelled, continue")
	color.Cyan("    }")
	color.Cyan("")
	color.Cyan("    // Start the operation")
	color.Cyan("    for {")
	color.Cyan("        // Do some work...")
	color.Cyan("")
	color.Cyan("        // Periodically check for cancellation")
	color.Cyan("        select {")
	color.Cyan("        case <-ctx.Done():")
	color.Cyan("            return ctx.Err()")
	color.Cyan("        default:")
	color.Cyan("            // Not cancelled, continue")
	color.Cyan("        }")
	color.Cyan("    }")
	color.Cyan("}")
	fmt.Println()

	pressEnterToContinue()

	printSection("Demonstration: Timeout")
	fmt.Println("Let's see a timeout in action. The following operation takes 2 seconds,")
	fmt.Println("but we'll give it a 1 second timeout:")
	fmt.Println()
	color.Yellow("Running operation with a 1 second timeout...")
	simulateTimeoutOperation()
	fmt.Println()

	pressEnterToContinue()

	printSection("Demonstration: Cancellation")
	fmt.Println("Now let's see manual cancellation in action. We'll start an operation")
	fmt.Println("and then cancel it after 500 milliseconds:")
	fmt.Println()
	color.Yellow("Running operation with manual cancellation...")
	simulateCancellationOperation()
	fmt.Println()

	pressEnterToContinue()

	printSection("Practical Example: HTTP Request")
	fmt.Println("A common use is making HTTP requests cancellable:")
	color.Cyan("func fetchURL(ctx context.Context, url string) ([]byte, error) {")
	color.Cyan("    // Create a request with the context")
	color.Cyan("    req, err := http.NewRequestWithContext(ctx, \"GET\", url, nil)")
	color.Cyan("    if err != nil {")
	color.Cyan("        return nil, err")
	color.Cyan("    }")
	color.Cyan("")
	color.Cyan("    // Send the request")
	color.Cyan("    resp, err := http.DefaultClient.Do(req)")
	color.Cyan("    if err != nil {")
	color.Cyan("        // Check if it was a context error")
	color.Cyan("        if ctxErr := ctx.Err(); ctxErr != nil {")
	color.Cyan("            return nil, fmt.Errorf(\"request failed due to %v: %w\", ctxErr, err)")
	color.Cyan("        }")
	color.Cyan("        return nil, err")
	color.Cyan("    }")
	color.Cyan("    defer resp.Body.Close()")
	color.Cyan("")
	color.Cyan("    // Read the response body with context awareness")
	color.Cyan("    // ... implementation details ...")
	color.Cyan("}")
	fmt.Println()

	pressEnterToContinue()

	printSection("Best Practices")
	fmt.Println("1. Always pass a context as the first parameter to functions that may block")
	fmt.Println("2. Always defer the cancel function to prevent resource leaks")
	fmt.Println("3. Check for context cancellation regularly in long-running operations")
	fmt.Println("4. Don't store contexts in structs; pass them explicitly")
	fmt.Println("5. Only use context.Background() at the highest level; otherwise pass down contexts")
	fmt.Println()

	pressEnterToContinue()

	printSection("Summary")
	fmt.Println("Context-based error handling in Go provides:")
	fmt.Println("- A standardized way to handle timeouts and cancellation")
	fmt.Println("- Clear error types for different cancellation reasons")
	fmt.Println("- A mechanism to propagate cancellation across API boundaries")
	fmt.Println("- A way to make concurrent operations cancellable")
	fmt.Println()
	fmt.Println("This pattern is essential for writing robust, responsive services in Go,")
	fmt.Println("especially when dealing with external resources or long-running operations.")
	fmt.Println()

	color.Green("To continue learning, try the next command:")
	color.Green("goerrors errgroup    # Learn about error groups for concurrent operations")
	fmt.Println()
}
