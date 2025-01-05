package db

import (
	"fmt"
	"log"
	"time"
)

// Retry supplied function with exponential backoff.
// Mainly used for first attempts to connect to database while CockroachDB might be still starting.
func retry[T any](attempts int, duration time.Duration, f func() (T, error)) (result T, err error) {
	for i := 0; i < attempts; i++ {
		if i > 0 {
			log.Println("retrying after error:", err)
			time.Sleep(time.Duration(duration))
			duration *= 2
		}
		result, err = f()
		if err == nil {
			return result, nil
		}
	}
	return result, fmt.Errorf("after %d attempts, last error: %s", attempts, err)
}
