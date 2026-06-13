package utils

import (
	"context"
	"fmt"
	"log"
)

func TryExecute(ctx context.Context, fn func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic recovered: %v", r)
			log.Printf("critical panic: %v", r)
		}
	}()

	done := make(chan error, 1)
	go func() {
		done <- fn()
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("context cancelled before execution: %w", ctx.Err())
	case err := <-done:
		return err
	}
}