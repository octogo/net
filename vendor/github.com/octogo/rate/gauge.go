package rate

import (
	"time"

	"github.com/octogo/rate/meter"
)

// Gauge implements a rate-limiter in combination with a meter for measuring rate.
type Gauge struct {
	meter *meter.Speedo
	limit *Limit
}

// NewGauge returns an initialized Gauge{}.
func NewGauge(b, i int, f time.Duration) *Gauge {
	gauge := &Gauge{
		meter: meter.New(f),
		limit: NewBurstLimit(b, i, f),
	}

	return gauge
}

// Wait blocks until the next free token is available.
// This method also adds a single tick before returning.
func (g Gauge) Wait() {
	defer g.meter.Tick()
	g.limit.Wait()
}

// Try returns nil if a free token is available
func (g Gauge) Try() error {
	return g.limit.Try()
}

// Tick ...
func (g Gauge) Tick() {
	g.meter.Tick()
}

// Measure ...
func (g Gauge) Measure() float64 {
	return g.meter.Measure()
}

// Close ...
func (g Gauge) Close() {
	g.limit.Close()
	g.meter.Close()
}
