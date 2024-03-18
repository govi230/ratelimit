package ratelimit

import (
	"fmt"
	"sync"
	"time"
)

type FixedWindow struct {
	// Duration specifies the length of the time window. It must be greater then zero.
	Duration uint64
	// Unit specifies the time unit for the duration (e.g., "second", "minute", "hour").
	// Supported Units-  "second", "minute", "hour"
	Unit string
	// Limit is the maximum number of requests allowed during the time window.
	// It must be greater then zero.
	Limit uint64

	// counter keeps track of the number of requests accepted within the current window.
	counter uint64
	// ticker is used to reset the counter after each time window.
	ticker *time.Ticker
	// stop is used to prevent resetting the counter for further time windows.
	stop bool
	// mu is a mutex to prevent data race conditions in concurrent goroutines.
	mu *sync.RWMutex
}

// Accept checks whether a request will be accepted or not.
// It verifies if the current number of requests within the time window has reached its limit.
// If the limit is reached, it returns false. Otherwise it will return true and also increase the counter with one.
func (fw *FixedWindow) Accept() (accepted bool) {
	// Lock with mutex
	fw.mu.Lock()
	defer fw.mu.Unlock()

	if accepted = !fw.stop && fw.counter != fw.Limit; accepted {
		fw.counter++
	}

	return
}

// Counter returns the total number of accepted requests within the current time window.
func (fw *FixedWindow) Counter() uint64 {
	return fw.counter
}

// Validate checks the rate limiter configuration for validity.
// It ensures that the duration, limit, and time unit are properly configured.
// Returns an error if any of the configurations are invalid.
func (fw *FixedWindow) Validate() error {
	if fw.Duration == 0 {
		return fmt.Errorf("duration must be greater than zero")
	}

	if fw.Limit == 0 {
		return fmt.Errorf("limit must be greater than zero")
	}

	if fw.Unit != "second" && fw.Unit != "minute" && fw.Unit != "hour" {
		return fmt.Errorf("expected one of them: 'second', 'minute', 'hour' but got '%s'", fw.Unit)
	}

	return nil
}

// Do validates the rate limiter configuration and initializes the internal fields.
// It starts a goroutine to reset the counter to zero for each time window.
// If an error occurs during validation or initialization of rate limiter fields, it returns an error object.
// Otherwise, it returns nil.
// NOTE:
// As discussed, it starts a goroutine to reset the counter. To stop resetting the counter, you need to use the Stop method.
// If you lose the reference of FixedWindow typed variables without stopping the resetter,
// it will continue running until the program terminates, leading to unnecessary CPU/RAM usage.
// For example, if a function is defined as follows and the "rl" object is not returned:
//
//	func limiter() {
//		rl := FixedWindow(Duration: 10, Unit: "second", Limit: 100)
//		if err := rl.Do(); err != nil {
//			panic(err)
//		}
//	}
//
// It will start a goroutine for resetting the counter, but you lose the reference of rl variable.
// Therefore, ensure to retain a reference to the FixedWindow object if you need to stop resetting the counter later.
func (fw *FixedWindow) Do() (err error) {

	// Validate rate limiter configuration
	if err = fw.Validate(); err != nil {
		return err
	}

	// Set mutex for prevent data race condition
	fw.mu = &sync.RWMutex{}

	// Counter reset ticker
	fw.ticker = time.NewTicker(fw.duration())

	// Start a goroutine to reset the counter
	go fw.reset()

	return
}

// duration returns a time.Duration object for the provided rate limiter configuration.
func (fw *FixedWindow) duration() time.Duration {
	switch fw.Unit {
	default:
		panic(fmt.Sprintf("Unsupported time unit '%s'", fw.Unit))
	case "second":
		return time.Duration(fw.Duration) * time.Second
	case "minute":
		return time.Duration(fw.Duration) * time.Minute
	case "hour":
		return time.Duration(fw.Duration) * time.Hour
	}
}

// It helps to stop to reset counter for further time windows.
func (fw *FixedWindow) Stop() {
	fw.stop = true
	fw.ticker.Stop()
}

// It resets the counter to zero for each time window.
func (fw *FixedWindow) reset() {
	for !fw.stop {
		<-fw.ticker.C
		fw.counter = 0
	}
}
