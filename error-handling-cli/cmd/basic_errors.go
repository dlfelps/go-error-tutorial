package cmd

import (
        "fmt"
        "strconv"

        "github.com/fatih/color"
        "github.com/spf13/cobra"
)

// AddBasicErrorHandlingCmd adds the basic error handling command to the root command
func AddBasicErrorHandlingCmd(rootCmd *cobra.Command) {
        basicCmd := &cobra.Command{
                Use:   "basic [number1] [number2]",
                Short: "Demonstrates basic error handling patterns",
                Long: `
BASIC ERROR HANDLING IN GO
--------------------------
This command demonstrates the most fundamental error handling pattern in Go:
the "if err != nil" pattern.

In Go, functions that can fail typically return an error as their last return value.
It's the caller's responsibility to check if this error is nil (indicating success)
or non-nil (indicating failure).

EXAMPLE:
  goerrors basic 10 2    # Divides 10 by 2
  goerrors basic 10 0    # Attempts to divide by zero (will show error handling)
`,
                Args: cobra.MinimumNArgs(0),
                Run:  runBasicErrorDemo,
        }

        rootCmd.AddCommand(basicCmd)
}

func runBasicErrorDemo(cmd *cobra.Command, args []string) {
        // If no arguments provided, run the interactive learning mode
        if len(args) == 0 {
                runBasicErrorTutorial()
                return
        }

        // If arguments provided, use them for the demonstration
        if len(args) != 2 {
                fmt.Println("Error: This command requires exactly 2 numeric arguments")
                return
        }

        // Parse the arguments as integers
        a, err := strconv.Atoi(args[0])
        if err != nil {
                fmt.Printf("Error parsing first number: %v\n", err)
                return
        }

        b, err := strconv.Atoi(args[1])
        if err != nil {
                fmt.Printf("Error parsing second number: %v\n", err)
                return
        }

        // Perform the division and handle potential errors
        result, err := divide(a, b)
        if err != nil {
                // This demonstrates the basic error handling pattern in Go
                color.Red("Error occurred: %v\n", err)
                printErrorHandlingExplanation()
        } else {
                color.Green("Result of %d ÷ %d = %d\n", a, b, result)
        }
}

func divide(a, b int) (int, error) {
        if b == 0 {
                return 0, fmt.Errorf("cannot divide by zero")
        }
        return a / b, nil
}

// runBasicErrorTutorial provides a step-by-step tutorial on basic error handling
func runBasicErrorTutorial() {
        ClearScreen()
        PrintTitle("Basic Error Handling in Go")

        fmt.Println("Welcome to the interactive tutorial on basic error handling in Go!")
        fmt.Println()

        PrintSection("The Fundamentals")
        fmt.Println("Go has a unique approach to error handling compared to many other languages:")
        fmt.Println("1. Errors are just values (typically of the built-in 'error' interface type)")
        fmt.Println("2. Functions that can fail return an error as their last return value")
        fmt.Println("3. Callers are expected to check and handle these errors explicitly")
        fmt.Println()

        PressEnterToContinue()

        PrintSection("The Error Interface")
        fmt.Println("In Go, the error interface is very simple:")
        color.Cyan("type error interface {")
        color.Cyan("    Error() string")
        color.Cyan("}")
        fmt.Println("\nAny type that implements the Error() method is an error.")
        fmt.Println()

        PressEnterToContinue()

        PrintSection("The Basic Pattern")
        fmt.Println("Here's the most common error handling pattern in Go:")
        color.Cyan("result, err := someFunction()")
        color.Cyan("if err != nil {")
        color.Cyan("    // Handle the error")
        color.Cyan("    return err  // Or take other appropriate action")
        color.Cyan("}")
        color.Cyan("// Use result if no error occurred")
        fmt.Println()

        PressEnterToContinue()

        PrintSection("Creating Errors")
        fmt.Println("Go provides several ways to create errors:")
        fmt.Println()
        color.Cyan("// Using the errors package")
        color.Cyan("import \"errors\"")
        color.Cyan("err := errors.New(\"something went wrong\")")
        fmt.Println()
        color.Cyan("// Using fmt.Errorf() for formatted errors")
        color.Cyan("err := fmt.Errorf(\"invalid value: %d\", value)")
        fmt.Println()

        PressEnterToContinue()

        PrintSection("Practical Example")
        fmt.Println("Let's look at a simple division function:")
        color.Cyan("func divide(a, b int) (int, error) {")
        color.Cyan("    if b == 0 {")
        color.Cyan("        return 0, fmt.Errorf(\"cannot divide by zero\")")
        color.Cyan("    }")
        color.Cyan("    return a / b, nil")
        color.Cyan("}")
        fmt.Println()
        fmt.Println("And how we would use it:")
        color.Cyan("result, err := divide(10, 2)")
        color.Cyan("if err != nil {")
        color.Cyan("    fmt.Printf(\"Error: %v\\n\", err)")
        color.Cyan("    return")
        color.Cyan("}")
        color.Cyan("fmt.Printf(\"Result: %d\\n\", result)")
        fmt.Println()

        PressEnterToContinue()

        PrintSection("Try It Yourself")
        fmt.Println("You can try the divide function with:")
        color.Green("goerrors basic 10 2    # Should work fine")
        color.Green("goerrors basic 10 0    # Should show error handling")
        fmt.Println()

        PrintSection("Best Practices")
        fmt.Println("1. Always check errors returned from functions")
        fmt.Println("2. Provide context in error messages")
        fmt.Println("3. Don't discard errors silently (unless it's intentional)")
        fmt.Println("4. Consider wrapping errors for better context (covered in the 'wrap' command)")
        fmt.Println()

        PressEnterToContinue()

        PrintSection("Summary")
        fmt.Println("Basic error handling in Go is characterized by:")
        fmt.Println("- Explicit error checking with 'if err != nil'")
        fmt.Println("- Treating errors as regular values")
        fmt.Println("- Functions that can fail return an error as the last return value")
        fmt.Println("- The caller is responsible for checking and handling errors")
        fmt.Println()
        fmt.Println("This approach makes error handling very explicit and forces developers")
        fmt.Println("to think about failure cases, leading to more robust programs.")
        fmt.Println()

        color.Green("To continue learning, try the other commands:")
        color.Green("goerrors custom    # Learn about custom error types")
        color.Green("goerrors wrap      # Learn about error wrapping")
        fmt.Println()
}

// printErrorHandlingExplanation prints an explanation of what just happened
func printErrorHandlingExplanation() {
        fmt.Println("\nWhat just happened?")
        fmt.Println("------------------")
        fmt.Println("1. The divide() function detected a division by zero")
        fmt.Println("2. It returned an error with a descriptive message")
        fmt.Println("3. The calling function checked for the error using 'if err != nil'")
        fmt.Println("4. Since the error was not nil, the error handling branch was executed")
        fmt.Println("5. The program displayed the error message to the user")
        fmt.Println()
        fmt.Println("This is the most basic error handling pattern in Go!")
}
