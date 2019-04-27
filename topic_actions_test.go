package bus_test

import (
	"testing"

	"github.com/mustafaturan/bus"
	"github.com/stretchr/testify/assert"
)

func TestListTopics(t *testing.T) {
	setup()
	defer tearDown()

	topicNames := []string{"user.created", "user.deleted"}
	defer bus.DeregisterTopics(topicNames...)

	bus.RegisterTopics(topicNames...)
	topics := []*bus.Topic{fetchTopic("user.created"), fetchTopic("user.deleted")}
	assert.Equal(t, topics, bus.ListTopics())
}

func TestRegisterTopics(t *testing.T) {
	setup()
	defer tearDown()

	topicNames := []string{"user.created", "user.deleted"}
	defer bus.DeregisterTopics(topicNames...)

	t.Run("register topics", func(t *testing.T) {
		bus.RegisterTopics(topicNames...)
		assert := assert.New(t)
		for _, n := range topicNames {
			assert.NotNil(fetchTopic(n))
		}
	})
	t.Run("does not register a topic twice", func(t *testing.T) {
		assert := assert.New(t)
		assert.Len(bus.ListTopics(), 2)
		bus.RegisterTopics("user.created")
		assert.Len(bus.ListTopics(), 2)
	})
}

func TestDeregisterTopics(t *testing.T) {
	setup()
	defer tearDown()

	topicNames := []string{"user.created", "user.deleted", "user.updated"}
	defer bus.DeregisterTopics(topicNames...)

	bus.RegisterTopics(topicNames...)
	bus.DeregisterTopics("user.created", "user.updated")
	assert := assert.New(t)
	assert.Nil(fetchTopic("user.created"))
	assert.Nil(fetchTopic("user.updated"))
	assert.NotNil(fetchTopic("user.deleted"))
}
