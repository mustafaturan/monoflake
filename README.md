# Monoflake

[![Build Status](https://travis-ci.org/mustafaturan/monoflake.svg?branch=master)](https://travis-ci.org/mustafaturan/monoflake)
[![Go Report Card](https://goreportcard.com/badge/github.com/mustafaturan/monoflake)](https://goreportcard.com/report/github.com/mustafaturan/monoflake)
[![GoDoc](https://godoc.org/github.com/mustafaturan/monoflake?status.svg)](https://godoc.org/github.com/mustafaturan/monoflake)

Highly scalable, single/multi node, predictable and incremental 64 bits (8 bytes) unique id
generator with zero allocation magic. It is using snowflake idea with flexible options.

If you are interested in 16 bytes ids with optional sequencers, checkout the
[monoton](http://github.com/mustafaturan/monoton) package.

## Installation

Via go packages:
```go get github.com/mustafaturan/monoflake```

## API

The method names and arities/args are stable now. No change should be expected
on the package for the version `1.x.x` except any bug fixes.

## Usage

### Using with Singleton

Create a new package like below, and then call `Next()` method:

```go
package uniqid

// Import packages
import (
	"github.com/mustafaturan/monoflake"
)

var m *monoflake.MonoFlake

// On init configure the monoflake
func init() {
	m = newIDGenerator()
}

func newIDGenerator() *monoflake.MonoFlake {
	// Fetch your node id from a config server or generate from MAC/IP address
	node := uint16(0)

	// If we want to init the time with 2024-06-01 00:00:00 UTC (min allowed)
	epoch := time.Unix(1717200000, 0)

	// Configure monoflake with a node and epoch
	m, err = monoflake.New(node, monoflake.WithEpoch(epoch))
	if err != nil{
		panic(err)
	}

	return m
}

func Generate() int64 {
	return m.Next().Int64()
}

func GeneateBytes() []byte {
	return m.Next().Bytes()
}

func GeneateString() string {
	return m.Next().String()
}
```


In any other package generate the ids like below:

```go
import (
	"fmt"
	"uniqid" // your local uniqid package from your project
)

func main() {
	for i := 0; i < 100; i++ {
		fmt.Println(uniqid.Generate())
	}
}
```

### Using with Dependency Injection

```go
package main

// Import packages
import (
	"fmt"
	"github.com/mustafaturan/monoflake"
)

func NewIDGenerator() *monoflake.MonoFlake {
	// Fetch your node id from a config server or generate from MAC/IP address
	node := uint16(0)

	// Configure monoflake with a sequencer and the node
	m, err := monoflake.New(node)
	if err != nil{
		panic(err)
	}

	return m
}

func main() {
	g := NewIDGenerator()

	for i := 0; i < 100; i++ {
		fmt.Println(g.Next())
	}
}
```

### Initilization options

#### monoflake.WithNodeBits

```
# monoflake.WithNodeBits sets the max node bits to 8 which modules the node id with 256 limits the val to [0, 256).
# maximum node bits allowed is 2^13, you can go up to 13 by setting it
# You can set node id between 8 to 13 bits inclusive, the rest will be automatically used for sequencer.
monoflake.New(node, monoflake.WithNodeBits(13), ...)
```

#### monoflake.WithEpoch

```
# monoflake.WithEpoch sets the epoch start time, minimum epoch value is 2024-06-01 00:00:00 UTC.

epoch := time.Unix(...)

monoflake.New(node, monoflake.WithEpoch(epoch), ...)
```

## How does it work?

**Default bit allocations**

```
[1 bit(reserved) | 40 bits (time in milliseconds) | 13 bits (sequencer) | 10 bits (node id)]
[ [0, 1)         | up to 34 years in milliseconds | [0, 8192)           | [0, 1024)        ]
```

### Max sequencer case

When the sequencer reaches to maximum value in the same milliseconds then milliseconds increased automatically and the
sequence set to 0.

### Thread safety

Thread safety achieved with mutex locks.

## Features

### Time ordered

The `monoflake` package provides sequences based on the `monotonic` time which
represents the absolute elapsed wall-clock time since some arbitrary, fixed
point in the past. It isn't affected by changes in the system time-of-day clock.

### Epoch time

Epoch time value opens space for time value by subtracting the given value from the time sequence.

### Readable

It comes with `String()` method which encodes the ids into base62 as string and allows padded with zeros to 11 bytes.

### Ready to use bytes

`Bytes()` method allows converting the id directly into static 11 bytes.

### Multi Node Support

The `monoflake` package can be used on single/multiple nodes without the need for
machine coordination. It uses configured node identifier to generate ids by
attaching the node identifier to the end of the sequences.

### Zero allocation

Zero allocation magic with blazing fast results.

## Performance benchmark

Command:

```
go test -benchtime 10000000x -benchmem -run=^$ -bench=. github.com/mustafaturan/monoflake
```

Results:
```
goos: darwin
goarch: arm64
pkg: github.com/mustafaturan/monoflake
BenchmarkNext-8                 10000000                48.66 ns/op            0 B/op          0 allocs/op
BenchmarkNextCompare-8          10000000                92.34 ns/op            0 B/op          0 allocs/op
BenchmarkNextBase62-8           10000000                51.08 ns/op            0 B/op          0 allocs/op
BenchmarkNextBytes-8            10000000                29.26 ns/op            0 B/op          0 allocs/op
PASS
ok      github.com/mustafaturan/monoflake       2.442s
```

## Contributing

All contributors should follow [Contributing Guidelines](CONTRIBUTING.md) before creating pull requests.

## Credits

[Mustafa Turan](https://github.com/mustafaturan)

## License

Apache License 2.0

Copyright (c) 2024 Mustafa Turan

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
