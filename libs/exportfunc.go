package main

import (
	"C"

	"time"

	"golang.org/x/time/rate"
)

var lm *Limiter = NewLimiter()

//export Allow
func Allow(resource string, interval int64, burst int64) bool {
	var (
		now         time.Time     = time.Now()
		rateLimiter *rate.Limiter = lm.GetRateLimiter(resource, now, interval, burst)
	)
	return rateLimiter.AllowN(now, 1)
}

func main() {
}
