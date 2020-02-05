// Copyright 2020 Mustafa Turan. All rights reserved.
// Use of this source code is governed by a Apache License 2.0 license that can
// be found in the LICENSE file.

package bus

import (
	"context"
	"fmt"
	"regexp"
	"sync"
	"time"
)

// Bus is a message bus
type Bus struct {
	sync.Mutex
	idgen    IDGenerator
	topics   map[string]*topic
	handlers map[string]*Handler
}

// Next is a sequential unique id generator func type
type Next func() string

// IDGenerator is a sequential unique id generator interface
type IDGenerator interface {
	Generate() string
}

// Event is data structure for any logs
type Event struct {
	ctx        context.Context
	ID         string      // identifier
	TxID       string      // transaction identifier
	Topic      string      // topic name
	Data       interface{} // actual event data
	OccurredAt int64       // creation time in nanoseconds
}

// Handler is a receiver for event reference with the given regex pattern
type Handler struct {
	Handle  func(e *Event) // handler func to process events
	Matcher string         // topic matcher as regex pattern
}

// topic structure
type topic struct {
	name     string
	handlers []*Handler
}

type ctxKey string

const (
	// CtxKeyTxID tx id context key
	CtxKeyTxID = ctxKey("bus.txID")

	// Version syncs with package version
	Version = "1.0.2"
)

// NewBus inits a new bus
func NewBus(g IDGenerator) (*Bus, error) {
	if g == nil {
		return nil, fmt.Errorf("bus: Next() id generator func can't be nil")
	}

	return &Bus{
		idgen:    g,
		topics:   make(map[string]*topic),
		handlers: make(map[string]*Handler),
	}, nil
}

// Emit inits a new event and delivers to the interested in handlers
func (b *Bus) Emit(ctx context.Context, topicName string, data interface{}) (*Event, error) {
	if _, ok := b.topics[topicName]; !ok {
		return nil, fmt.Errorf("bus: topic(%s) not found", topicName)
	}

	e := b.newEvent(ctx, topicName, data)
	b.emit(e)
	return e, nil
}

// Topics lists the all registered topics
func (b *Bus) Topics() []string {
	topics, index := make([]string, len(b.topics)), 0

	for topicName := range b.topics {
		topics[index] = topicName
		index++
	}
	return topics
}

// RegisterTopics registers topics and fullfills handlers
func (b *Bus) RegisterTopics(topicNames ...string) {
	for _, n := range topicNames {
		b.registerTopic(n)
	}
}

// DeregisterTopics deletes topic
func (b *Bus) DeregisterTopics(topicNames ...string) {
	for _, n := range topicNames {
		b.deregisterTopic(n)
	}
}

// TopicHandlers returns all handlers for the topic
func (b *Bus) TopicHandlers(topicName string) []*Handler {
	return b.topics[topicName].handlers
}

// HandlerKeys returns list of registered handler keys
func (b *Bus) HandlerKeys() []string {
	keys, index := make([]string, len(b.handlers)), 0

	for k := range b.handlers {
		keys[index] = k
		index++
	}
	return keys
}

// HandlerTopicSubscriptions returns all topic subscriptions of the handler
func (b *Bus) HandlerTopicSubscriptions(handlerKey string) []string {
	var subscriptions []string
	h, ok := b.handlers[handlerKey]
	if !ok {
		return subscriptions
	}

	for _, t := range b.topics {
		if matched, _ := regexp.MatchString(h.Matcher, t.name); matched {
			subscriptions = append(subscriptions, t.name)
		}
	}
	return subscriptions
}

// RegisterHandler re/register the handler to the registry
func (b *Bus) RegisterHandler(key string, h *Handler) {
	b.Lock()
	defer b.Unlock()

	b.registerHandler(key, h)
}

// DeregisterHandler deletes handler from the registry
func (b *Bus) DeregisterHandler(key string) {
	b.Lock()
	defer b.Unlock()

	b.deregisterHandler(key)
}

// Context returns event's context
func (e *Event) Context() context.Context {
	return e.ctx
}

// Generate is an implementation of IDGenerator for bus.Next fn type
func (n Next) Generate() string {
	return n()
}

func (b *Bus) registerHandler(key string, h *Handler) {
	b.deregisterHandler(key)
	b.handlers[key] = h
	for _, t := range b.HandlerTopicSubscriptions(key) {
		b.registerTopicHandler(b.topics[t], h)
	}
}

func (b *Bus) deregisterHandler(handlerKey string) {
	if h, ok := b.handlers[handlerKey]; ok {
		for _, t := range b.HandlerTopicSubscriptions(handlerKey) {
			b.deregisterTopicHandler(b.topics[t], h)
		}
		delete(b.handlers, handlerKey)
	}
}

func (b *Bus) registerTopicHandlers(t *topic) {
	for _, h := range b.handlers {
		if matched, _ := regexp.MatchString(h.Matcher, t.name); matched {
			b.registerTopicHandler(t, h)
		}
	}
}

func (b *Bus) registerTopicHandler(t *topic, h *Handler) {
	t.handlers = append(t.handlers, h)
}

func (b *Bus) deregisterTopicHandler(t *topic, h *Handler) {
	for i, handler := range t.handlers {
		if handler == h {
			t.handlers[i] = t.handlers[len(t.handlers)-1]
			t.handlers[len(t.handlers)-1] = nil
			t.handlers = t.handlers[:len(t.handlers)-1]
			break
		}
	}
}

func (b *Bus) newEvent(ctx context.Context, topicName string, data interface{}) *Event {
	e := &Event{
		ctx:        ctx,
		ID:         b.idgen.Generate(),
		Topic:      topicName,
		Data:       data,
		OccurredAt: time.Now().UnixNano(),
	}
	val := ctx.Value(CtxKeyTxID)
	if txID, ok := val.(string); ok {
		e.TxID = txID
	} else {
		e.TxID = b.idgen.Generate()
	}
	return e
}

func (b *Bus) emit(e *Event) {
	for _, h := range b.topics[e.Topic].handlers {
		h.Handle(e)
	}
}

func (b *Bus) registerTopic(topicName string) {
	b.Lock()
	defer b.Unlock()

	if _, ok := b.topics[topicName]; ok {
		return
	}
	t := &topic{name: topicName, handlers: []*Handler{}}
	b.registerTopicHandlers(t)
	b.topics[topicName] = t
}

func (b *Bus) deregisterTopic(topicName string) {
	b.Lock()
	defer b.Unlock()

	delete(b.topics, topicName)
}
