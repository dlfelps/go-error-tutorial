package dbops

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
	"github.com/pkg/errors"
)

// Custom errors for database operations
var (
	ErrUserNotFound      = errors.New("user not found")
	ErrDatabaseOperation = errors.New("database operation failed")
)

// User represents a user in the database
type User struct {
	ID        int
	Username  string
	Email     string
	CreatedAt time.Time
}

// InitDatabase initializes the database connection and sets up tables
func InitDatabase(dbPath string) (*sql.DB, error) {
	// Open the database connection
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open database connection")
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		// Make sure to close the database if we can't connect
		db.Close()
		return nil, errors.Wrap(err, "failed to ping database")
	}

	// Create tables if they don't exist
	if err := createTables(db); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

// createTables creates the necessary tables in the database
func createTables(db *sql.DB) error {
	// Define the users table
	createUsersTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY,
		username TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	// Execute the table creation statement
	_, err := db.Exec(createUsersTableSQL)
	if err != nil {
		return errors.Wrap(err, "failed to create users table")
	}

	return nil
}

// InsertUser inserts a new user into the database with transaction support
func InsertUser(ctx context.Context, db *sql.DB, user struct {
	ID       int
	Username string
	Email    string
}) error {
	// Start a transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}

	// Defer a rollback in case anything fails
	// If the transaction completes successfully, this rollback will be a no-op
	defer tx.Rollback()

	// Insert the user
	insertSQL := `INSERT INTO users (id, username, email) VALUES (?, ?, ?)`
	_, err = tx.ExecContext(ctx, insertSQL, user.ID, user.Username, user.Email)
	if err != nil {
		return errors.Wrap(err, "failed to insert user")
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}

// GetUser retrieves a user from the database by ID
func GetUser(ctx context.Context, db *sql.DB, id int) (*User, error) {
	// Create a context with timeout for the query
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Query the user
	row := db.QueryRowContext(queryCtx, `SELECT id, username, email, created_at FROM users WHERE id = ?`, id)

	// Scan the result into a User struct
	var user User
	var createdAtStr string
	err := row.Scan(&user.ID, &user.Username, &user.Email, &createdAtStr)
	if err != nil {
		// Check for no rows error
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, errors.Wrap(err, "failed to scan user row")
	}

	// Parse the created_at timestamp
	user.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse created_at timestamp")
	}

	return &user, nil
}

// UpdateUser updates an existing user in the database
func UpdateUser(ctx context.Context, db *sql.DB, user User) error {
	// Start a transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}
	defer tx.Rollback()

	// First check if the user exists
	var exists bool
	err = tx.QueryRowContext(ctx, "SELECT 1 FROM users WHERE id = ?", user.ID).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrUserNotFound
		}
		return errors.Wrap(err, "failed to check if user exists")
	}

	// Update the user
	updateSQL := `UPDATE users SET username = ?, email = ? WHERE id = ?`
	result, err := tx.ExecContext(ctx, updateSQL, user.Username, user.Email, user.ID)
	if err != nil {
		return errors.Wrap(err, "failed to update user")
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to get rows affected")
	}
	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}

// DeleteUser deletes a user from the database
func DeleteUser(ctx context.Context, db *sql.DB, id int) error {
	// Execute the delete with context
	result, err := db.ExecContext(ctx, `DELETE FROM users WHERE id = ?`, id)
	if err != nil {
		return errors.Wrap(err, "failed to delete user")
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to get rows affected")
	}
	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

// QueryUsersWithCancellation demonstrates query cancellation with context
func QueryUsersWithCancellation(ctx context.Context, db *sql.DB) ([]*User, error) {
	// Execute query with context
	rows, err := db.QueryContext(ctx, `SELECT id, username, email, created_at FROM users`)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query users")
	}
	defer rows.Close()

	// Collect users
	var users []*User
	for rows.Next() {
		// Check for context cancellation periodically
		select {
		case <-ctx.Done():
			return nil, errors.Wrap(ctx.Err(), "context cancelled during user query")
		default:
			// Continue processing
		}

		var user User
		var createdAtStr string
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &createdAtStr); err != nil {
			return nil, errors.Wrap(err, "failed to scan user row")
		}

		// Parse created_at timestamp
		user.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse created_at timestamp")
		}

		users = append(users, &user)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "error iterating over user rows")
	}

	return users, nil
}

// ExecuteInTransaction executes a function within a database transaction
func ExecuteInTransaction(ctx context.Context, db *sql.DB, fn func(*sql.Tx) error) error {
	// Start a transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}

	// Defer a rollback in case anything fails
	defer func() {
		// If the transaction was already committed, this will be a no-op
		tx.Rollback()
	}()

	// Execute the function
	if err := fn(tx); err != nil {
		// Transaction will be rolled back by the deferred function
		return err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}
