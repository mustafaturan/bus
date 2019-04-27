package bus_test

import (
	"testing"

	"github.com/mustafaturan/bus"
	"github.com/stretchr/testify/assert"
)

func TestListHandlerKeys(t *testing.T) {
	setup("comment.created", "comment.deleted")
	defer tearDown("comment.created", "comment.deleted")
	defer bus.DeregisterHandler("test.key.1")
	defer bus.DeregisterHandler("test.key.2")

	h := fakeHandler(".*")
	bus.RegisterHandler("test.key.1", &h)
	bus.RegisterHandler("test.key.2", &h)

	want := []string{"test.key.1", "test.key.2"}
	assert.ElementsMatch(t, want, bus.ListHandlerKeys())
}

func TestRegisterHandler(t *testing.T) {
	setup("comment.created", "comment.deleted")
	defer tearDown("comment.created", "comment.deleted")
	defer bus.DeregisterHandler("test.handler")

	h := fakeHandler(".*created$")
	bus.RegisterHandler("test.handler", &h)

	t.Run("registers handler key", func(t *testing.T) {
		assert.True(t, isHandlerKey("test.handler"))
	})
	t.Run("adds handler references to the matched topics", func(t *testing.T) {
		t.Run("when topic is matched", func(t *testing.T) {
			topic := fetchTopic("comment.created")
			assert.True(t, isTopicHandler(topic, &h))
		})
		t.Run("when topic is not matched", func(t *testing.T) {
			topic := fetchTopic("comment.deleted")
			assert.False(t, isTopicHandler(topic, &h))
		})
	})
}

func TestDeregisterHandler(t *testing.T) {
	setup("comment.created", "comment.deleted")
	defer tearDown("comment.created", "comment.deleted")

	h := fakeHandler(".*")
	bus.RegisterHandler("test.handler", &h)
	bus.DeregisterHandler("test.handler")

	t.Run("deletes handler key", func(t *testing.T) {
		assert.False(t, isHandlerKey("test.handler"))
	})
	t.Run("deletes handler references from the topics", func(t *testing.T) {
		assert := assert.New(t)
		for _, topic := range bus.ListTopics() {
			assert.False(isTopicHandler(topic, &h))
		}
	})
}
