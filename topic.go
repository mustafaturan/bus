// Copyright 2019 Mustafa Turan. All rights reserved.
// Use of this source code is governed by a Apache License 2.0 license that can
// be found in the LICENSE file.

package bus

import (
	"regexp"
)

// Topic structure
type Topic struct {
	Name     string
	handlers []*Handler
}

// Handlers returns all handlers of the topic
func (t *Topic) Handlers() []*Handler {
	return t.handlers
}

func (t *Topic) registerHandlers() {
	for _, h := range handlers {
		if matched, _ := regexp.MatchString(h.Matcher, t.Name); matched {
			t.registerHandler(h)
		}
	}
}

func (t *Topic) registerHandler(h *Handler) {
	t.handlers = append(t.handlers, h)
}

func (t *Topic) deregisterHandler(h *Handler) {
	for i, handler := range t.handlers {
		if handler == h {
			t.handlers[i] = t.handlers[len(t.handlers)-1]
			t.handlers[len(t.handlers)-1] = nil
			t.handlers = t.handlers[:len(t.handlers)-1]
			break
		}
	}
}
