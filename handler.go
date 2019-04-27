// Copyright 2019 Mustafa Turan. All rights reserved.
// Use of this source code is governed by a Apache License 2.0 license that can
// be found in the LICENSE file.

package bus

import (
	"regexp"
)

// Handler is a receiver for event reference with the given regex pattern
type Handler struct {
	Handle  func(e *Event) // handler func to process events
	Matcher string         // topic matcher as regex pattern
}

// Subscriptions returns all subscriptions of the handler
func (h *Handler) Subscriptions() []*Topic {
	var subscriptions []*Topic
	for _, t := range topics {
		if matched, _ := regexp.MatchString(h.Matcher, t.Name); matched {
			subscriptions = append(subscriptions, t)
		}
	}
	return subscriptions
}
