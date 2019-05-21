package rate

// Limiter defines the interface of a rate-limiter.
type Limiter interface {
	// Close frees the underlying resources.
	Close()

	// Try returns nil if a token was available and an error when there wasn't.
	Try() error

	// Wait blocks until the next free token is available.
	Wait()
}
