// Copyright 2020 Mustafa Turan. All rights reserved.
// Use of this source code is governed by a Apache License 2.0 license that can
// be found in the LICENSE file.

/*
Package bus is a minimalist event/message bus implementation for internal
communication

The package requires a unique id generator to assign ids to events. You can
write your own function to generate unique ids or use a package that provides
unique id generation functionality.

The `bus` package respect to software design choice of the packages/projects. It
supports both singleton and dependency injection to init a `bus` instance.

Here is a sample initilization using `monoton` id generator:

Example code for configuration:

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

Register Topics

To emit events to the topics, topic names should be registered first:

Example code:


	// register topics
	b.RegisterTopics("order.received", "order.fulfilled")
	// ...

Register Handlers

To receive topic events you need to register handlers; A handler basically
requires two vals which are a `Handle` function and topic `Matcher` regex
pattern.

Example code:

	handler := bus.Handler{
		Handle: func(e *Event) {
			// do something
			// NOTE: Highly recommended to process the event in an async way
		},
		Matcher: ".*", // matches all topics
	}
	b.RegisterHandler("a unique key for the handler", &handler)

Emit Event

Example code:

	// if txID val is blank, bus package generates one using the id generator
	ctx := context.Background()
	ctx = context.WithValue(ctx, bus.CtxKeyTxID, "a-transaction-id")

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

	// in case of need to do anything with event on caller, a ref is also
	// returning on `Emit` call.
	fmt.Println(event)

Processing Events

When an event is emitted, the topic handlers receive the event synchronously.
It is highly recommended to process events asynchronous. Package leave the
decision to the packages/projects to use concurrency abstractions depending on
use-cases. Each handlers receive the same event as ref of `bus.Event` struct.

*/
package bus
