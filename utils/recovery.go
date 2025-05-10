package utils

import (
	"fmt"
	"runtime/debug"

	"github.com/sirupsen/logrus"
)

// SafeGo runs a function in a goroutine with panic recovery
func SafeGo(log *logrus.Logger, fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				stack := debug.Stack()
				log.WithFields(logrus.Fields{
					"panic": r,
					"stack": string(stack),
				}).Error("Recovered from panic in goroutine")
			}
		}()
		fn()
	}()
}

// RecoverMiddleware is a middleware function that recovers from panics
// It's useful for HTTP handlers or any function that needs panic recovery
func RecoverMiddleware(log *logrus.Logger, next func()) {
	defer func() {
		if r := recover(); r != nil {
			stack := debug.Stack()
			log.WithFields(logrus.Fields{
				"panic": r,
				"stack": string(stack),
			}).Error("Recovered from panic")
		}
	}()
	next()
}

// RecoverWithCallback recovers from panics and calls a callback function
// This is useful when you need to do custom handling after a panic
func RecoverWithCallback(callback func(interface{}, []byte)) {
	if r := recover(); r != nil {
		stack := debug.Stack()
		callback(r, stack)
	}
}

// SafeExecute executes a function with panic recovery and returns an error if a panic occurs
func SafeExecute(fn func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			stack := debug.Stack()
			err = fmt.Errorf("panic recovered: %v\nstack: %s", r, stack)
		}
	}()
	return fn()
}
