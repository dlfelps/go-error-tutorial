package models

import (
	"time"

	"github.com/pkg/errors"
)

// User represents a user in the system
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// UserValidator defines validation functions for a User
type UserValidator interface {
	Validate() error
}

// Validate checks if the user data is valid
func (u *User) Validate() error {
	// Create a multi-error to collect multiple validation errors
	var errs []error

	// Validate Username
	if u.Username == "" {
		errs = append(errs, errors.New("username cannot be empty"))
	} else if len(u.Username) < 3 {
		errs = append(errs, errors.New("username must be at least 3 characters long"))
	} else if len(u.Username) > 50 {
		errs = append(errs, errors.New("username must be less than 50 characters long"))
	}

	// Validate Email
	if u.Email == "" {
		errs = append(errs, errors.New("email cannot be empty"))
	} else if !isValidEmail(u.Email) {
		errs = append(errs, errors.New("email format is invalid"))
	}

	// If there are any errors, return them
	if len(errs) > 0 {
		return NewValidationError("user validation failed", errs)
	}

	return nil
}

// isValidEmail is a simple function to validate email format
func isValidEmail(email string) bool {
	// This is a very simplified check - in real code use a proper validation
	return len(email) > 5 && (email[len(email)-4:] == ".com" || email[len(email)-4:] == ".org")
}

// ValidationError represents multiple validation errors
type ValidationError struct {
	Message string
	Errors  []error
}

// Error implements the error interface
func (ve *ValidationError) Error() string {
	if len(ve.Errors) == 0 {
		return ve.Message
	}

	errMessages := make([]string, len(ve.Errors))
	for i, err := range ve.Errors {
		errMessages[i] = err.Error()
	}

	return ve.Message + ": " + strings.Join(errMessages, ", ")
}

// NewValidationError creates a new ValidationError
func NewValidationError(message string, errs []error) *ValidationError {
	return &ValidationError{
		Message: message,
		Errors:  errs,
	}
}

// UserRepository defines operations for working with users
type UserRepository interface {
	FindByID(id int) (*User, error)
	FindByUsername(username string) (*User, error)
	Create(user *User) error
	Update(user *User) error
	Delete(id int) error
}
