package bus_test

import (
	"fmt"
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

	t.Run("with valid generator", func(t *testing.T) {
		c := bus.Config{Next: nil}
		assert.Error(t, fmt.Errorf("id generator func can't be nil"), bus.Configure(c))
	})
}
