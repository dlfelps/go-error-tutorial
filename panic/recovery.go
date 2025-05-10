package panic

import (
	"fmt"
	"runtime/debug"

	"github.com/pkg/errors"
)

// PanicError represents an error that was recovered from a panic
type PanicError struct {
	Panic     interface{} // The value passed to panic()
	StackTrace string     // Stack trace of the panic
}

// Error implements the error interface
func (e *PanicError) Error() string {
	return fmt.Sprintf("panic occurred: %v", e.Panic)
}

// ExecuteWithRecover runs a function with panic recovery
func ExecuteWithRecover(fn func() (string, error)) (result string, err error) {
	// Set up recovery
	defer func() {
		if r := recover(); r != nil {
			// Create a PanicError with the panic value and stack trace
			err = &PanicError{
				Panic:     r,
				StackTrace: string(debug.Stack()),
			}
			// Log the stack trace for debugging
			fmt.Printf("Recovered from panic: %v\nStack trace: %s\n", r, debug.Stack())
		}
	}()

	// Execute the function
	return fn()
}

// SomethingThatPanics is a function that will panic
func SomethingThatPanics() {
	// This function will panic with a division by zero
	a := 10
	b := 0
	fmt.Println(a / b) // This will cause a panic
}

// GetValueSafely safely accesses a slice with panic recovery
func GetValueSafely(slice []int, index int) (value int, err error) {
	// Set up recovery
	defer func() {
		if r := recover(); r != nil {
			err = &PanicError{
				Panic:     r,
				StackTrace: string(debug.Stack()),
			}
		}
	}()

	// Check bounds to avoid panic
	if index < 0 || index >= len(slice) {
		return 0, errors.Errorf("index out of bounds: %d (length: %d)", index, len(slice))
	}

	// Access the slice
	return slice[index], nil
}

// ProcessData processes data with safety checks and panic recovery
func ProcessData(data interface{}) (result string, err error) {
	// Set up recovery
	defer func() {
		if r := recover(); r != nil {
			err = &PanicError{
				Panic:     r,
				StackTrace: string(debug.Stack()),
			}
		}
	}()

	// Check for nil data
	if data == nil {
		return "", errors.New("data is nil")
	}

	// Type assertion with safety
	strData, ok := data.(string)
	if !ok {
		return "", errors.Errorf("expected string, got %T", data)
	}

	// Process the data
	if strData == "" {
		return "", errors.New("empty string")
	}

	return fmt.Sprintf("Processed: %s", strData), nil
}

// IsPanicError checks if the error is a panic error
func IsPanicError(err error) bool {
	var panicErr *PanicError
	return errors.As(err, &panicErr)
}

// SafeMapAccess safely accesses a map with panic recovery
func SafeMapAccess(m map[string]string, key string) (value string, exists bool, err error) {
	// Set up recovery for nil map
	defer func() {
		if r := recover(); r != nil {
			err = &PanicError{
				Panic:     r,
				StackTrace: string(debug.Stack()),
			}
		}
	}()

	// Check if map is nil
	if m == nil {
		return "", false, errors.New("map is nil")
	}

	// Access the map
	value, exists = m[key]
	return value, exists, nil
}

// SafeCall safely calls a function with panic recovery
func SafeCall(fn func()) (err error) {
	// Set up recovery
	defer func() {
		if r := recover(); r != nil {
			err = &PanicError{
				Panic:     r,
				StackTrace: string(debug.Stack()),
			}
		}
	}()

	// Call the function
	fn()
	return nil
}
