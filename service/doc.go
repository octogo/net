/*
Package service implements a simple service manager for managing a net.Listener.

Features:

- limit the maximum number of allowed connections
- limit the maximum number of allowed connections per IP
- limit the rate at which new connections are handled (with support for bursts)
- monitoring utilization

*/
package service
