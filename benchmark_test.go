// Copyright 2021 Mustafa Turan. All rights reserved.
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

	const (
		txID       = "BENCHMARK"
		topic      = "order.created"
		handlerKey = "test.bench.handler"
	)

	ebus := setup(topic)
	defer tearDown(ebus, topic)

	h := fakeHandler(topic)
	ebus.RegisterHandler(handlerKey, &h)

	ctx := context.Background()
	ctx = context.WithValue(ctx, bus.CtxKeyTxID, txID)
	for n := 0; n < b.N; n++ {
		data := n
		if _, err := ebus.Emit(ctx, topic, data); err != nil {
			panic(err)
		}
	}
}

func BenchmarkEmitWithoutTxID(b *testing.B) {
	b.ReportAllocs()

	const (
		topic      = "order.created"
		handlerKey = "test.bench.handler"
	)

	ebus := setup(topic)
	defer tearDown(ebus, topic)

	h := fakeHandler(topic)
	ebus.RegisterHandler(handlerKey, &h)

	ctx := context.Background()
	for n := 0; n < b.N; n++ {
		data := n
		if _, err := ebus.Emit(ctx, topic, data); err != nil {
			panic(err)
		}
	}
}
