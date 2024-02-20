package cms

import "context"

type Hook[T any] interface {
	Trigger(context.Context, T, ...func(context.Context, T) error) error
}

type StubHook[T any] struct{}

func (h StubHook[T]) Trigger(ctx context.Context, event T, oneOff ...func(context.Context, T) error) error {
	for _, f := range oneOff {
		if err := f(ctx, event); err != nil {
			return err
		}
	}
	return nil
}
