// Copyright 2019 Mustafa Turan. All rights reserved.
// Use of this source code is governed by a Apache License 2.0 license that can
// be found in the LICENSE file.

package bus

var handlers map[string]*Handler

func init() {
	handlers = make(map[string]*Handler)
}

// ListHandlerKeys returns list of registered handler keys
func ListHandlerKeys() []string {
	var keys []string
	for k := range handlers {
		keys = append(keys, k)
	}
	return keys
}

// RegisterHandler re/register the handler to the registry
func RegisterHandler(key string, h *Handler) {
	b.Lock()
	defer b.Unlock()

	registerHandler(key, h)
}

// DeregisterHandler deletes handler from the registry
func DeregisterHandler(key string) {
	b.Lock()
	defer b.Unlock()

	deregisterHandler(key)
}

func registerHandler(key string, h *Handler) {
	deregisterHandler(key)
	handlers[key] = h
	for _, s := range h.Subscriptions() {
		s.registerHandler(h)
	}
}

func deregisterHandler(key string) {
	if h, ok := handlers[key]; ok {
		for _, s := range h.Subscriptions() {
			s.deregisterHandler(h)
		}
		delete(handlers, key)
	}
}
