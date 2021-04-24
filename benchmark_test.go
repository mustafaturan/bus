// Copyright 2021 Mustafa Turan. All rights reserved.
// Use of this source code is governed by a Apache License 2.0 license that can
// be found in the LICENSE file.

package bus_test

import (
	"context"
	"testing"
	"time"

	"github.com/mustafaturan/bus/v3"
)

func BenchmarkEmit(b *testing.B) {
	b.ReportAllocs()

	const (
		txID       = "BENCHMARK"
		topic      = "order.created"
		handlerKey = "test.bench.handler"
		source     = "source"
	)

	ebus := setup(topic)
	defer tearDown(ebus, topic)

	h := fakeHandler(topic)
	ebus.RegisterHandler(handlerKey, h)

	ctx := context.WithValue(context.Background(), bus.CtxKeyTxID, txID)
	ctx = context.WithValue(ctx, bus.CtxKeySource, source)
	for n := 0; n < b.N; n++ {
		_ = ebus.Emit(ctx, topic, n)
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
	ebus.RegisterHandler(handlerKey, h)

	ctx := context.Background()
	for n := 0; n < b.N; n++ {
		_ = ebus.Emit(ctx, topic, n)
	}
}

func BenchmarkEmitWithOpts(b *testing.B) {
	b.ReportAllocs()

	const (
		topic      = "order.created"
		handlerKey = "test.bench.handler"
	)

	ebus := setup(topic)
	defer tearDown(ebus, topic)

	h := fakeHandler(topic)
	ebus.RegisterHandler(handlerKey, h)

	ctx := context.Background()
	now := time.Now()
	for n := 0; n < b.N; n++ {
		_ = ebus.EmitWithOpts(
			ctx,
			topic,
			n,
			bus.WithTxID("tx"),
			bus.WithSource("source"),
			bus.WithID("id"),
			bus.WithOccurredAt(now),
		)
	}
}

func BenchmarkEmitWithOptsUnspecified(b *testing.B) {
	b.ReportAllocs()

	const (
		topic      = "order.created"
		handlerKey = "test.bench.handler"
	)

	ebus := setup(topic)
	defer tearDown(ebus, topic)

	h := fakeHandler(topic)
	ebus.RegisterHandler(handlerKey, h)

	ctx := context.Background()
	for n := 0; n < b.N; n++ {
		_ = ebus.EmitWithOpts(ctx, topic, n)
	}
}
