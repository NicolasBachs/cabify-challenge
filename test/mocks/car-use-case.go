package mocks

import (
	"github.com/stretchr/testify/mock"
	appError "gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/domain/app-error"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/domain/entity"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/domain/usecase"
)

type carUseCaseMock struct {
	mock.Mock
}

func NewCarUseCaseMock() *carUseCaseMock {
	return &carUseCaseMock{}
}

func (m *carUseCaseMock) SetJourneyUseCase(journeyUseCase usecase.JourneyUseCase) {
	m.Called(journeyUseCase)
}

func (m *carUseCaseMock) GetCarByID(carID uint) (*entity.Car, appError.CommonError) {
	args := m.Called(carID)
	var arg0 *entity.Car
	if args.Get(0) != nil {
		arg0 = args.Get(0).(*entity.Car)
	}
	var arg1 appError.CommonError
	if args.Get(1) != nil {
		arg1 = args.Get(1).(appError.CommonError)
	}
	return arg0, arg1
}

func (m *carUseCaseMock) CreateCar(car *entity.Car) (uint, appError.CommonError) {
	args := m.Called(car)
	var arg0 uint
	if args.Get(0) != nil {
		arg0 = args.Get(0).(uint)
	}
	var arg1 appError.CommonError
	if args.Get(1) != nil {
		arg1 = args.Get(1).(appError.CommonError)
	}
	return arg0, arg1
}

func (m *carUseCaseMock) DeleteCarsAndInsert(cars []*entity.Car) ([]uint, appError.CommonError) {
	args := m.Called(cars)
	var arg0 []uint
	if args.Get(0) != nil {
		arg0 = args.Get(0).([]uint)
	}
	var arg1 appError.CommonError
	if args.Get(1) != nil {
		arg1 = args.Get(1).(appError.CommonError)
	}
	return arg0, arg1
}

func (m *carUseCaseMock) GetCarsWithSeatingCapacityForAtLeast(passengers uint) ([]*entity.Car, appError.CommonError) {
	args := m.Called(passengers)
	var arg0 []*entity.Car
	if args.Get(0) != nil {
		arg0 = args.Get(0).([]*entity.Car)
	}
	var arg1 appError.CommonError
	if args.Get(1) != nil {
		arg1 = args.Get(1).(appError.CommonError)
	}
	return arg0, arg1
}

func (m *carUseCaseMock) DispatchNewCarAvailableEvent(newCarAvailableID uint) appError.CommonError {
	args := m.Called(newCarAvailableID)
	var arg0 appError.CommonError
	if args.Get(0) != nil {
		arg0 = args.Get(0).(appError.CommonError)
	}
	return arg0
}
