package mutex

import (
	"errors"

	"github.com/go-redsync/redsync/v4"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/infrastructure/app"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/infrastructure/mutex/resources"
)

const RedisSyncServiceName = "REDIS_SYNC"

type redisSync struct {
	redsync *redsync.Redsync
}

type redisMutex struct {
	resource resources.MutexResource
	mutex    *redsync.Mutex
}

func NewRedisSync(rs *redsync.Redsync) DistributedResourceSync {
	return &redisSync{
		redsync: rs,
	}
}

func (rs *redisSync) NewMutex(res resources.MutexResource) (DistributedMutex, error) {
	mutex := rs.redsync.NewMutex(res.ToString())

	return &redisMutex{
		resource: res,
		mutex:    mutex,
	}, nil
}

func (m *redisMutex) Lock() error {
	err := m.mutex.Lock()

	if err != nil {
		app.Logger.Error(RedisSyncServiceName, "Error on acquire lock for resource '%s', error: %s", m.resource, err.Error())
		return errors.New("Error on acquire lock")
	}

	app.Logger.Debug(RedisSyncServiceName, "Lock for resource '%s', acquired", m.resource)

	return nil
}

func (m *redisMutex) Unlock() error {
	ok, err := m.mutex.Unlock()

	if err != nil || !ok {
		app.Logger.Error(RedisSyncServiceName, "Error on release lock for resource '%s', error: %s", m.resource, err.Error())
		return errors.New("Error on release lock")
	}

	app.Logger.Debug(RedisSyncServiceName, "Lock for resource '%s', released", m.resource)

	return nil
}
