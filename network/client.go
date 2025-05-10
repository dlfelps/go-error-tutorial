package network

import (
        "context"
        "fmt"
        "io"
        "net/http"
        "os"
        "time"

        "github.com/pkg/errors"
        "github.com/sirupsen/logrus"
)

// NetworkError represents an error that occurred during a network operation
type NetworkError struct {
        URL     string
        Op      string
        Err     error
        Retries int
}

// Error implements the error interface
func (e *NetworkError) Error() string {
        return fmt.Sprintf("network error during %s on %s (after %d retries): %v", e.Op, e.URL, e.Retries, e.Err)
}

// Unwrap returns the underlying error
func (e *NetworkError) Unwrap() error {
        return e.Err
}

// FetchURL fetches a URL with retries and timeout
func FetchURL(ctx context.Context, url string, maxRetries int) (*http.Response, error) {
        var lastErr error
        
        // Create a custom HTTP client with sensible defaults
        client := &http.Client{
                Timeout: 10 * time.Second, // Default timeout
                Transport: &http.Transport{
                        MaxIdleConns:        10,
                        IdleConnTimeout:     30 * time.Second,
                        DisableCompression:  false,
                        TLSHandshakeTimeout: 5 * time.Second,
                },
        }

        // Initialize logger for this function
        log := logrus.New()
        log.SetFormatter(&logrus.JSONFormatter{})

        // Try the request with retries
        for retry := 0; retry <= maxRetries; retry++ {
                // Check if context is cancelled before making the request
                if ctx.Err() != nil {
                        return nil, &NetworkError{
                                URL:     url,
                                Op:      "fetch",
                                Err:     ctx.Err(),
                                Retries: retry,
                        }
                }

                // Log retry attempt
                if retry > 0 {
                        log.WithFields(logrus.Fields{
                                "url":   url,
                                "retry": retry,
                                "max":   maxRetries,
                        }).Info("Retrying request")

                        // Add exponential backoff before retrying
                        backoffTime := time.Duration(1<<uint(retry-1)) * 100 * time.Millisecond
                        
                        // Create a timer that will be cancelled if context is cancelled
                        timer := time.NewTimer(backoffTime)
                        select {
                        case <-ctx.Done():
                                timer.Stop()
                                return nil, &NetworkError{
                                        URL:     url,
                                        Op:      "fetch_backoff",
                                        Err:     ctx.Err(),
                                        Retries: retry,
                                }
                        case <-timer.C:
                                // Continue with the retry
                        }
                }

                // Create a new request with the provided context
                req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
                if err != nil {
                        lastErr = err
                        continue
                }

                // Set common headers
                req.Header.Set("User-Agent", "GoErrorHandlingDemo/1.0")
                
                // Execute the request
                resp, err := client.Do(req)
                if err != nil {
                        lastErr = err
                        log.WithError(err).WithField("url", url).Error("Request failed")
                        continue
                }

                // Check for successful status code
                if resp.StatusCode >= 400 {
                        // Read response body for error details
                        body, readErr := io.ReadAll(resp.Body)
                        resp.Body.Close()
                        
                        if readErr != nil {
                                log.WithError(readErr).Error("Failed to read error response body")
                                // Continue with original error
                        }
                        
                        lastErr = fmt.Errorf("bad status code: %d, body: %s", resp.StatusCode, string(body))
                        log.WithFields(logrus.Fields{
                                "status_code": resp.StatusCode,
                                "url":         url,
                        }).Error("Request returned error status")
                        continue
                }

                // Success!
                return resp, nil
        }

        // If we got here, all retries failed
        return nil, &NetworkError{
                URL:     url,
                Op:      "fetch",
                Err:     errors.Wrap(lastErr, "all retries failed"),
                Retries: maxRetries,
        }
}

// PostJSON sends a JSON payload to a URL with retries
func PostJSON(ctx context.Context, url string, payload []byte, maxRetries int) (*http.Response, error) {
        var lastErr error
        
        // Create a custom HTTP client with sensible defaults
        client := &http.Client{
                Timeout: 10 * time.Second,
        }

        // Try the request with retries
        for retry := 0; retry <= maxRetries; retry++ {
                // Check if context is cancelled before making the request
                if ctx.Err() != nil {
                        return nil, &NetworkError{
                                URL:     url,
                                Op:      "post_json",
                                Err:     ctx.Err(),
                                Retries: retry,
                        }
                }

                // Add exponential backoff before retrying
                if retry > 0 {
                        backoffTime := time.Duration(1<<uint(retry-1)) * 100 * time.Millisecond
                        time.Sleep(backoffTime)
                }

                // Create a context for this specific request
                reqCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
                defer cancel()

                // Create a new POST request with the payload
                req, err := http.NewRequestWithContext(reqCtx, "POST", url, nil)
                if err != nil {
                        lastErr = err
                        continue
                }

                // Set content type to JSON
                req.Header.Set("Content-Type", "application/json")
                
                // Execute the request
                resp, err := client.Do(req)
                if err != nil {
                        lastErr = err
                        continue
                }

                // Check for successful status code
                if resp.StatusCode >= 400 {
                        resp.Body.Close()
                        lastErr = fmt.Errorf("bad status code: %d", resp.StatusCode)
                        continue
                }

                // Success!
                return resp, nil
        }

        // If we got here, all retries failed
        return nil, &NetworkError{
                URL:     url,
                Op:      "post_json",
                Err:     errors.Wrap(lastErr, "all retries failed"),
                Retries: maxRetries,
        }
}

// DownloadFile downloads a file with proper error handling
func DownloadFile(ctx context.Context, url string, destPath string) error {
        // Get the data with retry
        resp, err := FetchURL(ctx, url, 3)
        if err != nil {
                return errors.Wrap(err, "failed to download file")
        }
        defer resp.Body.Close()

        // Create the file
        out, err := createFileWithErrorHandling(destPath)
        if err != nil {
                return err
        }
        defer out.Close()

        // Copy the response body to the file
        _, err = io.Copy(out, resp.Body)
        if err != nil {
                // If copy fails, try to remove the partial file
                out.Close()
                os.Remove(destPath)
                return errors.Wrap(err, "failed to write downloaded content to file")
        }

        return nil
}

// createFileWithErrorHandling creates a file with proper error handling
func createFileWithErrorHandling(filePath string) (*os.File, error) {
        file, err := os.Create(filePath)
        if err != nil {
                if os.IsPermission(err) {
                        return nil, errors.Wrap(err, "permission denied when creating file")
                }
                if os.IsExist(err) {
                        return nil, errors.Wrap(err, "file already exists")
                }
                return nil, errors.Wrap(err, "failed to create file")
        }
        return file, nil
}
