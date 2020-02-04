# 🔊 Bus

[![GoDoc](https://godoc.org/github.com/mustafaturan/bus?status.svg)](https://godoc.org/github.com/mustafaturan/bus)
[![Build Status](https://travis-ci.org/mustafaturan/bus.svg?branch=master)](https://travis-ci.org/mustafaturan/bus)
[![Coverage Status](https://coveralls.io/repos/github/mustafaturan/bus/badge.svg?branch=master)](https://coveralls.io/github/mustafaturan/bus?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/mustafaturan/bus)](https://goreportcard.com/report/github.com/mustafaturan/bus)
[![GitHub license](https://img.shields.io/github/license/mustafaturan/bus.svg)](https://github.com/mustafaturan/bus/blob/master/LICENSE)

Bus is a minimalist event/message bus implementation for internal communication.
It is heavily inspired from my [event_bus](https://github.com/otobus/event_bus)
package for Elixir language.

## API

The method names and arities/args are stable now. No change should be expected
on the package for the version `1.x.x` except any bug fixes.

## Installation

Via go packages:
```go get github.com/mustafaturan/bus```

## Usage

### Configure

The package requires a unique id generator to assign ids to events. You can
write your own function to generate unique ids or use a package that provides
unique id generation functionality.

The `bus` package respect to software design choice of the packages/projects. It
supports both singleton and dependency injection to init a `bus` instance.

*Hint:*
Check the [demo project](https://github.com/mustafaturan/bus-sample-project) for
the singleton configuration.

Here is a sample initilization using `monoton` id generator:

```go
import (
    "github.com/mustafaturan/bus"
    "github.com/mustafaturan/monoton"
    "github.com/mustafaturan/monoton/sequencer"
)

func NewBus() *bus.Bus {
    // configure id generator (it doesn't have to be monoton)
    node        := uint64(1)
    initialTime := uint64(0)
    monoton.Configure(sequencer.NewMillisecond(), node, initialTime)

    // init an id generator
    var idGenerator bus.Next = monoton.Next

    // create a new bus instance
    b, err := bus.NewBus(idGenerator)
    if err != nil {
        panic(err)
    }

    // maybe register topics in here
    b.RegisterTopics("order.received", "order.fulfilled")

    return b
}
```

### Register Event Topics

To emit events to the topics, topic names need to be registered first:

```go
// register topics
b.RegisterTopics("order.received", "order.fulfilled")
```

### Register Event Handlers

To receive topic events you need to register handlers; A handler basically
requires two vals which are a `Handle` function and topic `Matcher` regex
pattern.

```go
handler := bus.Handler{
    Handle: func(e *bus.Event) {
        // do something
        // NOTE: Highly recommended to process the event in an async way
    },
    Matcher: ".*", // matches all topics
}
b.RegisterHandler("a unique key for the handler", &handler)
```

### Emit Events

```go
// if txID val is blank, bus package generates one using the id generator
ctx := context.Background()
ctx = context.WithValue(ctx, bus.CtxKeyTxID, "some-transaction-id-if-exists")

// event topic name (must be registered before)
topic := "order.received"

// interface{} data for event
order := make(map[string]string)
order["orderID"]     = "123456"
order["orderAmount"] = "112.20"
order["currency"]    = "USD"

// emit the event
event, err := b.Emit(ctx, topic, order)

if err != nil {
    // report the err
    fmt.Println(err)
}

// if the caller needs the event, a ref for the event is returning as result of
// the `Emit` call.
fmt.Println(event)
```

### Processing Events

When an event is emitted, the topic handlers receive the event synchronously.
It is highly recommended to process events asynchronous. Package leave the
decision to the packages/projects to use concurrency abstractions depending on
use-cases. Each handlers receive the same event as ref of `bus.Event` struct:

```go
// Event data structure
type Event struct {
	ID         string      // identifier
	TxID       string      // transaction identifier
	Topic      string      // topic name
	Data       interface{} // actual event data
	OccurredAt int64       // creation time in nanoseconds
}
```

### Sample Project

A [demo project](https://github.com/mustafaturan/bus-sample-project) with three
consumers which increments a `counter` for each event topic, `printer` consumer
which prints all events and lastly `calculator` consumer which sums amounts.

### Benchmarks

```
BenchmarkEmit-4   	 5983903	       200 ns/op	     104 B/op	       2 allocs/op
```

## Contributing

All contributors should follow [Contributing Guidelines](CONTRIBUTING.md) before creating pull requests.

## Credits

[Mustafa Turan](https://github.com/mustafaturan)

## License

Apache License 2.0

Copyright (c) 2020 Mustafa Turan

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
