package bus_test

import (
	"testing"

	"github.com/mustafaturan/bus"
	"github.com/stretchr/testify/assert"
)

func TestSubscriptions(t *testing.T) {
	setup("comment.created", "comment.deleted")
	defer tearDown("comment.created", "comment.deleted")
	defer bus.DeregisterHandler("test.handler")

	h := fakeHandler(".*")
	bus.RegisterHandler("test.handler", &h)

	assert.Equal(t, bus.ListTopics(), h.Subscriptions())
}
