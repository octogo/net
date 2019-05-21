/*
Package rate implements various rate-limiting utilities.

- Limit

  A primitive rate-limiter, regenerating a new token at the given rate.

- BurstLimit

  A rate-limiter with burst token bucket support.
  Once all tokens are consumed, new ones are generated at the given rate.

- Gauge

  A wrapper for a limiter and a rate-limiter and a meter.
  A Gauge can limit rate and has support for monitoring utilization.
  (See: github.com/octogo/rate/meter)
*/
package rate
