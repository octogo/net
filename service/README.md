# OctoNetService

A Golang package implementing a simple service-manager.

Services takes a net.Listener and a handler-function and then manages
the Listener in terms of connection- and rate-limiting.

Supported limits are:

- new connection rate *(with support for bursts)*
- max connections
- max connections per IP

## Getting Started

### Installation

```bash
go get github.com/octogo/net/service
```

### Usage

```go
import (
    "log"

    "github.com/octogo/net/foo"
)

func handler(c net.Conn) error {
    c.Write([]byte("hello world!"))
    return nil
}

func main() {
    l, err := net.Listen("tcp", "")
    if err != nil {
        log.Fatal(err)
    }

    // initialize new service
    s := service.New(l, handler)

    // configure it for max 10 connections/s with a burst size of 10
    s.ConnectBurst = 10
    s.ConnectCount = 10
    s.ConnectRate = time.Second

    // and run it
    s.Run()
}
```