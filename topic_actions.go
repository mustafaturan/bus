// Copyright 2019 Mustafa Turan. All rights reserved.
// Use of this source code is governed by a Apache License 2.0 license that can
// be found in the LICENSE file.

package bus

var topics map[string]*Topic

func init() {
	topics = make(map[string]*Topic)
}

// ListTopics list registered topics
func ListTopics() []*Topic {
	var l []*Topic
	for _, t := range topics {
		l = append(l, t)
	}
	return l
}

// RegisterTopics registers topics and fullfills handlers
func RegisterTopics(topicNames ...string) {
	for _, n := range topicNames {
		registerTopic(n)
	}
}

// DeregisterTopics deletes topic
func DeregisterTopics(topicNames ...string) {
	for _, n := range topicNames {
		deregisterTopic(n)
	}
}

func registerTopic(topicName string) {
	b.Lock()
	defer b.Unlock()

	if _, ok := topics[topicName]; ok {
		return
	}
	t := &Topic{Name: topicName, handlers: make([]*Handler, 0)}
	t.registerHandlers()
	topics[topicName] = t
}

func deregisterTopic(topicName string) {
	b.Lock()
	defer b.Unlock()

	if _, ok := topics[topicName]; ok {
		delete(topics, topicName)
	}
}
