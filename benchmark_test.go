// Copyright 2020 Mustafa Turan. All rights reserved.
// Use of this source code is governed by a Apache License 2.0 license that can
// be found in the LICENSE file.

package bus_test

import (
	"context"
	"testing"

	"github.com/mustafaturan/bus/v2"
)

func BenchmarkEmit(b *testing.B) {
	b.ReportAllocs()

	topic := "order.created"
	ebus := setup(topic)
	defer tearDown(ebus, topic)
	h := fakeHandler(topic)
	ebus.RegisterHandler("test.bench.handler", &h)

	ctx := context.Background()
	ctx = context.WithValue(ctx, bus.CtxKeyTxID, "BENCHMARK")
	for n := 0; n < b.N; n++ {
		data := n
		if _, err := ebus.Emit(ctx, topic, data); err != nil {
			panic(err)
		}
	}
}
