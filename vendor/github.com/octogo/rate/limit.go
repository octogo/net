package rate

import "time"

// Limit is a simple rate-limiter.
// It counts ticks and provides the most recent value of ticks/freq.
type Limit struct {
	count int
	freq  time.Duration

	chTick  chan struct{}
	chClose chan struct{}
}

// NewLimit returns an initialized Limit{} that limits the rate to 1/f ticks.
func NewLimit(f time.Duration) *Limit {
	if f.Nanoseconds() <= int64(0) {
		f = time.Duration(1 * time.Nanosecond)
	}

	limit := &Limit{
		freq:    f,
		chTick:  make(chan struct{}),
		chClose: make(chan struct{}),
	}

	go func() {
		defer func() {
			close(limit.chTick)
			close(limit.chClose)
		}()
		tick := time.Tick(f)
		for {
			select {
			case <-limit.chClose:
				return

			case <-tick:
				limit.chTick <- struct{}{}
			}
		}
	}()

	return limit
}

// Wait blocks until the next free token becomes available.
func (l Limit) Wait() {
	<-l.chTick
}

// Try returns nil if a free token was available and an errTooFast, otherwise.
func (l Limit) Try() error {
	select {
	case <-l.chTick:
		return nil
	default:
		return errTooFast
	}
}

// Close frees the underlying resources.
func (l Limit) Close() {
	l.chClose <- struct{}{}
}
