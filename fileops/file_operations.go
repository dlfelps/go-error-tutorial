package fileops

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/pkg/errors"
)

// WriteFile writes data to a file with proper error handling
func WriteFile(filename string, content string) error {
	// Create the file with appropriate permissions
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return errors.Wrap(err, "failed to create file")
	}
	// Use defer to ensure the file is closed properly when the function returns
	defer func() {
		closeErr := file.Close()
		if closeErr != nil {
			// Log the error but don't override the original error if there was one
			// In a real application, you might want to use a logger here
			println("Error closing file:", closeErr.Error())
		}
	}()

	// Write the content to the file
	_, err = file.WriteString(content)
	if err != nil {
		return errors.Wrap(err, "failed to write to file")
	}

	// Explicitly sync to ensure data is written to disk
	err = file.Sync()
	if err != nil {
		return errors.Wrap(err, "failed to sync file")
	}

	return nil
}

// ReadFileWithContext reads a file with a context for cancellation
func ReadFileWithContext(ctx context.Context, filename string) (string, error) {
	// Create a channel to communicate the result
	resultCh := make(chan struct {
		content string
		err     error
	})

	// Start a goroutine to read the file
	go func() {
		// Open the file
		file, err := os.Open(filename)
		if err != nil {
			resultCh <- struct {
				content string
				err     error
			}{"", errors.Wrap(err, "failed to open file")}
			return
		}
		defer file.Close()

		// Read the file content
		content, err := ioutil.ReadAll(file)
		if err != nil {
			resultCh <- struct {
				content string
				err     error
			}{"", errors.Wrap(err, "failed to read file")}
			return
		}

		// Simulate a long-running operation
		time.Sleep(100 * time.Millisecond)

		// Send the result
		resultCh <- struct {
			content string
			err     error
		}{string(content), nil}
	}()

	// Wait for either the context to be cancelled or the read to complete
	select {
	case <-ctx.Done():
		return "", errors.Wrap(ctx.Err(), "context cancelled while reading file")
	case result := <-resultCh:
		return result.content, result.err
	}
}

// CopyFileWithProgress copies a file with progress tracking and proper error handling
func CopyFileWithProgress(src, dst string, progressFn func(bytesRead int64, total int64)) error {
	// Open source file
	sourceFile, err := os.Open(src)
	if err != nil {
		return errors.Wrap(err, "failed to open source file")
	}
	defer sourceFile.Close()

	// Get file size for progress reporting
	fileInfo, err := sourceFile.Stat()
	if err != nil {
		return errors.Wrap(err, "failed to get source file info")
	}
	totalSize := fileInfo.Size()

	// Create destination file
	destFile, err := os.Create(dst)
	if err != nil {
		return errors.Wrap(err, "failed to create destination file")
	}
	// Use defer and anonymous function to handle errors from Close
	defer func() {
		if err := destFile.Close(); err != nil {
			println("Error closing destination file:", err.Error())
		}
	}()

	// Create buffer to optimize copy
	buffer := make([]byte, 32*1024) // 32KB buffer
	var bytesRead int64 = 0

	// Copy the file
	for {
		n, err := sourceFile.Read(buffer)
		if err != nil && err != io.EOF {
			return errors.Wrap(err, "error reading from source file")
		}
		if n == 0 {
			break
		}

		// Write to destination file
		if _, err := destFile.Write(buffer[:n]); err != nil {
			return errors.Wrap(err, "error writing to destination file")
		}

		// Update progress
		bytesRead += int64(n)
		if progressFn != nil {
			progressFn(bytesRead, totalSize)
		}
	}

	// Sync to ensure all data is written to disk
	if err := destFile.Sync(); err != nil {
		return errors.Wrap(err, "failed to sync destination file")
	}

	return nil
}
