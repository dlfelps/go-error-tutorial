package file

import (
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
)

// WriteToFile writes data to a file with proper error handling
func WriteToFile(filename, data string) error {
	// Open file with proper flags and permissions
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		// Wrap the error with context
		return errors.Wrap(err, fmt.Sprintf("failed to open file for writing: %s", filename))
	}
	// Ensure the file is closed when function completes
	defer func() {
		// Close the file, but don't overwrite the original error if there was one
		cerr := file.Close()
		if err == nil && cerr != nil {
			err = errors.Wrap(cerr, fmt.Sprintf("failed to close file: %s", filename))
		}
	}()

	// Write data to file
	_, err = file.WriteString(data)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to write to file: %s", filename))
	}
	
	// Ensure data is written to disk
	err = file.Sync()
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to sync file: %s", filename))
	}

	return nil
}

// ReadFromFile reads data from a file with proper error handling
func ReadFromFile(filename string) (string, error) {
	// Check if file exists before attempting to read
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return "", errors.Wrap(err, fmt.Sprintf("file does not exist: %s", filename))
	}

	// Read the entire file
	data, err := os.ReadFile(filename)
	if err != nil {
		// Different error handling based on error type
		if os.IsPermission(err) {
			return "", errors.Wrap(err, fmt.Sprintf("permission denied for file: %s", filename))
		}
		return "", errors.Wrap(err, fmt.Sprintf("failed to read file: %s", filename))
	}

	return string(data), nil
}

// SafeCopyFile safely copies a file with proper error handling
func SafeCopyFile(src, dst string) error {
	// Check if source file exists
	sourceInfo, err := os.Stat(src)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.Wrap(err, fmt.Sprintf("source file does not exist: %s", src))
		}
		return errors.Wrap(err, fmt.Sprintf("failed to get source file info: %s", src))
	}

	// Ensure source is a regular file
	if !sourceInfo.Mode().IsRegular() {
		return errors.New(fmt.Sprintf("source is not a regular file: %s", src))
	}

	// Open source file
	sourceFile, err := os.Open(src)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to open source file: %s", src))
	}
	// Ensure source file is closed when function completes
	defer sourceFile.Close()

	// Create destination file
	// We use a temporary file and then rename to ensure atomicity
	tempDst := dst + ".tmp"
	destFile, err := os.Create(tempDst)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to create destination file: %s", dst))
	}

	// Setup deferred cleanup in case of failure
	success := false
	defer func() {
		// Close the file
		destFile.Close()
		
		// If the operation was not successful, remove the temporary file
		if !success {
			os.Remove(tempDst)
		}
	}()

	// Copy the content
	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return errors.Wrap(err, "failed to copy file content")
	}

	// Ensure data is written to disk
	err = destFile.Sync()
	if err != nil {
		return errors.Wrap(err, "failed to sync destination file")
	}

	// Close the file before renaming
	err = destFile.Close()
	if err != nil {
		return errors.Wrap(err, "failed to close destination file")
	}

	// Rename the temporary file to the actual destination
	err = os.Rename(tempDst, dst)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to rename temporary file to destination: %s", dst))
	}

	// Mark as successful to prevent cleanup of the temporary file
	success = true
	return nil
}

// AppendToFile appends data to a file with proper error handling
func AppendToFile(filename, data string) error {
	// Open file with append flag
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to open file for appending: %s", filename))
	}
	defer file.Close()

	// Write data to file
	_, err = file.WriteString(data)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to append to file: %s", filename))
	}

	return nil
}

// DeleteFile safely deletes a file with proper error handling
func DeleteFile(filename string) error {
	// Check if file exists before attempting to delete
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// Not an error if file doesn't exist - it's already deleted
		return nil
	}

	// Delete the file
	err := os.Remove(filename)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to delete file: %s", filename))
	}

	return nil
}
