# clock

[![godoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/brandoneprice31/clock)
[![License](https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/brandoneprice31/clock)](https://goreportcard.com/report/github.com/brandoneprice31/clock)

Go package useful for dealing with time (hours, minutes, and seconds).  The native `time.Time` struct built-in to Go also deals with dates and timezones.  This package is for the more narrow and simpler use case for when all you care about is clock time.

## Quick Start

```go
package main

import (
    "fmt"

    "github.com/brandoneprice31/clock"
)

func main() {
    noon := clock.NewTime(12, 0, 0)

    now := clock.Now("America/Los_Angeles")

    if now.After(noon) {
        fmt.Println("It is afternoon!")
    }
}
```

## Installation / Usage

To install `clock`, use `go get`:
```
go get github.com/brandoneprice31/clock
```

Import the `brandoneprice31/clock` package into your code:
```go
import "github.com/brandoneprice31/clock"
```

## Staying Up to Date

To update `clock` to the latest version, use `go get -u github.com/brandoneprice31/clock`.

## Contributing

Please feel free to submit issues, fork the repository and send pull requests!
