package bus_test

import (
	"github.com/mustafaturan/bus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandlers(t *testing.T) {
	setup()
	defer tearDown()
	defer bus.DeregisterHandler("test.handler/1")
	defer bus.DeregisterHandler("test.handler/2")

	handler := fakeHandler(".*")
	bus.RegisterHandler("test.handler/1", &handler)
	bus.RegisterHandler("test.handler/2", &handler)
	bus.RegisterTopics("user.created")

	assert := assert.New(t)
	for _, h := range fetchTopic("user.created").Handlers() {
		assert.Equal(&handler, h)
	}
}
