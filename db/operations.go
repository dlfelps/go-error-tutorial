package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/pkg/errors"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// User represents a database user entity
type User struct {
	ID        int
	Name      string
	Email     string
	CreatedAt time.Time
}

// DBError represents a database operation error
type DBError struct {
	Op  string
	SQL string
	Err error
}

// Error implements the error interface
func (e *DBError) Error() string {
	return fmt.Sprintf("database error during %s: %v", e.Op, e.Err)
}

// Unwrap returns the underlying error
func (e *DBError) Unwrap() error {
	return e.Err
}

// ErrNotFound is returned when a record is not found
var ErrNotFound = errors.New("record not found")

// OpenDatabase opens a database connection with proper error handling
func OpenDatabase(ctx context.Context, dbPath string) (*sql.DB, error) {
	// Open database connection
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, &DBError{
			Op:  "open",
			Err: errors.Wrap(err, "failed to open database"),
		}
	}

	// Set connection pool parameters
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	// Verify connection with ping
	err = db.PingContext(ctx)
	if err != nil {
		// Close the database if ping fails
		db.Close()
		return nil, &DBError{
			Op:  "ping",
			Err: errors.Wrap(err, "failed to ping database"),
		}
	}

	return db, nil
}

// CreateSchema creates the database schema
func CreateSchema(ctx context.Context, db *sql.DB) error {
	// SQL to create the users table
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	// Execute the SQL
	_, err := db.ExecContext(ctx, createTableSQL)
	if err != nil {
		return &DBError{
			Op:  "create_schema",
			SQL: createTableSQL,
			Err: errors.Wrap(err, "failed to create users table"),
		}
	}

	return nil
}

// InsertUser inserts a new user into the database
func InsertUser(ctx context.Context, db *sql.DB, name, email string) (int64, error) {
	// Validate input
	if name == "" || email == "" {
		return 0, errors.New("name and email are required")
	}

	// Prepare insert statement
	insertSQL := "INSERT INTO users (name, email) VALUES (?, ?)"
	
	// Begin a transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return 0, &DBError{
			Op:  "insert_user_begin_tx",
			Err: errors.Wrap(err, "failed to begin transaction"),
		}
	}
	
	// Ensure transaction is rolled back if function returns with error
	defer func() {
		if tx != nil {
			tx.Rollback()
		}
	}()

	// Execute insert
	result, err := tx.ExecContext(ctx, insertSQL, name, email)
	if err != nil {
		return 0, &DBError{
			Op:  "insert_user",
			SQL: insertSQL,
			Err: errors.Wrap(err, "failed to insert user"),
		}
	}

	// Get the inserted ID
	id, err := result.LastInsertId()
	if err != nil {
		return 0, &DBError{
			Op:  "insert_user_last_id",
			Err: errors.Wrap(err, "failed to get last insert ID"),
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return 0, &DBError{
			Op:  "insert_user_commit",
			Err: errors.Wrap(err, "failed to commit transaction"),
		}
	}
	
	// Set tx to nil to prevent rollback in defer
	tx = nil

	return id, nil
}

// GetUser retrieves a user by ID
func GetUser(ctx context.Context, db *sql.DB, id int64) (*User, error) {
	// Prepare query
	querySQL := "SELECT id, name, email, created_at FROM users WHERE id = ?"
	
	// Execute query with context
	row := db.QueryRowContext(ctx, querySQL, id)
	
	// Scan results into User struct
	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			// Specific error for "not found" case
			return nil, ErrNotFound
		}
		return nil, &DBError{
			Op:  "get_user",
			SQL: querySQL,
			Err: errors.Wrapf(err, "failed to scan user with id %d", id),
		}
	}

	return &user, nil
}

// UpdateUser updates a user's information
func UpdateUser(ctx context.Context, db *sql.DB, id int64, name, email string) error {
	// Validate input
	if id <= 0 {
		return errors.New("invalid user ID")
	}
	if name == "" && email == "" {
		return errors.New("no fields to update")
	}

	// Determine which fields to update
	updateSQL := "UPDATE users SET "
	args := make([]interface{}, 0)
	
	if name != "" {
		updateSQL += "name = ?"
		args = append(args, name)
		if email != "" {
			updateSQL += ", "
		}
	}
	
	if email != "" {
		updateSQL += "email = ?"
		args = append(args, email)
	}
	
	updateSQL += " WHERE id = ?"
	args = append(args, id)

	// Execute update
	result, err := db.ExecContext(ctx, updateSQL, args...)
	if err != nil {
		return &DBError{
			Op:  "update_user",
			SQL: updateSQL,
			Err: errors.Wrapf(err, "failed to update user with id %d", id),
		}
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return &DBError{
			Op:  "update_user_rows_affected",
			Err: errors.Wrap(err, "failed to get rows affected"),
		}
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// DeleteUser deletes a user by ID
func DeleteUser(ctx context.Context, db *sql.DB, id int64) error {
	// Validate input
	if id <= 0 {
		return errors.New("invalid user ID")
	}

	// Prepare delete statement
	deleteSQL := "DELETE FROM users WHERE id = ?"
	
	// Execute delete
	result, err := db.ExecContext(ctx, deleteSQL, id)
	if err != nil {
		return &DBError{
			Op:  "delete_user",
			SQL: deleteSQL,
			Err: errors.Wrapf(err, "failed to delete user with id %d", id),
		}
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return &DBError{
			Op:  "delete_user_rows_affected",
			Err: errors.Wrap(err, "failed to get rows affected"),
		}
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// ExecuteTransaction demonstrates transaction with error handling
func ExecuteTransaction(ctx context.Context, db *sql.DB) error {
	// Start a transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return &DBError{
			Op:  "begin_transaction",
			Err: errors.Wrap(err, "failed to begin transaction"),
		}
	}
	
	// Ensure transaction is rolled back if function returns with error
	defer func() {
		if tx != nil {
			tx.Rollback()
		}
	}()

	// Insert multiple users in the transaction
	users := []struct {
		name  string
		email string
	}{
		{"Alice", "alice@example.com"},
		{"Bob", "bob@example.com"},
		{"Charlie", "charlie@example.com"},
	}

	// Insert each user
	insertSQL := "INSERT INTO users (name, email) VALUES (?, ?)"
	for _, user := range users {
		_, err := tx.ExecContext(ctx, insertSQL, user.name, user.email)
		if err != nil {
			// No need to rollback here, the defer will handle it
			return &DBError{
				Op:  "transaction_insert",
				SQL: insertSQL,
				Err: errors.Wrapf(err, "failed to insert user %s", user.name),
			}
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return &DBError{
			Op:  "commit_transaction",
			Err: errors.Wrap(err, "failed to commit transaction"),
		}
	}
	
	// Set tx to nil to prevent rollback in defer
	tx = nil

	return nil
}

// IsDBError checks if the error is a database error
func IsDBError(err error) bool {
	var dbErr *DBError
	return errors.As(err, &dbErr)
}

// IsNotFoundError checks if the error is a "not found" error
func IsNotFoundError(err error) bool {
	return errors.Is(err, ErrNotFound)
}
