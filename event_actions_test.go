package bus_test

import (
	"testing"
	"time"

	"github.com/mustafaturan/bus"
	"github.com/stretchr/testify/assert"
)

func TestEmit(t *testing.T) {
	setup("comment.created", "comment.deleted")
	defer tearDown("comment.created", "comment.deleted")

	t.Run("correctly assigns fields", func(t *testing.T) {
		e, err := bus.Emit("comment.deleted", "my comment", "tx")

		assert := assert.New(t)
		assert.Equal("tx", e.TxID)
		assert.Equal("fakeid", e.ID)
		assert.Equal("comment.deleted", e.Topic.Name)
		assert.Equal("my comment", e.Data)
		assert.True(e.OccurredAt <= time.Now().UnixNano())
		assert.Nil(err)
	})
	t.Run("updates txID when empty", func(t *testing.T) {
		e, err := bus.Emit("comment.deleted", "my comment", "")

		assert := assert.New(t)
		assert.Equal("fakeid", e.TxID)
		assert.Nil(err)
	})
	t.Run("with handler", func(t *testing.T) {
		registerHandler("test", t)

		bus.Emit("comment.created", "my comment with handler", "tx")
		bus.DeregisterHandler("test")
	})
	t.Run("with unknown topic", func(t *testing.T) {
		e, err := bus.Emit("comment.updated", "my comment", "tx")

		assert := assert.New(t)
		assert.Nil(e)
		assert.NotNil(err)
	})
}

func registerHandler(key string, t *testing.T) {
	fn := func(e *bus.Event) {
		t.Run("receives right event", func(t *testing.T) {
			assert := assert.New(t)
			assert.Equal("tx", e.TxID)
			assert.Equal("fakeid", e.ID)
			assert.Equal("comment.created", e.Topic.Name)
			assert.Equal("my comment with handler", e.Data)
			assert.True(e.OccurredAt <= time.Now().UnixNano())
		})
	}
	h := bus.Handler{Handle: fn, Matcher: ".*created$"}
	bus.RegisterHandler(key, &h)
}
