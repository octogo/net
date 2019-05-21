package service

import "errors"

var (
	errTooManyConnections      = errors.New("too many connections")
	errTooManyConnectionsPerIP = errors.New("too many connections from IP address")
)
