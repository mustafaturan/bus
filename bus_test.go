// Copyright 2021 Mustafa Turan. All rights reserved.
// Use of this source code is governed by a Apache License 2.0 license that can
// be found in the LICENSE file.

package bus_test

import (
	"context"
	"testing"
	"time"

	"github.com/mustafaturan/bus/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCtxKeyTxID(t *testing.T) {
	assert.EqualValues(t, bus.CtxKeyTxID, rune(66))
}

func TestVersion(t *testing.T) {
	assert.Equal(t, bus.Version, "2.0.4")
}

func TestNew(t *testing.T) {
	var fn bus.Next = func() string { return "afakeid" }

	t.Run("with valid generator", func(t *testing.T) {
		b, err := bus.NewBus(fn)
		require.Nil(t, err)
		assert.IsType(t, &bus.Bus{}, b)
	})

	t.Run("with invalid generator", func(t *testing.T) {
		b, err := bus.NewBus(nil)
		require.Nil(t, b)
		assert.NotNil(t, err)
		if assert.Error(t, err) {
			want := "bus: Next() id generator func can't be nil"
			assert.Equal(t, want, err.Error())
		}
	})
}

func TestEmit(t *testing.T) {
	b := setup("comment.created", "comment.deleted")
	defer tearDown(b, "comment.created", "comment.deleted")

	t.Run("correctly assigns fields", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, bus.CtxKeyTxID, "tx")
		e, err := b.Emit(ctx, "comment.deleted", "my comment")

		assert := assert.New(t)
		assert.Equal("tx", e.TxID)
		assert.Equal("fakeid", e.ID)
		assert.Equal("comment.deleted", e.Topic)
		assert.Equal("my comment", e.Data)
		assert.True(e.OccurredAt <= time.Now().UnixNano())
		assert.Nil(err)
	})
	t.Run("updates txID when empty", func(t *testing.T) {
		ctx := context.Background()
		e, err := b.Emit(ctx, "comment.deleted", "my comment")

		assert := assert.New(t)
		assert.Equal("fakeid", e.TxID)
		assert.Nil(err)
	})
	t.Run("with handler", func(t *testing.T) {
		ctx := context.Background()
		registerFakeHandler(b, "test", t)

		_, err := b.Emit(ctx, "comment.created", "my comment with handler")
		if err != nil {
			t.Fatalf("emit failed: %v", err)
		}
		b.DeregisterHandler("test")
	})
	t.Run("with unknown topic", func(t *testing.T) {
		ctx := context.Background()
		e, err := b.Emit(ctx, "comment.updated", "my comment")

		assert := assert.New(t)
		assert.Nil(e)
		assert.NotNil(err)
		assert.Equal("bus: topic(comment.updated) not found", err.Error())
	})
}

func TestTopics(t *testing.T) {
	topicNames := []string{"user.created", "user.deleted"}
	b := setup(topicNames...)
	defer tearDown(b, topicNames...)

	assert.ElementsMatch(t, topicNames, b.Topics())
}

func TestRegisterTopics(t *testing.T) {
	b := setup()
	defer tearDown(b)

	topicNames := []string{"user.created", "user.deleted"}
	defer b.DeregisterTopics(topicNames...)

	t.Run("register topics", func(t *testing.T) {
		b.RegisterTopics(topicNames...)
		assert.ElementsMatch(t, topicNames, b.Topics())
	})
	t.Run("does not register a topic twice", func(t *testing.T) {
		assert := assert.New(t)
		assert.Len(b.Topics(), 2)
		b.RegisterTopics("user.created")
		assert.Len(b.Topics(), 2)
		assert.ElementsMatch(topicNames, b.Topics())
	})
}

func TestDeregisterTopics(t *testing.T) {
	b := setup()
	defer tearDown(b)

	topicNames := []string{"user.created", "user.deleted", "user.updated"}
	defer b.DeregisterTopics(topicNames...)

	b.RegisterTopics(topicNames...)
	b.DeregisterTopics("user.created", "user.updated")
	assert := assert.New(t)
	assert.ElementsMatch([]string{"user.deleted"}, b.Topics())
}

func TestTopicHandlers(t *testing.T) {
	b := setup()
	defer tearDown(b)
	defer b.DeregisterHandler("test.handler/1")
	defer b.DeregisterHandler("test.handler/2")

	handler := fakeHandler(".*")
	b.RegisterHandler("test.handler/1", &handler)
	b.RegisterHandler("test.handler/2", &handler)
	b.RegisterTopics("user.created")

	assert := assert.New(t)
	for _, h := range b.TopicHandlers("user.created") {
		assert.Equal(&handler, h)
	}
}

func TestHandlerKeys(t *testing.T) {
	b := setup("comment.created", "comment.deleted")
	defer tearDown(b, "comment.created", "comment.deleted")
	defer b.DeregisterHandler("test.key.1")
	defer b.DeregisterHandler("test.key.2")

	h := fakeHandler(".*")
	b.RegisterHandler("test.key.1", &h)
	b.RegisterHandler("test.key.2", &h)

	want := []string{"test.key.1", "test.key.2"}
	assert.ElementsMatch(t, want, b.HandlerKeys())
}

func TestHandlerTopicSubscriptions(t *testing.T) {
	b := setup("comment.created", "comment.deleted")
	defer tearDown(b, "comment.created", "comment.deleted")

	tests := []struct {
		handler          bus.Handler
		handlerKey       string
		handlerLookupKey string
		want             []string
	}{
		{fakeHandler(".*"), "test.handler.1", "test.handler.1", []string{"comment.created", "comment.deleted"}},
		{fakeHandler("user.updated"), "test.handler.2", "test.handler.2", []string{}},
		{fakeHandler(".*"), "test.handler.3", "test.handler.NA", []string{}},
	}

	for _, test := range tests {
		b.RegisterHandler(test.handlerKey, &test.handler)

		assert.ElementsMatch(t, test.want, b.HandlerTopicSubscriptions(test.handlerLookupKey))
	}
}

func TestRegisterHandler(t *testing.T) {
	b := setup("comment.created", "comment.deleted")
	defer tearDown(b, "comment.created", "comment.deleted")
	defer b.DeregisterHandler("test.handler")

	h := fakeHandler(".*created$")
	b.RegisterHandler("test.handler", &h)

	t.Run("registers handler key", func(t *testing.T) {
		assert.True(t, isHandlerKeyExists(b, "test.handler"))
	})
	t.Run("adds handler references to the matched topics", func(t *testing.T) {
		t.Run("when topic is matched", func(t *testing.T) {
			assert.True(t, isTopicHandler(b, "comment.created", &h))
		})
		t.Run("when topic is not matched", func(t *testing.T) {
			assert.False(t, isTopicHandler(b, "comment.deleted", &h))
		})
	})
}

func TestDeregisterHandler(t *testing.T) {
	b := setup("comment.created", "comment.deleted")
	defer tearDown(b, "comment.created", "comment.deleted")

	h := fakeHandler(".*")
	b.RegisterHandler("test.handler", &h)
	b.DeregisterHandler("test.handler")

	t.Run("deletes handler key", func(t *testing.T) {
		assert.False(t, isHandlerKeyExists(b, "test.handler"))
	})
	t.Run("deletes handler references from the topics", func(t *testing.T) {
		assert := assert.New(t)
		for _, topic := range b.Topics() {
			assert.False(isTopicHandler(b, topic, &h))
		}
	})
}

func setup(topicNames ...string) *bus.Bus {
	var fn bus.Next = func() string { return "fakeid" }
	b, _ := bus.NewBus(fn)
	b.RegisterTopics(topicNames...)
	return b
}

func tearDown(b *bus.Bus, topicNames ...string) {
	b.DeregisterTopics(topicNames...)
}

func fakeHandler(matcher string) bus.Handler {
	return bus.Handler{Handle: func(context.Context, *bus.Event) {}, Matcher: matcher}
}

func registerFakeHandler(b *bus.Bus, key string, t *testing.T) {
	fn := func(ctx context.Context, e *bus.Event) {
		t.Run("receives right event", func(t *testing.T) {
			assert := assert.New(t)
			assert.Equal("fakeid", e.ID)
			assert.Equal("comment.created", e.Topic)
			assert.Equal("my comment with handler", e.Data)
			assert.True(e.OccurredAt <= time.Now().UnixNano())
		})
	}
	h := bus.Handler{Handle: fn, Matcher: ".*created$"}
	b.RegisterHandler(key, &h)
}

func isTopicHandler(b *bus.Bus, topicName string, h *bus.Handler) bool {
	for _, th := range b.TopicHandlers(topicName) {
		if h == th {
			return true
		}
	}
	return false
}

func isHandlerKeyExists(b *bus.Bus, key string) bool {
	for _, k := range b.HandlerKeys() {
		if k == key {
			return true
		}
	}
	return false
}
