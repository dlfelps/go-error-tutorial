package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// AddPanicRecoveryCmd adds the panic and recovery command to the root command
func AddPanicRecoveryCmd(rootCmd *cobra.Command) {
	panicCmd := &cobra.Command{
		Use:   "panic",
		Short: "Demonstrates panic handling and recovery",
		Long: `
PANIC HANDLING AND RECOVERY IN GO
--------------------------------
This command demonstrates how to handle panics in Go using recover().

In Go:
- Panics are for unrecoverable situations and programming errors
- Regular error handling should be used for expected failure conditions
- recover() can be used to catch panics and convert them to errors
- defer is used with recover to ensure it runs even after a panic

EXAMPLE:
  goerrors panic    # Start interactive tutorial
`,
		Run: runPanicRecoveryDemo,
	}

	rootCmd.AddCommand(panicCmd)
}

func runPanicRecoveryDemo(cmd *cobra.Command, args []string) {
	// Run the full interactive tutorial
	runPanicRecoveryTutorial()
}

// safeDivide demonstrates panic recovery by attempting division
func safeDivide(a, b int) (result int, err error) {
	// Set up a deferred function to recover from panics
	defer func() {
		if r := recover(); r != nil {
			// Convert the panic to an error
			err = fmt.Errorf("panic in division operation: %v", r)
		}
	}()

	// This will panic if b is 0
	if b == 0 {
		panic("division by zero")
	}

	return a / b, nil
}

// nestedPanicExample demonstrates how panics propagate up the call stack
func nestedPanicExample() {
	defer func() {
		if r := recover(); r != nil {
			color.Yellow("Recovered in nestedPanicExample: %v", r)
		}
	}()

	fmt.Println("Starting nested function calls...")
	level1()
	fmt.Println("This line won't be reached if level1() panics")
}

func level1() {
	fmt.Println("In level1 function")
	level2()
	fmt.Println("Exiting level1 function")
}

func level2() {
	fmt.Println("In level2 function")
	level3()
	fmt.Println("Exiting level2 function")
}

func level3() {
	fmt.Println("In level3 function")
	panic("intentional panic in level3")
	fmt.Println("This line won't be reached")
}

// runPanicRecoveryTutorial provides a step-by-step tutorial on panic and recovery
func runPanicRecoveryTutorial() {
	clearScreen()
	printTitle("Panic Handling and Recovery in Go")

	fmt.Println("Welcome to the interactive tutorial on panic handling and recovery in Go!")
	fmt.Println()
	
	printSection("What are Panics?")
	fmt.Println("In Go, a panic is for exceptional situations where normal error handling isn't appropriate:")
	fmt.Println("1. Unrecoverable programming errors (nil pointer dereference, index out of range)")
	fmt.Println("2. Unexpected states that shouldn't happen in correctly written programs")
	fmt.Println("3. Situations where continuing execution would be dangerous")
	fmt.Println()
	fmt.Println("Panics are NOT for normal error conditions. For those, use regular error handling.")
	fmt.Println()
	
	pressEnterToContinue()
	
	printSection("How Panics Work")
	fmt.Println("When a panic occurs:")
	fmt.Println("1. Normal execution stops")
	fmt.Println("2. Deferred functions are executed")
	fmt.Println("3. The panic propagates up the call stack")
	fmt.Println("4. If not recovered, the program terminates with a stack trace")
	fmt.Println()
	
	pressEnterToContinue()
	
	printSection("Causing a Panic")
	fmt.Println("You can explicitly cause a panic:")
	color.Cyan("func dangerousOperation() {")
	color.Cyan("    // Something went terribly wrong")
	color.Cyan("    panic(\"critical error occurred\")")
	color.Cyan("}")
	fmt.Println()
	fmt.Println("Panics can also happen automatically:")
	color.Cyan("var ptr *int      // nil pointer")
	color.Cyan("*ptr = 42         // This causes a panic: nil pointer dereference")
	color.Cyan("")
	color.Cyan("arr := []int{1, 2, 3}")
	color.Cyan("value := arr[10]  // This causes a panic: index out of range")
	fmt.Println()
	
	pressEnterToContinue()
	
	printSection("Recovering from Panics")
	fmt.Println("The recover() function allows you to catch and handle panics:")
	color.Cyan("func doSomething() (err error) {")
	color.Cyan("    defer func() {")
	color.Cyan("        if r := recover(); r != nil {")
	color.Cyan("            // Convert the panic to an error")
	color.Cyan("            err = fmt.Errorf(\"panic occurred: %v\", r)")
	color.Cyan("        }")
	color.Cyan("    }()")
	color.Cyan("    ")
	color.Cyan("    // ... potentially panicking code ...")
	color.Cyan("    ")
	color.Cyan("    return nil")
	color.Cyan("}")
	fmt.Println()
	
	pressEnterToContinue()
	
	printSection("Practical Example: Safe Division")
	fmt.Println("Let's implement a division function that converts panics to errors:")
	color.Cyan("func safeDivide(a, b int) (result int, err error) {")
	color.Cyan("    // Set up a deferred function to recover from panics")
	color.Cyan("    defer func() {")
	color.Cyan("        if r := recover(); r != nil {")
	color.Cyan("            // Convert the panic to an error")
	color.Cyan("            err = fmt.Errorf(\"panic in division operation: %v\", r)")
	color.Cyan("        }")
	color.Cyan("    }()")
	color.Cyan("")
	color.Cyan("    // This will panic if b is 0")
	color.Cyan("    if b == 0 {")
	color.Cyan("        panic(\"division by zero\")")
	color.Cyan("    }")
	color.Cyan("")
	color.Cyan("    return a / b, nil")
	color.Cyan("}")
	fmt.Println()
	
	pressEnterToContinue()
	
	printSection("Demonstration")
	fmt.Println("Let's see our safeDivide function in action:")
	
	// Safe case
	result, err := safeDivide(10, 2)
	if err != nil {
		color.Red("Error: %v\n", err)
	} else {
		color.Green("Result of 10 ÷ 2 = %d\n", result)
	}
	
	// Panic case
	result, err = safeDivide(10, 0)
	if err != nil {
		color.Red("Error: %v\n", err)
	} else {
		color.Green("Result = %d\n", result)
	}
	
	fmt.Println()
	color.Yellow("Notice how the panic was converted to a regular error that we can handle!")
	fmt.Println()
	
	pressEnterToContinue()
	
	printSection("How Panics Propagate")
	fmt.Println("Panics propagate up the call stack until recovered:")
	color.Cyan("func main() {")
	color.Cyan("    // Set up recovery")
	color.Cyan("    defer func() {")
	color.Cyan("        if r := recover(); r != nil {")
	color.Cyan("            fmt.Println(\"Recovered:\", r)")
	color.Cyan("        }")
	color.Cyan("    }()")
	color.Cyan("    ")
	color.Cyan("    // Call a function that will eventually panic")
	color.Cyan("    level1()")
	color.Cyan("}")
	color.Cyan("")
	color.Cyan("func level1() {")
	color.Cyan("    level2()")
	color.Cyan("}")
	color.Cyan("")
	color.Cyan("func level2() {")
	color.Cyan("    level3()")
	color.Cyan("}")
	color.Cyan("")
	color.Cyan("func level3() {")
	color.Cyan("    panic(\"something went wrong\")")
	color.Cyan("}")
	fmt.Println()
	
	pressEnterToContinue()
	
	printSection("Demonstration: Nested Panic")
	fmt.Println("Let's see how a panic propagates through nested function calls:")
	fmt.Println()
	color.Yellow("Starting demonstration - watch the call stack unwind...")
	nestedPanicExample()
	fmt.Println()
	
	pressEnterToContinue()
	
	printSection("Stack Traces")
	fmt.Println("You can capture stack traces when recovering from panics:")
	color.Cyan("defer func() {")
	color.Cyan("    if r := recover(); r != nil {")
	color.Cyan("        fmt.Printf(\"Panic: %v\\n\", r)")
	color.Cyan("        fmt.Printf(\"Stack trace:\\n%s\\n\", debug.Stack())")
	color.Cyan("    }")
	color.Cyan("}()")
	fmt.Println()
	
	pressEnterToContinue()
	
	printSection("Best Practices")
	fmt.Println("1. Use panics only for truly exceptional conditions")
	fmt.Println("2. For expected errors, use error returns instead of panics")
	fmt.Println("3. Only recover from panics in high-level functions")
	fmt.Println("4. When recovering, log or return enough information to diagnose the issue")
	fmt.Println("5. Consider converting panics to errors at API boundaries")
	fmt.Println()
	
	pressEnterToContinue()
	
	printSection("HTTP Handler Example")
	fmt.Println("A common use of recover is in HTTP handlers:")
	color.Cyan("func safeHandler(handler http.HandlerFunc) http.HandlerFunc {")
	color.Cyan("    return func(w http.ResponseWriter, r *http.Request) {")
	color.Cyan("        defer func() {")
	color.Cyan("            if err := recover(); err != nil {")
	color.Cyan("                // Log the error and stack trace")
	color.Cyan("                log.Printf(\"Handler panic: %v\\n%s\", err, debug.Stack())")
	color.Cyan("                // Return a 500 error to the client")
	color.Cyan("                http.Error(w, \"Internal server error\", http.StatusInternalServerError)")
	color.Cyan("            }")
	color.Cyan("        }()")
	color.Cyan("        ")
	color.Cyan("        // Call the original handler")
	color.Cyan("        handler(w, r)")
	color.Cyan("    }")
	color.Cyan("}")
	fmt.Println()
	
	pressEnterToContinue()
	
	printSection("Summary")
	fmt.Println("In Go, panic and recover provide a mechanism for handling exceptional cases:")
	fmt.Println("- Panics are for unrecoverable errors and programmer mistakes")
	fmt.Println("- Regular error handling is for expected failure conditions")
	fmt.Println("- recover() can convert panics to errors for graceful handling")
	fmt.Println("- defer is essential for setting up recovery")
	fmt.Println()
	fmt.Println("Remember: panic should be rare in well-written Go code. Use error returns")
	fmt.Println("for normal error handling, and panics only for truly exceptional conditions.")
	fmt.Println()
	
	color.Green("To continue learning, try the next command:")
	color.Green("goerrors context    # Learn about context-based cancellation and error handling")
	fmt.Println()
}