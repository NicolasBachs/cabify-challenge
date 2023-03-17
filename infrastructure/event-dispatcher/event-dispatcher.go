package eventDispatcher

import "context"

type EventDispatcher interface {
	Dispatch(ctx context.Context, topic string, event interface{}) error
	Close() error
}
