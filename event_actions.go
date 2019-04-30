// Copyright 2019 Mustafa Turan. All rights reserved.
// Use of this source code is governed by a Apache License 2.0 license that can
// be found in the LICENSE file.

package bus

import (
	"fmt"
	"time"
)

// Emit inits a new event and delivers to the interested in handlers
func Emit(topicName string, data interface{}, txID string) (*Event, error) {
	if topic, ok := topics[topicName]; ok {
		e := newEvent(txID, topic, data)
		e.transactionize()
		e.emit()
		return e, nil
	}

	return nil, fmt.Errorf("bus: topic(%s) not found", topicName)
}

func newEvent(txID string, t *Topic, data interface{}) *Event {
	return &Event{
		TxID:       txID,
		ID:         b.next(),
		Topic:      t,
		Data:       data,
		OccurredAt: time.Now().UnixNano(),
	}
}
