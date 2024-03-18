# Rate Limit
This package provides rate-limiting implementation in Golang. Several algorithms are available for rate limiting such as -
* Fixed Window
* Sliding Window
* Leaky Bucket
* Token Bucket

Currently, this package only supports a fixed-window algorithm. In the upcoming releases, more algorithms such as leaky-bucket will be added.

## Fixed Window Algorithm
Create a Fixed Window rate limiter with FixedWindow struct with maximum requests/operations limit for a specific time duration. Now Use Do() method to start a process to reset the counter for each time window.

Note: Do() internally start a goroutine to reset the counter for each time window. Goroutine will be executed at the termination of the program If you manually don't stop it. Use Stop() method to stop the goroutine.

```
package main

import (
	"fmt"

	"github.com/osfbeast/ratelimit"
)

func main() {
	// Initialize rate limiter
	rl := ratelimit.FixedWindow{Duration: 5, Unit: "second", Limit: 2}

	// Validate and apply rate limit configuration
	rl.Do()

	for i := 0; i < 5; i++ {
		if rl.Accept() {
			// Display "ACCEPTED" if the request/operation is allowed by the rate limiter
			fmt.Println("Request", i, "ACCEPTED")
			continue
		}

		// Display "REJECTED" if the request/operation is blocked by the rate limiter
		fmt.Println("Request", i, "REJECTED")
	}
}

```