// Copyright 2019 Mustafa Turan. All rights reserved.
// Use of this source code is governed by a Apache License 2.0 license that can
// be found in the LICENSE file.

/*
Package bus is a minimalist event/message bus implementation for internal
communication

Configure

The package requires a unique id generator to assign ids to events. You can
write your own function to generate unique ids or use a package that provides
unique id generation functionality. Here is a sample configuration using
`monoton` id generator:

Example code for configuration:

	import (
		"github.com/mustafaturan/bus"
		"github.com/mustafaturan/monoton"
		"github.com/mustafaturan/monoton/sequencer"
	)

	func init() {
		// configure id generator (it doesn't have to be monoton)
		node        := uint(1)
		initialTime := uint(0)
		monoton.Configure(sequencer.NewMillisecond(), node, initialTime)

		// configure bus
		if err := bus.Configure(bus.Config{Next: monoton.Next}); err != nil {
			panic("whoops")
		}
		// ...
	}

Register Topics

To emit events to the topics, topic names should be registered first:

Example code:

	func init() {
		// ...
		// register topics
		bus.RegisterTopics("order.received", "order.fulfilled")
		// ...
	}

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
	bus.RegisterHandler("a unique key for the handler", &handler)

Emit Event

Example code:

	txID  := "some-transaction-id-if-exists" // if it is blank, package will generate one
	topic := "order.received" // event topic name (must be registered before)
	order := make(map[string]string) // interface{} data for the event

	order["orderID"]     = "123456"
	order["orderAmount"] = "112.20"
	order["currency"]    = "USD"

	// emit the event for the topic with the transaction id
	bus.Emit(topic, order, txID)

Processing Events

When an event is emitted, the topic handlers will receive events synchronously.
It is highly recommended to process events asynchronous. Package leave the
decision to the packages/projects to use concurrency abstractions depending on
use-cases.

*/
package bus
