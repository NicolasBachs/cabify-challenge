package mocks

import (
	"github.com/stretchr/testify/mock"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/domain/entity"
)

type carRepositoryMock struct {
	mock.Mock
}

func NewCarRepositoryMock() *carRepositoryMock {
	return &carRepositoryMock{}
}

func (m *carRepositoryMock) GetCarByID(carID uint) (*entity.Car, error) {
	args := m.Called(carID)
	var arg0 *entity.Car
	if args.Get(0) != nil {
		arg0 = args.Get(0).(*entity.Car)
	}
	var arg1 error
	if args.Error(1) != nil {
		arg1 = args.Error(1)
	}
	return arg0, arg1
}

func (m *carRepositoryMock) CreateCar(car *entity.Car) (uint, error) {
	args := m.Called(car)
	var arg0 uint
	if args.Get(0) != nil {
		arg0 = args.Get(0).(uint)
	}
	var arg1 error
	if args.Error(1) != nil {
		arg1 = args.Error(1)
	}
	return arg0, arg1
}

func (m *carRepositoryMock) DeleteCarsAndInsert(cars []*entity.Car) ([]uint, error) {
	args := m.Called(cars)
	var arg0 []uint
	if args.Get(0) != nil {
		arg0 = args.Get(0).([]uint)
	}
	var arg1 error
	if args.Error(1) != nil {
		arg1 = args.Error(1)
	}
	return arg0, arg1
}

func (m *carRepositoryMock) GetCarsWithAvailableSeatsGreaterOrEqualTo(passengers uint) ([]*entity.Car, error) {
	args := m.Called(passengers)
	var arg0 []*entity.Car
	if args.Get(0) != nil {
		arg0 = args.Get(0).([]*entity.Car)
	}
	var arg1 error
	if args.Error(1) != nil {
		arg1 = args.Error(1)
	}
	return arg0, arg1
}
