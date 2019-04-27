// Copyright 2019 Mustafa Turan. All rights reserved.
// Use of this source code is governed by a Apache License 2.0 license that can
// be found in the LICENSE file.

package bus

// Event is data structure for any logs
type Event struct {
	ID         string      // identifier
	TxID       string      // transaction identifier
	Topic      *Topic      // topic name
	Data       interface{} // actual event data
	OccurredAt int64       // creation time in nanoseconds
}

func (e *Event) emit() {
	for _, h := range e.Topic.handlers {
		h.Handle(e)
	}
}

func (e *Event) transactionize() {
	if e.TxID == "" {
		e.TxID = b.next()
	}
}
