package bus_test

import (
	"github.com/mustafaturan/bus"
)

func setup(topicNames ...string) {
	fn := func() string { return "fakeid" }
	if err := bus.Configure(bus.Config{Next: fn}); err != nil {
		panic(err)
	}
	bus.RegisterTopics(topicNames...)
}

func tearDown(topicNames ...string) {
	bus.DeregisterTopics(topicNames...)
}

func fakeHandler(matcher string) bus.Handler {
	return bus.Handler{Handle: func(*bus.Event) {}, Matcher: matcher}
}

func fetchTopic(topicName string) *bus.Topic {
	for _, t := range bus.ListTopics() {
		if t.Name == topicName {
			return t
		}
	}
	return nil
}

func isTopicHandler(t *bus.Topic, h *bus.Handler) bool {
	for _, th := range t.Handlers() {
		if h == th {
			return true
		}
	}
	return false
}

func isHandlerKey(key string) bool {
	for _, k := range bus.ListHandlerKeys() {
		if k == key {
			return true
		}
	}
	return false
}
