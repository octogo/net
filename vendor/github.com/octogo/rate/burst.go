package rate

import "time"

// NewBurstLimit is a rate-limiter with support for bursts.
// It behaves just like a primitive rate-limiter, except that it allows for
// bursts of size b. Once all burst tokens are consumed, the bucket is
// refilled at the rate i / f ticks.
func NewBurstLimit(b, i int, f time.Duration) *Limit {
	if b < 0 {
		b = 0
	}
	if i <= 0 {
		i = 0
	}
	if f.Nanoseconds() <= int64(0) {
		f = time.Nanosecond
	}

	burst := &Limit{
		count: i,
		freq:  f,

		chTick:  make(chan struct{}, b),
		chClose: make(chan struct{}),
	}

	// prefill the bucket
	for j := 0; j < b; j++ {
		burst.chTick <- struct{}{}
	}

	go func(burst *Limit) {
		defer func() {
			close(burst.chTick)
			close(burst.chClose)
		}()
		for {
			select {
			case <-time.After(
				time.Duration(burst.freq.Nanoseconds() / int64(burst.count)),
			):
				select {
				case burst.chTick <- struct{}{}:
				default: // noop
				}
			case <-burst.chClose:
				return
			}
		}
	}(burst)

	return burst
}
