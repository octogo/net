package meter

import (
	"time"
)

// Speedo is a simple meter for counting hits per unit.
type Speedo struct {
	freq     time.Duration
	counter  float64
	lastRate float64
	lastTick time.Time

	chTick    chan struct{}
	chMeasure chan float64
	chClose   chan struct{}
}

// New returns an initialized meter.Speedo{}.
func New(freq time.Duration) *Speedo {
	meter := &Speedo{
		freq:    freq,
		counter: 0,

		chTick:    make(chan struct{}),
		chMeasure: make(chan float64),
		chClose:   make(chan struct{}),
	}

	go func() {
		defer func() {
			close(meter.chTick)
			close(meter.chMeasure)
			close(meter.chClose)
		}()
		nextTick := time.After(meter.freq)
		for {
			select {
			case <-meter.chClose:
				return
			case <-nextTick:
				meter.lastRate = meter.counter
				meter.counter = 0
				nextTick = time.After(meter.freq)
			case meter.chMeasure <- meter.lastRate:
			case <-meter.chTick:
				meter.counter++
			}
		}
	}()

	return meter
}

// Tick adds a tick to this speedometer's counter.
func (r Speedo) Tick() {
	r.chTick <- struct{}{}
}

// Measure returns the recent meter of ticks per freq.
func (r Speedo) Measure() float64 {
	return <-r.chMeasure
}

// Close frees the underlying resources.
func (r Speedo) Close() {
	r.chClose <- struct{}{}
}
