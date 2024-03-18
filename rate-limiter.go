package ratelimit

type RateLimiter interface {
	Accept() bool
	Do() error
	Stop()
}
