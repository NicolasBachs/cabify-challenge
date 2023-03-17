package mutex

import (
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/infrastructure/mutex/resources"
)

type DistributedResourceSync interface {
	NewMutex(resource resources.MutexResource) (DistributedMutex, error)
}

type DistributedMutex interface {
	Lock() error
	Unlock() error
}
