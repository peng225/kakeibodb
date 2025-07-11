package fake

import (
	"context"
	"fmt"
)

type FakeTransaction struct{}

func NewFakeTransaction() *FakeTransaction {
	return &FakeTransaction{}
}

func (ft *FakeTransaction) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	err := fn(ctx)
	if err != nil {
		return fmt.Errorf("function failed: %w", err)
	}
	return nil
}
