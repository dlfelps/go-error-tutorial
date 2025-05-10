package cmd

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// AddErrorGroupsCmd adds the error groups command to the root command
func AddErrorGroupsCmd(rootCmd *cobra.Command) {
	errorGroupCmd := &cobra.Command{
		Use:   "errgroup",
		Short: "Demonstrates error groups for concurrent error handling",
		Long: `
ERROR GROUPS IN GO
----------------
This command demonstrates using error groups for handling errors in concurrent operations.

Error groups allow you to:
1. Run multiple operations concurrently
2. Collect the first error that occurs
3. Cancel all operations when one fails
4. Wait for all operations to complete

EXAMPLE:
  goerrors errgroup    # Start interactive tutorial
`,
		Run: runErrorGroupsDemo,
	}

	rootCmd.AddCommand(errorGroupCmd)
}

func runErrorGroupsDemo(cmd *cobra.Command, args []string) {
	// Run the full interactive tutorial
	runErrorGroupsTutorial()
}

// simulateConcurrentOperations demonstrates using a simple error group pattern
func simulateConcurrentOperations() error {
	// Set up a context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Use wait group to wait for all workers
	var wg sync.WaitGroup

	// Use mutex to protect the error
	var mu sync.Mutex
	var firstErr error

	// Launch multiple workers
	for i := 1; i <= 3; i++ {
		workerId := i // Create a local copy to avoid closure problems

		wg.Add(1)
		go func() {
			defer wg.Done()

			// Run the worker
			err := simulateWorker(ctx, workerId)

			// If error occurs, store it and cancel the context
			if err != nil {
				mu.Lock()
				if firstErr == nil {
					firstErr = err
					cancel() // Cancel all other workers
				}
				mu.Unlock()
			}
		}()
	}

	// Wait for all workers to complete
	wg.Wait()

	// Return the first error that occurred, if any
	return firstErr
}

// simulateWorker simulates a worker that might fail
func simulateWorker(ctx context.Context, id int) error {
	// Initialize random number generator
	rand.Seed(time.Now().UnixNano() + int64(id))

	// Simulate different execution times
	workTime := time.Duration(500+rand.Intn(1000)) * time.Millisecond

	// Worker has a 30% chance of failing
	willFail := rand.Float32() < 0.3

	// Print worker starting
	color.Cyan("Worker %d: Starting (will take %v)...\n", id, workTime)

	// Simulate work with context cancellation checks
	select {
	case <-time.After(workTime):
		// Work completed, check if it should fail
		if willFail {
			err := fmt.Errorf("worker %d: encountered a critical error", id)
			color.Red("Worker %d: Failed with error: %v\n", id, err)
			return err
		}
		color.Green("Worker %d: Completed successfully\n", id)
		return nil
	case <-ctx.Done():
		color.Yellow("Worker %d: Cancelled due to context cancellation\n", id)
		return ctx.Err()
	}
}

// simulateMultipleErrors demonstrates collecting multiple errors
func simulateMultipleErrors() error {
	// Create a container for multiple errors
	var errList []error

	// Run operations that might generate errors
	for i := 1; i <= 3; i++ {
		// Simulate operation with 50% chance of failure
		if rand.Float32() < 0.5 {
			err := fmt.Errorf("operation %d failed", i)
			color.Red("Operation %d: Failed with error: %v\n", i, err)
			errList = append(errList, err)
		} else {
			color.Green("Operation %d: Completed successfully\n", i)
		}
	}

	// If any errors occurred, combine them
	if len(errList) > 0 {
		return fmt.Errorf("multiple errors: %s", joinErrors(errList))
	}

	return nil
}

// joinErrors combines multiple errors into a single error message
func joinErrors(errs []error) string {
	if len(errs) == 0 {
		return ""
	}

	errorMessages := make([]string, len(errs))
	for i, err := range errs {
		errorMessages[i] = err.Error()
	}

	return strings.Join(errorMessages, "; ")
}

// runErrorGroupsTutorial provides a step-by-step tutorial on error groups
func runErrorGroupsTutorial() {
	clearScreen()
	printTitle("Error Groups in Go")

	fmt.Println("Welcome to the interactive tutorial on error groups in Go!")
	fmt.Println()

	printSection("What are Error Groups?")
	fmt.Println("Error groups provide synchronization, error propagation, and context")
	fmt.Println("cancellation for groups of goroutines working on subtasks of a common task.")
	fmt.Println()
	fmt.Println("Error groups are especially useful when you need to:")
	fmt.Println("1. Run multiple operations concurrently")
	fmt.Println("2. Stop all operations when one fails")
	fmt.Println("3. Collect the first error that occurs")
	fmt.Println("4. Wait for all operations to complete")
	fmt.Println()

	pressEnterToContinue()

	printSection("Basic Pattern")
	fmt.Println("The basic pattern for implementing an error group:")
	color.Cyan("// Create a context, wait group, and error tracking")
	color.Cyan("ctx, cancel := context.WithCancel(context.Background())")
	color.Cyan("defer cancel()")
	color.Cyan("var wg sync.WaitGroup")
	color.Cyan("var mu sync.Mutex")
	color.Cyan("var firstErr error")
	color.Cyan("")
	color.Cyan("// Launch multiple goroutines")
	color.Cyan("for i := 0; i < 3; i++ {")
	color.Cyan("    id := i  // Local copy for closure")
	color.Cyan("    wg.Add(1)")
	color.Cyan("    go func() {")
	color.Cyan("        defer wg.Done()")
	color.Cyan("        err := doWork(ctx, id)")
	color.Cyan("        if err != nil {")
	color.Cyan("            mu.Lock()")
	color.Cyan("            if firstErr == nil {")
	color.Cyan("                firstErr = err")
	color.Cyan("                cancel() // Cancel other operations")
	color.Cyan("            }")
	color.Cyan("            mu.Unlock()")
	color.Cyan("        }")
	color.Cyan("    }()")
	color.Cyan("}")
	color.Cyan("")
	color.Cyan("// Wait for all goroutines to complete")
	color.Cyan("wg.Wait()")
	color.Cyan("// Handle the first error")
	color.Cyan("if firstErr != nil {")
	color.Cyan("    // Handle error")
	color.Cyan("}")
	fmt.Println()

	pressEnterToContinue()

	printSection("Demonstration: Concurrent Workers")
	fmt.Println("Let's see error groups in action with some simulated workers:")
	fmt.Println("- We'll launch 3 concurrent workers")
	fmt.Println("- Each has a 30% chance of failing")
	fmt.Println("- If one fails, the context will be cancelled for all")
	fmt.Println()

	color.Yellow("Starting workers...")
	rand.Seed(time.Now().UnixNano())
	err := simulateConcurrentOperations()
	if err != nil {
		color.Red("\nError group failed: %v\n", err)
		fmt.Println("Notice how other workers were cancelled after the first error!")
	} else {
		color.Green("\nAll workers completed successfully!\n")
	}

	fmt.Println()
	pressEnterToContinue()

	printSection("Collecting Multiple Errors")
	fmt.Println("Sometimes you want to collect multiple errors rather than stopping at the first one:")
	color.Cyan("var errList []error")
	color.Cyan("")
	color.Cyan("// Collect errors from operations")
	color.Cyan("for i := 0; i < 3; i++ {")
	color.Cyan("    if err := doOperation(i); err != nil {")
	color.Cyan("        errList = append(errList, err)")
	color.Cyan("    }")
	color.Cyan("}")
	color.Cyan("")
	color.Cyan("// Combine errors if needed")
	color.Cyan("if len(errList) > 0 {")
	color.Cyan("    return fmt.Errorf(\"multiple errors: %s\", joinErrors(errList))")
	color.Cyan("}")
	fmt.Println()

	fmt.Println("Let's see multiple error collection in action:")
	err = simulateMultipleErrors()
	if err != nil {
		color.Red("\nCombined errors: %v\n", err)
	} else {
		color.Green("\nAll operations successful!\n")
	}

	fmt.Println()
	pressEnterToContinue()

	printSection("Practical Example: Parallel Downloads")
	fmt.Println("A common use case is downloading multiple resources in parallel:")
	color.Cyan("func downloadFiles(urls []string) error {")
	color.Cyan("    // Set up cancellation and synchronization")
	color.Cyan("    ctx, cancel := context.WithCancel(context.Background())")
	color.Cyan("    defer cancel()")
	color.Cyan("    var wg sync.WaitGroup")
	color.Cyan("    var mu sync.Mutex")
	color.Cyan("    var firstErr error")
	color.Cyan("")
	color.Cyan("    // Launch a goroutine for each URL")
	color.Cyan("    for _, url := range urls {")
	color.Cyan("        url := url  // Create local copy for closure")
	color.Cyan("        wg.Add(1)")
	color.Cyan("        go func() {")
	color.Cyan("            defer wg.Done()")
	color.Cyan("            err := downloadFile(ctx, url)")
	color.Cyan("            if err != nil {")
	color.Cyan("                mu.Lock()")
	color.Cyan("                if firstErr == nil {")
	color.Cyan("                    firstErr = err")
	color.Cyan("                    cancel() // Cancel other downloads")
	color.Cyan("                }")
	color.Cyan("                mu.Unlock()")
	color.Cyan("            }")
	color.Cyan("        }()")
	color.Cyan("    }")
	color.Cyan("")
	color.Cyan("    // Wait for all downloads to complete")
	color.Cyan("    wg.Wait()")
	color.Cyan("    return firstErr")
	color.Cyan("}")
	fmt.Println()

	pressEnterToContinue()

	printSection("Best Practices")
	fmt.Println("1. Always use context cancellation to stop goroutines early")
	fmt.Println("2. Protect shared state (like errors) with a mutex")
	fmt.Println("3. Be cautious about error handling semantics - first error vs. all errors")
	fmt.Println("4. Consider resource limits when launching many goroutines")
	fmt.Println("5. Use proper error combining techniques for multiple errors")
	fmt.Println()

	pressEnterToContinue()

	printSection("Summary")
	fmt.Println("Error groups in Go provide:")
	fmt.Println("- Structured concurrency with error handling")
	fmt.Println("- Automatic cancellation when an operation fails")
	fmt.Println("- A clean way to wait for multiple operations")
	fmt.Println("- Integration with the context package for timeouts and cancellation")
	fmt.Println()
	fmt.Println("When combined with context and other error handling patterns,")
	fmt.Println("error groups make concurrent error handling in Go much cleaner and safer.")
	fmt.Println()

	color.Green("You've completed all the error handling tutorials!")
	color.Green("Run 'goerrors --help' to see all available commands again.")
	fmt.Println()
}
