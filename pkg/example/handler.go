package example

import (
	"context"

	"github.com/ridebeam/go-common/kv"
	"github.com/ridebeam/golang-skeleton/pkg/serde"
)

type ExampleEmitter func(context.Context, serde.ExampleOutput) error

type Example struct {
	kv      kv.VersionedKeyValueStore
	emitter ExampleEmitter
}

func NewExample(kv kv.VersionedKeyValueStore, emitter ExampleEmitter) *Example {
	return &Example{
		kv:      kv,
		emitter: emitter,
	}
}

func (e *Example) Handle(ctx context.Context, m serde.ExampleModel) error {
	// we setup the actual handler to receive and emit in-memory representation and let serialization be handled outside
	// this decouples the actual business logic, allowing for better testability
	return e.emitter(ctx, serde.ExampleOutput{
		ID:     m.ID,
		FooBar: m.Bar,
	})
}
