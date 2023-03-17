package mocks

import (
	"github.com/stretchr/testify/mock"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/infrastructure/mutex"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/infrastructure/mutex/resources"
)

type distributedResourceSyncMock struct {
	mock.Mock
	mutex *distributedMutexMock
}

type distributedMutexMock struct {
	mock.Mock
}

func NewDistributedResourceSyncMock(mutex *distributedMutexMock) *distributedResourceSyncMock {
	return &distributedResourceSyncMock{
		mutex: mutex,
	}
}

func NewDistributedMutexMock() *distributedMutexMock {
	return &distributedMutexMock{}
}

func (m *distributedResourceSyncMock) NewMutex(res resources.MutexResource) (mutex.DistributedMutex, error) {
	args := m.Called(res)
	var arg0 mutex.DistributedMutex
	if args.Get(0) != nil {
		arg0 = args.Get(0).(mutex.DistributedMutex)
	}
	var arg1 error
	if args.Get(1) != nil {
		arg1 = args.Error(1)
	}
	return arg0, arg1
}

func (m *distributedMutexMock) Lock() error {
	args := m.Called()
	var arg0 error
	if args.Error(0) != nil {
		arg0 = args.Error(0)
	}
	return arg0
}

func (m *distributedMutexMock) Unlock() error {
	args := m.Called()
	var arg0 error
	if args.Error(0) != nil {
		arg0 = args.Error(0)
	}
	return arg0
}
