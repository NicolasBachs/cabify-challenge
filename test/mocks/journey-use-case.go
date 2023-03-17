package mocks

import (
	"github.com/stretchr/testify/mock"
	appError "gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/domain/app-error"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/domain/entity"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/domain/enum"
)

type journeyUseCaseMock struct {
	mock.Mock
}

func NewJourneyUseCaseMock() *journeyUseCaseMock {
	return &journeyUseCaseMock{}
}

func (m *journeyUseCaseMock) GetJourneyByID(journeyID uint) (*entity.Journey, appError.CommonError) {
	args := m.Called(journeyID)
	var arg0 *entity.Journey
	var arg1 appError.CommonError
	if args.Get(0) != nil {
		arg0 = args.Get(0).(*entity.Journey)
	}
	if args.Get(1) != nil {
		arg1 = args.Get(1).(appError.CommonError)
	}
	return arg0, arg1
}

func (m *journeyUseCaseMock) GetJourneysByCarID(carID uint) ([]*entity.Journey, appError.CommonError) {
	args := m.Called(carID)
	var arg0 []*entity.Journey
	var arg1 appError.CommonError
	if args.Get(0) != nil {
		arg0 = args.Get(0).([]*entity.Journey)
	}
	if args.Get(1) != nil {
		arg1 = args.Get(1).(appError.CommonError)
	}
	return arg0, arg1
}

func (m *journeyUseCaseMock) CreateJourney(journey *entity.Journey) (uint, appError.CommonError) {
	args := m.Called(journey)
	var arg0 uint
	var arg1 appError.CommonError
	if args.Get(0) != nil {
		arg0 = args.Get(0).(uint)
	}
	if args.Get(1) != nil {
		arg1 = args.Get(1).(appError.CommonError)
	}
	return arg0, arg1
}

func (m *journeyUseCaseMock) CleanJourneys() appError.CommonError {
	args := m.Called()
	var arg0 appError.CommonError
	if args.Get(0) != nil {
		arg0 = args.Get(0).(appError.CommonError)
	}
	return arg0
}

func (m *journeyUseCaseMock) GetLastJourneyByGroupID(groupID uint) (*entity.Journey, appError.CommonError) {
	args := m.Called(groupID)
	var arg0 *entity.Journey
	var arg1 appError.CommonError
	if args.Get(0) != nil {
		arg0 = args.Get(0).(*entity.Journey)
	}
	if args.Get(1) != nil {
		arg1 = args.Get(1).(appError.CommonError)
	}
	return arg0, arg1
}

func (m *journeyUseCaseMock) UpdateJourneyStatus(journeyID uint, newStatus enum.JourneyStatus) appError.CommonError {
	args := m.Called(journeyID, newStatus)
	var arg0 appError.CommonError
	if args.Get(0) != nil {
		arg0 = args.Get(0).(appError.CommonError)
	}
	return arg0
}

func (m *journeyUseCaseMock) DropoffJourney(groupID uint) appError.CommonError {
	args := m.Called(groupID)
	var arg0 appError.CommonError
	if args.Get(0) != nil {
		arg0 = args.Get(0).(appError.CommonError)
	}
	return arg0
}

func (m *journeyUseCaseMock) AssignCarToJourney(journeyID uint, carID uint) appError.CommonError {
	args := m.Called(journeyID, carID)
	var arg0 appError.CommonError
	if args.Get(0) != nil {
		arg0 = args.Get(0).(appError.CommonError)
	}
	return arg0
}

func (m *journeyUseCaseMock) DispatchNewPendingJourneyEvent(newPendingJourneyID uint) appError.CommonError {
	args := m.Called(newPendingJourneyID)
	var arg0 appError.CommonError
	if args.Get(0) != nil {
		arg0 = args.Get(0).(appError.CommonError)
	}
	return arg0
}

func (m *journeyUseCaseMock) LocateJourney(groupID uint) (*entity.Journey, appError.CommonError) {
	args := m.Called(groupID)
	var arg0 *entity.Journey
	var arg1 appError.CommonError
	if args.Get(0) != nil {
		arg0 = args.Get(0).(*entity.Journey)
	}
	if args.Get(1) != nil {
		arg1 = args.Get(1).(appError.CommonError)
	}
	return arg0, arg1
}

func (m *journeyUseCaseMock) TryToAssignAvailableCarToJourney(journeyID uint) appError.CommonError {
	args := m.Called(journeyID)
	var arg0 appError.CommonError
	if args.Get(0) != nil {
		arg0 = args.Get(0).(appError.CommonError)
	}
	return arg0
}

func (m *journeyUseCaseMock) TryToAssignWaitlistedJourneysToCar(carID uint) appError.CommonError {
	args := m.Called(carID)
	var arg0 appError.CommonError
	if args.Get(0) != nil {
		arg0 = args.Get(0).(appError.CommonError)
	}
	return arg0
}
