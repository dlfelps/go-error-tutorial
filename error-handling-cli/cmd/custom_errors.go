package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// AddCustomErrorsCmd adds the custom error types command to the root command
func AddCustomErrorsCmd(rootCmd *cobra.Command) {
	customCmd := &cobra.Command{
		Use:   "custom [email]",
		Short: "Demonstrates custom error types",
		Long: `
CUSTOM ERROR TYPES IN GO
-----------------------
This command demonstrates creating and using custom error types in Go.

Custom error types allow you to:
1. Include structured data in your errors
2. Type-assert errors to check for specific error types
3. Implement domain-specific error handling logic

EXAMPLE:
  goerrors custom                     # Start interactive tutorial
  goerrors custom user@example.com    # Test validation with a specific email
`,
		Args: cobra.MinimumNArgs(0),
		Run:  runCustomErrorsDemo,
	}

	rootCmd.AddCommand(customCmd)
}

func runCustomErrorsDemo(cmd *cobra.Command, args []string) {
	// If no arguments provided, run the interactive learning mode
	if len(args) == 0 {
		runCustomErrorsTutorial()
		return
	}

	// If an email is provided, validate it
	email := args[0]

	// Create a test user with the given email
	user := User{
		Username: "testuser",
		Email:    email,
		Age:      17, // Intentionally set to a value that will cause a validation error
	}

	// Validate the user
	err := validateUser(user)
	if err != nil {
		// Type assertion to check for ValidationError
		if validationErr, ok := err.(*ValidationError); ok {
			color.Red("Validation Error: Field '%s' - %s\n", validationErr.Field, validationErr.Message)
			printCustomErrorExplanation()
		} else {
			color.Red("Unknown error: %v\n", err)
		}
	} else {
		color.Green("User validation successful!\n")
	}
}

// User represents a user for validation
type User struct {
	Username string
	Email    string
	Age      int
}

// ValidationError is a custom error type for validation errors
type ValidationError struct {
	Field   string
	Message string
}

// Error implements the error interface
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error for field '%s': %s", e.Field, e.Message)
}

// validateUser validates a user and returns a custom error type if validation fails
func validateUser(user User) error {
	// Validate username
	if user.Username == "" {
		return &ValidationError{
			Field:   "username",
			Message: "username cannot be empty",
		}
	}

	// Validate email
	if !isValidEmail(user.Email) {
		return &ValidationError{
			Field:   "email",
			Message: "email is not valid",
		}
	}

	// Validate age
	if user.Age < 18 {
		return &ValidationError{
			Field:   "age",
			Message: "user must be at least 18 years old",
		}
	}

	return nil
}

// isValidEmail performs a simple email validation
func isValidEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

// runCustomErrorsTutorial provides a step-by-step tutorial on custom error types
func runCustomErrorsTutorial() {
	clearScreen()
	printTitle("Custom Error Types in Go")

	fmt.Println("Welcome to the interactive tutorial on custom error types in Go!")
	fmt.Println()

	printSection("Why Custom Error Types?")
	fmt.Println("While simple string errors are often enough, custom error types provide:")
	fmt.Println("1. Structured error data with fields")
	fmt.Println("2. Type-based error handling with type assertions")
	fmt.Println("3. The ability to implement behavior on errors")
	fmt.Println("4. Domain-specific error hierarchies")
	fmt.Println()

	pressEnterToContinue()

	printSection("Defining a Custom Error Type")
	fmt.Println("A custom error type can be any type that implements the error interface:")
	color.Cyan("type ValidationError struct {")
	color.Cyan("    Field   string")
	color.Cyan("    Message string")
	color.Cyan("}")
	color.Cyan("")
	color.Cyan("// Error implements the error interface")
	color.Cyan("func (e *ValidationError) Error() string {")
	color.Cyan("    return fmt.Sprintf(\"validation error for field '%s': %s\", e.Field, e.Message)")
	color.Cyan("}")
	fmt.Println()

	pressEnterToContinue()

	printSection("Using Custom Error Types")
	fmt.Println("When a function returns an error, you can check for specific error types:")
	color.Cyan("err := validateUser(user)")
	color.Cyan("if err != nil {")
	color.Cyan("    // Type assertion to check if it's a ValidationError")
	color.Cyan("    if validationErr, ok := err.(*ValidationError); ok {")
	color.Cyan("        fmt.Printf(\"Field: %s, Message: %s\\n\", validationErr.Field, validationErr.Message)")
	color.Cyan("    } else {")
	color.Cyan("        fmt.Printf(\"Unknown error: %v\\n\", err)")
	color.Cyan("    }")
	color.Cyan("}")
	fmt.Println()

	pressEnterToContinue()

	printSection("Practical Example")
	fmt.Println("Let's look at a user validation function:")
	color.Cyan("func validateUser(user User) error {")
	color.Cyan("    if user.Username == \"\" {")
	color.Cyan("        return &ValidationError{")
	color.Cyan("            Field:   \"username\",")
	color.Cyan("            Message: \"username cannot be empty\",")
	color.Cyan("        }")
	color.Cyan("    }")
	color.Cyan("    ")
	color.Cyan("    if !isValidEmail(user.Email) {")
	color.Cyan("        return &ValidationError{")
	color.Cyan("            Field:   \"email\",")
	color.Cyan("            Message: \"email is not valid\",")
	color.Cyan("        }")
	color.Cyan("    }")
	color.Cyan("    ")
	color.Cyan("    // More validation...")
	color.Cyan("    ")
	color.Cyan("    return nil  // No error if validation passes")
	color.Cyan("}")
	fmt.Println()

	pressEnterToContinue()

	printSection("Error Type Hierarchies")
	fmt.Println("You can create hierarchies of error types:")
	color.Cyan("type AppError struct {")
	color.Cyan("    Err error")
	color.Cyan("    Message string")
	color.Cyan("    Code int")
	color.Cyan("}")
	color.Cyan("")
	color.Cyan("func (e *AppError) Error() string {")
	color.Cyan("    return e.Message")
	color.Cyan("}")
	color.Cyan("")
	color.Cyan("func (e *AppError) Unwrap() error {")
	color.Cyan("    return e.Err")
	color.Cyan("}")
	fmt.Println()

	pressEnterToContinue()

	printSection("Error Sentinel Values")
	fmt.Println("Go also supports predefined error values (sentinel errors):")
	color.Cyan("var (")
	color.Cyan("    ErrNotFound = errors.New(\"not found\")")
	color.Cyan("    ErrPermissionDenied = errors.New(\"permission denied\")")
	color.Cyan(")")
	color.Cyan("")
	color.Cyan("func FindUser(id int) (*User, error) {")
	color.Cyan("    // ...")
	color.Cyan("    return nil, ErrNotFound")
	color.Cyan("}")
	color.Cyan("")
	color.Cyan("// Then check for specific errors:")
	color.Cyan("if err == ErrNotFound {")
	color.Cyan("    // Handle not found case")
	color.Cyan("}")
	fmt.Println()

	pressEnterToContinue()

	printSection("Try It Yourself")
	fmt.Println("You can try the validation with:")
	color.Green("goerrors custom invalid-email    # Should fail email validation")
	color.Green("goerrors custom user@example.com  # Should pass email validation but fail age validation")
	fmt.Println()

	printSection("Best Practices")
	fmt.Println("1. Use custom error types for domain-specific errors")
	fmt.Println("2. Include enough context in errors to be helpful")
	fmt.Println("3. Consider implementing the Unwrap() method for error chains")
	fmt.Println("4. Use error sentinel values for expected errors that don't need context")
	fmt.Println()

	pressEnterToContinue()

	printSection("Summary")
	fmt.Println("Custom error types in Go allow you to:")
	fmt.Println("- Include structured data in your errors")
	fmt.Println("- Create domain-specific error types")
	fmt.Println("- Handle errors based on their type")
	fmt.Println("- Build error hierarchies")
	fmt.Println()
	fmt.Println("This approach provides much more context and structure than simple string errors.")
	fmt.Println()

	color.Green("To continue learning, try the next command:")
	color.Green("goerrors wrap    # Learn about error wrapping")
	fmt.Println()
}

// printCustomErrorExplanation prints an explanation of what just happened
func printCustomErrorExplanation() {
	fmt.Println("\nWhat just happened?")
	fmt.Println("------------------")
	fmt.Println("1. The validateUser() function detected a validation issue")
	fmt.Println("2. It returned a custom ValidationError with field and message details")
	fmt.Println("3. The calling code used type assertion (err.(*ValidationError)) to check the error type")
	fmt.Println("4. Since it was a ValidationError, it extracted the field and message info")
	fmt.Println("5. This allowed for more precise handling than a generic string error would")
	fmt.Println()
	fmt.Println("Custom error types make your errors more informative and allow for type-based handling!")
}
