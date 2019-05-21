[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://img.shields.io/badge/license-MIT-blue.svg)
[![Build Status](https://travis-ci.org/octogo/rate.svg?branch=master)](https://travis-ci.org/octogo/rate)
[![GoDoc](https://godoc.org/github.com/octogo/rate?status.svg)](https://godoc.org/github.com/octogo/rate)

# OctoRate

Package `rate` implements various utilities for rate-limiting and monitoring
rate-limit utilization.

## Getting Started

### Installation

```bash
go get github.com/octogo/rate
```

### Usage

#### Primitive Rate-Limiter

```go
import (
    "fmt"
    "time"

    "github.com/octogo/rate"
)

func main() {
    rateLimit := rate.NewLimit(250 * time.Millisecond)
    defer rateLimit.Close()

    for i := 0; i < 10; i++ {
        rateLimit.Wait()
        fmt.Println(i)
    }
}
```

#### Burst Bucket Rate-Limiter

```go
import (
    "fmt"
    "time"

    "github.com/octogo/rate"
)

func main() {
    burstLimit := rate.NewBurstLimit(
        10,             // burst bucket size
        4,              // regen this many tokens
        time.Second,    // at this rate
    )
    defer rateLimit.Close()

    // burst through the bucket
    for i := 0; i < 10; i++ {
        rateLimit.Wait()
        fmt.Println(i)
    }

    // actually wait a tick
    rateLimit.Wait()
}
```

#### Gauge

```go
import (
    "fmt"
    "time"

    "github.com/octogo/rate"
)

func main() {
    gauge := rate.NewGauge(
        10,             // burst bucket size
        4,              // regen this many tokens
        time.Second,    // at this rate
    )
    defer rateLimit.Close()

    for i := 0; i < 10; i++ {
        gauge.Wait()
        fmt.Println(i)
    }

    // after one second, the Gauge should be showing 10/s ticks.
    <-time.After(time.Second)
    fmt.Println(gauge.Measure())
}
```
