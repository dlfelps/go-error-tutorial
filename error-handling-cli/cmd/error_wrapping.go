package cmd

import (
        "errors"
        "fmt"
        "strings"

        "github.com/fatih/color"
        "github.com/spf13/cobra"
)

// AddErrorWrappingCmd adds the error wrapping command to the root command
func AddErrorWrappingCmd(rootCmd *cobra.Command) {
        wrapCmd := &cobra.Command{
                Use:   "wrap",
                Short: "Demonstrates error wrapping and unwrapping",
                Long: `
ERROR WRAPPING IN GO
------------------
This command demonstrates how to wrap and unwrap errors in Go.

Error wrapping allows you to:
1. Add context to errors as they propagate up the call stack
2. Maintain the original error information
3. Create a chain of errors showing what went wrong at each level

Starting with Go 1.13, the standard library provides:
- fmt.Errorf() with %w verb for wrapping errors
- errors.Is() to check if an error is a specific error in the chain
- errors.As() to extract a specific error type from the chain
- errors.Unwrap() to get the wrapped error

EXAMPLE:
  goerrors wrap    # Start interactive tutorial
`,
                Run: runErrorWrappingDemo,
        }

        rootCmd.AddCommand(wrapCmd)
}

func runErrorWrappingDemo(cmd *cobra.Command, args []string) {
        // Run the full interactive tutorial
        runErrorWrappingTutorial()
}

// simulateErrorChain demonstrates error wrapping through multiple function calls
func simulateErrorChain() error {
        err := levelThree()
        if err != nil {
                return fmt.Errorf("top-level operation failed: %w", err)
        }
        return nil
}

func levelThree() error {
        err := levelTwo()
        if err != nil {
                return fmt.Errorf("level three processing failed: %w", err)
        }
        return nil
}

func levelTwo() error {
        err := levelOne()
        if err != nil {
                return fmt.Errorf("level two validation failed: %w", err)
        }
        return nil
}

func levelOne() error {
        // Base error
        return errors.New("level one encountered a critical error")
}

// runErrorWrappingTutorial provides a step-by-step tutorial on error wrapping
func runErrorWrappingTutorial() {
        ClearScreen()
        PrintTitle("Error Wrapping in Go")

        fmt.Println("Welcome to the interactive tutorial on error wrapping in Go!")
        fmt.Println()

        PrintSection("What is Error Wrapping?")
        fmt.Println("Error wrapping is a technique that allows you to:")
        fmt.Println("1. Add context to errors as they propagate up the call stack")
        fmt.Println("2. Preserve the original error information")
        fmt.Println("3. Create a chain of errors showing what happened at each level")
        fmt.Println()

        PressEnterToContinue()

        PrintSection("Basic Error Wrapping")
        fmt.Println("Starting with Go 1.13, you can wrap errors using fmt.Errorf() with the %w verb:")
        color.Cyan("originalErr := errors.New(\"database connection failed\")")
        color.Cyan("wrappedErr := fmt.Errorf(\"query failed: %w\", originalErr)")
        fmt.Println()
        fmt.Println("The wrapped error contains both the new context and the original error.")
        fmt.Println()

        PressEnterToContinue()

        PrintSection("Checking Wrapped Errors")
        fmt.Println("You can check if an error in the chain is a specific error:")
        color.Cyan("// Define a sentinel error")
        color.Cyan("var ErrNotFound = errors.New(\"not found\")")
        color.Cyan("")
        color.Cyan("// Later in the code...")
        color.Cyan("if errors.Is(err, ErrNotFound) {")
        color.Cyan("    // Handle not found case")
        color.Cyan("}")
        fmt.Println()

        PressEnterToContinue()

        PrintSection("Extracting Wrapped Errors")
        fmt.Println("You can also extract specific error types from the chain:")
        color.Cyan("var dbErr *DatabaseError")
        color.Cyan("if errors.As(err, &dbErr) {")
        color.Cyan("    // dbErr is now the DatabaseError in the chain")
        color.Cyan("    fmt.Printf(\"Database error code: %d\\n\", dbErr.Code)")
        color.Cyan("}")
        fmt.Println()

        PressEnterToContinue()

        PrintSection("Practical Example")
        fmt.Println("Let's look at a chain of function calls with error wrapping:")

        color.Cyan("func simulateErrorChain() error {")
        color.Cyan("    err := levelThree()")
        color.Cyan("    if err != nil {")
        color.Cyan("        return fmt.Errorf(\"top-level operation failed: %w\", err)")
        color.Cyan("    }")
        color.Cyan("    return nil")
        color.Cyan("}")
        color.Cyan("")
        color.Cyan("func levelThree() error {")
        color.Cyan("    err := levelTwo()")
        color.Cyan("    if err != nil {")
        color.Cyan("        return fmt.Errorf(\"level three processing failed: %w\", err)")
        color.Cyan("    }")
        color.Cyan("    return nil")
        color.Cyan("}")
        color.Cyan("")
        color.Cyan("// ... and so on")
        fmt.Println()

        PressEnterToContinue()

        PrintSection("Demonstration")
        fmt.Println("Let's see error wrapping in action:")

        // Generate an error chain
        err := simulateErrorChain()

        // Print the wrapped error
        color.Red("Wrapped error: %v\n", err)

        // Demonstrate errors.Is
        originalErr := errors.New("level one encountered a critical error")
        if errors.Is(err, originalErr) {
                color.Green("✓ errors.Is() confirms this error chain contains our original error\n")
        }

        fmt.Println()
        color.Yellow("Notice how each layer adds context, while preserving the original error!")
        fmt.Println()

        PressEnterToContinue()

        PrintSection("Unwrapping Errors")
        fmt.Println("You can manually unwrap errors using errors.Unwrap():")
        color.Cyan("func printErrorChain(err error) {")
        color.Cyan("    for err != nil {")
        color.Cyan("        fmt.Printf(\"Error: %v\\n\", err)")
        color.Cyan("        err = errors.Unwrap(err)")
        color.Cyan("    }")
        color.Cyan("}")
        fmt.Println()

        fmt.Println("Let's print our error chain:")
        printErrorChain(err)
        fmt.Println()

        PressEnterToContinue()

        PrintSection("Custom Unwrap Method")
        fmt.Println("You can implement the Unwrap method on your custom error types:")
        color.Cyan("type QueryError struct {")
        color.Cyan("    Query string")
        color.Cyan("    Err   error")
        color.Cyan("}")
        color.Cyan("")
        color.Cyan("func (e *QueryError) Error() string {")
        color.Cyan("    return fmt.Sprintf(\"query '%s' failed: %v\", e.Query, e.Err)")
        color.Cyan("}")
        color.Cyan("")
        color.Cyan("func (e *QueryError) Unwrap() error {")
        color.Cyan("    return e.Err")
        color.Cyan("}")
        fmt.Println()

        PressEnterToContinue()

        PrintSection("Best Practices")
        fmt.Println("1. Add context when wrapping errors to clarify what operation failed")
        fmt.Println("2. Use errors.Is() to check for specific errors in the chain")
        fmt.Println("3. Use errors.As() to extract specific error types from the chain")
        fmt.Println("4. Implement the Unwrap() method on custom error types")
        fmt.Println("5. Don't wrap errors unnecessarily if no context is being added")
        fmt.Println()

        PressEnterToContinue()

        PrintSection("Summary")
        fmt.Println("Error wrapping in Go allows you to:")
        fmt.Println("- Add context at each level of your call stack")
        fmt.Println("- Preserve the original error information")
        fmt.Println("- Check for specific errors in the chain")
        fmt.Println("- Extract specific error types from the chain")
        fmt.Println()
        fmt.Println("This leads to more informative and useful error messages while")
        fmt.Println("still allowing precise error checking in your code.")
        fmt.Println()

        color.Green("To continue learning, try the next command:")
        color.Green("goerrors panic    # Learn about panic handling and recovery")
        fmt.Println()
}

// printErrorChain prints all errors in an error chain
func printErrorChain(err error) {
        indent := 0
        for err != nil {
                // Print the current error with indentation
                spaces := strings.Repeat("  ", indent)
                color.Cyan("%s→ %v", spaces, err)

                // Unwrap for the next iteration
                err = errors.Unwrap(err)
                indent++
        }
}
