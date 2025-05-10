# Go Error Handling Patterns Demo

[![Go Build and Test](https://github.com/username/repo-name/actions/workflows/go.yml/badge.svg)](https://github.com/username/repo-name/actions/workflows/go.yml)

This repository provides a comprehensive demonstration of error handling patterns and best practices in Go.

## Overview

Go has a unique approach to error handling, treating errors as values and encouraging explicit error checks. This project showcases various error handling techniques and patterns that will help you write more robust Go code.

## Features

The demo covers the following error handling patterns:

1. **Basic Error Handling**: The fundamental `if err != nil` pattern
2. **Creating Errors**: Using `errors.New` and `fmt.Errorf`
3. **Custom Error Types**: Creating structured error types with additional context
4. **Error Wrapping/Unwrapping**: Using Go 1.13+ error wrapping capabilities
5. **Panic and Recovery**: Handling exceptional situations with panic/recover
6. **Context-Based Cancellation**: Using the context package for timeouts and cancellation
7. **Concurrent Error Handling**: Managing errors across multiple goroutines
8. **Sentinel Errors**: Using predefined error values for specific error conditions
9. **Best Practices**: Guidelines for effective error handling in Go

## Getting Started

### Prerequisites

- Go 1.13 or later (for error wrapping features)

### Running the Demo

1. Clone this repository
2. Navigate to the project directory
3. Run the demonstration:

```bash
cd error-demo
go run main.go
```

The program will automatically step through each error handling pattern, showing both the code examples and their actual execution.

## Error Handling Patterns in Detail

### Basic Error Handling

Go functions that can fail typically return an error as their last return value. The caller checks if this error is nil to determine if the operation succeeded.

```go
file, err := os.Open("file.txt")
if err != nil {
    // Handle the error
    return err
}
// Use the file
```

### Custom Error Types

By implementing the `error` interface, you can create custom error types that include additional context:

```go
type ValidationError struct {
    Field   string
    Message string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("validation failed for field '%s': %s", e.Field, e.Message)
}
```

### Error Wrapping (Go 1.13+)

Error wrapping allows you to add context while preserving the original error:

```go
if err := doSomething(); err != nil {
    return fmt.Errorf("operation failed: %w", err)
}
```

## Best Practices

Some key best practices demonstrated include:

1. Always check errors returned by functions
2. Add context when wrapping errors
3. Use custom error types for domain-specific errors
4. Use sentinel errors for expected error conditions
5. Make operations cancellable with contexts
6. Use panic only for exceptional situations

## CI/CD

This project uses GitHub Actions to:

1. **Automatically build and test the code** on every push to the main branch
2. **Build binaries for multiple platforms** when a tag is pushed
   - Binaries are built for Linux, macOS, and Windows

### Workflow Files

- `.github/workflows/go.yml` - Builds and tests the code on push
- `.github/workflows/release.yml` - Builds binaries for multiple platforms when a tag is pushed

## License

This project is open source and available for learning and reference purposes.

## Acknowledgments

- The Go community for establishing these error handling patterns
- The Go team for their work on improving error handling in Go 1.13+