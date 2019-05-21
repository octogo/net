package meter

// Meter defines the interface of a meter.
//
// With a meter you can measure arbitrary rate.
type Meter interface {
	Tick()
	Measure() float64
	Close()
}
