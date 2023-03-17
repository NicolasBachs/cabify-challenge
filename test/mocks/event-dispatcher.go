package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type eventDispatcherMock struct {
	mock.Mock
}

func NewEventDispatcherMock() *eventDispatcherMock {
	return &eventDispatcherMock{}
}

func (m *eventDispatcherMock) Dispatch(ctx context.Context, topic string, event interface{}) error {
	args := m.Called(ctx, topic, event)
	var arg0 error
	if args.Error(0) != nil {
		arg0 = args.Error(0)
	}
	return arg0
}

func (m *eventDispatcherMock) Close() error {
	args := m.Called()
	var arg0 error
	if args.Error(0) != nil {
		arg0 = args.Error(0)
	}
	return arg0
}
