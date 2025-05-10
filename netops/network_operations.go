package netops

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

// FetchWithRetry attempts to fetch data from a URL with retry logic
func FetchWithRetry(ctx context.Context, url string, maxRetries int) ([]byte, error) {
	var lastErr error
	backoff := 100 * time.Millisecond

	for attempt := 0; attempt <= maxRetries; attempt++ {
		// Create a new request with the provided context
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create request")
		}

		// Perform the request
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			lastErr = errors.Wrapf(err, "attempt %d: request failed", attempt)
			
			// Check if the context has been cancelled before retrying
			select {
			case <-ctx.Done():
				return nil, errors.Wrap(ctx.Err(), "context cancelled during network operation")
			default:
				// If this is not the last attempt, wait before retrying
				if attempt < maxRetries {
					// Exponential backoff
					time.Sleep(backoff)
					backoff *= 2 // Double the backoff time for next retry
				}
				continue
			}
		}

		// Always close the response body
		defer resp.Body.Close()

		// Check for non-successful status code
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			lastErr = errors.Errorf("attempt %d: non-successful status code: %d", attempt, resp.StatusCode)
			
			// If this is not the last attempt, wait before retrying
			if attempt < maxRetries {
				time.Sleep(backoff)
				backoff *= 2
				continue
			}
			return nil, lastErr
		}

		// Read the response body
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			lastErr = errors.Wrapf(err, "attempt %d: failed to read response body", attempt)
			
			// If this is not the last attempt, wait before retrying
			if attempt < maxRetries {
				time.Sleep(backoff)
				backoff *= 2
				continue
			}
			return nil, lastErr
		}

		// Success
		return data, nil
	}

	// If we've exhausted all retries, return the last error
	return nil, errors.Wrap(lastErr, "all retry attempts failed")
}

// FetchWithTimeout fetches data from a URL with a specific timeout
func FetchWithTimeout(url string, timeout time.Duration) ([]byte, error) {
	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel() // Ensure resources are cleaned up

	// Create the request with the timeout context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	// Perform the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// Check if the error is due to context timeout
		if ctx.Err() == context.DeadlineExceeded {
			return nil, errors.Wrapf(err, "request to %s timed out after %v", url, timeout)
		}
		return nil, errors.Wrapf(err, "request to %s failed", url)
	}
	defer resp.Body.Close()

	// Check for non-successful status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errors.Errorf("non-successful status code: %d from %s", resp.StatusCode, url)
	}

	// Read the response body
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}

	return data, nil
}

// PostJSON sends a POST request with JSON data and handles errors
func PostJSON(ctx context.Context, url string, jsonData []byte) ([]byte, error) {
	// Create a custom HTTP client with sensible defaults
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Create a request with the provided context
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		// Add contextual information to the error
		return nil, errors.Wrapf(err, "failed to send POST request to %s", url)
	}
	defer resp.Body.Close()

	// Check for non-successful status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Read error response body if available
		errorBody, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			errorBody = []byte("[unable to read error response body]")
		}
		
		return nil, errors.Errorf("non-successful status code: %d, body: %s", 
			resp.StatusCode, string(errorBody))
	}

	// Read the response body
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}

	return responseData, nil
}
