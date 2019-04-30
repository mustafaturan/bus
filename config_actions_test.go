package bus_test

import (
	"testing"

	"github.com/mustafaturan/bus"
	"github.com/stretchr/testify/assert"
)

func TestConfigure(t *testing.T) {
	fn := func() string { return "newfakeid" }

	t.Run("with valid generator", func(t *testing.T) {
		c := bus.Config{Next: fn}
		assert.Nil(t, bus.Configure(c))
	})

	t.Run("with invalid generator", func(t *testing.T) {
		c := bus.Config{Next: nil}
		err := bus.Configure(c)
		assert.NotNil(t, err)
		if assert.Error(t, err) {
			want := "bus: Next() id generator func can't be nil"
			assert.Equal(t, want, err.Error())
		}
	})
}
