[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://img.shields.io/badge/license-MIT-blue.svg)
[![Build Status](https://travis-ci.org/octogo/dict.svg?branch=master)](https://travis-ci.org/octogo/dict)
[![GoDoc](https://godoc.org/github.com/octogo/dict?status.svg)](https://godoc.org/github.com/octogo/dict)

# OctoDict

Package `dict` implements a goroutine-safe key-value-store in form of a
`map[string]interface{}`.

## Getting Started

### Installation

```bash
go get github.com/octogo/dict
```

### Usage

```go
package main

import (
    "log"

    "github.com/octogo/dict"
)

func main() {
    d := dict.New(nil)
    defer d.Close()

    // set some values
    d.Set("foo", 42)
    d.Set("hello", "world")

    // go-style get with second return value
    v, ok := d.Get("foo")
    if !ok {
        log.Fatal("not found")
    }
    log.Println(v)

    // primitive get with fallback value as default
    v = d.GetDefault("hello", nil)
    if v == nil {
        log.Fatal("got: nil, want: \"world\"")
    }
    log.Println(v)
}
```
