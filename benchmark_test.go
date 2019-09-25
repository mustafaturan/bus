package bus_test

import (
	"testing"

	"github.com/mustafaturan/bus"
)

func BenchmarkEmit(b *testing.B) {
	b.ReportAllocs()

	topic := "order.created"
	setup(topic)
	defer tearDown(topic)
	h := fakeHandler(topic)
	bus.RegisterHandler("test.bench.handler", &h)

	txID := "tx"
	for n := 0; n < b.N; n++ {
		data := n
		if _, err := bus.Emit(topic, data, txID); err != nil {
			panic(err)
		}
	}
}
