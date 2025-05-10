package utils

import (
	"context"
	"math/rand"
	"time"
)

// RetryOptions configures the retry behavior
type RetryOptions struct {
	MaxRetries    int           // Maximum number of retry attempts
	BaseDelay     time.Duration // Base delay between retries
	MaxDelay      time.Duration // Maximum delay between retries
	Factor        float64       // Factor to increase the delay with each retry
	Jitter        float64       // Randomness factor to add to the delay (0.0-1.0)
	RetryableFunc func(error) bool // Function to determine if an error is retryable
}

// DefaultRetryOptions provides sensible default retry options
func DefaultRetryOptions() RetryOptions {
	return RetryOptions{
		MaxRetries: 3,
		BaseDelay:  100 * time.Millisecond,
		MaxDelay:   10 * time.Second,
		Factor:     1.5,
		Jitter:     0.2,
		RetryableFunc: func(err error) bool {
			// By default, retry all errors
			return err != nil
		},
	}
}

// Retry executes the given function with exponential backoff retry logic
func Retry(ctx context.Context, fn func() error, opts RetryOptions) error {
	var err error
	
	// Initialize random number generator for jitter
	rand.Seed(time.Now().UnixNano())

	// Keep track of the current delay
	currentDelay := opts.BaseDelay

	// Try the operation up to MaxRetries times
	for attempt := 0; attempt <= opts.MaxRetries; attempt++ {
		// Execute the function
		err = fn()
		
		// If there was no error or the error is not retryable, return immediately
		if err == nil || (opts.RetryableFunc != nil && !opts.RetryableFunc(err)) {
			return err
		}

		// If this was the last attempt, return the error
		if attempt == opts.MaxRetries {
			return err
		}

		// Calculate the next delay with exponential backoff and jitter
		jitter := 1.0
		if opts.Jitter > 0 {
			jitter = 1.0 + (rand.Float64()*2-1)*opts.Jitter // Random value between (1-jitter) and (1+jitter)
		}
		
		nextDelay := time.Duration(float64(currentDelay) * opts.Factor * jitter)
		
		// Cap the delay at MaxDelay
		if nextDelay > opts.MaxDelay {
			nextDelay = opts.MaxDelay
		}
		
		// Update the current delay for the next iteration
		currentDelay = nextDelay

		// Wait for the delay or until the context is cancelled
		select {
		case <-ctx.Done():
			// Context was cancelled, return the context error
			return ctx.Err()
		case <-time.After(nextDelay):
			// Continue to the next attempt
		}
	}

	// This should never be reached due to the return in the loop
	return err
}

// RetryWithResult is like Retry but for functions that return a value and an error
func RetryWithResult[T any](ctx context.Context, fn func() (T, error), opts RetryOptions) (T, error) {
	var result T
	var err error
	
	// Initialize random number generator for jitter
	rand.Seed(time.Now().UnixNano())

	// Keep track of the current delay
	currentDelay := opts.BaseDelay

	// Try the operation up to MaxRetries times
	for attempt := 0; attempt <= opts.MaxRetries; attempt++ {
		// Execute the function
		result, err = fn()
		
		// If there was no error or the error is not retryable, return immediately
		if err == nil || (opts.RetryableFunc != nil && !opts.RetryableFunc(err)) {
			return result, err
		}

		// If this was the last attempt, return the error
		if attempt == opts.MaxRetries {
			return result, err
		}

		// Calculate the next delay with exponential backoff and jitter
		jitter := 1.0
		if opts.Jitter > 0 {
			jitter = 1.0 + (rand.Float64()*2-1)*opts.Jitter // Random value between (1-jitter) and (1+jitter)
		}
		
		nextDelay := time.Duration(float64(currentDelay) * opts.Factor * jitter)
		
		// Cap the delay at MaxDelay
		if nextDelay > opts.MaxDelay {
			nextDelay = opts.MaxDelay
		}
		
		// Update the current delay for the next iteration
		currentDelay = nextDelay

		// Wait for the delay or until the context is cancelled
		select {
		case <-ctx.Done():
			// Context was cancelled, return the context error
			return result, ctx.Err()
		case <-time.After(nextDelay):
			// Continue to the next attempt
		}
	}

	// This should never be reached due to the return in the loop
	return result, err
}
