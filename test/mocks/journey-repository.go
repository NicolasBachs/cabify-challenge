package mocks

import (
	"github.com/stretchr/testify/mock"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/domain/entity"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/domain/enum"
)

type journeyRepositoryMock struct {
	mock.Mock
}

func NewJourneyRepositoryMock() *journeyRepositoryMock {
	return &journeyRepositoryMock{}
}

func (m *journeyRepositoryMock) GetJourneyByID(journeyID uint) (*entity.Journey, error) {
	args := m.Called(journeyID)
	var arg0 *entity.Journey
	if args.Get(0) != nil {
		arg0 = args.Get(0).(*entity.Journey)
	}
	var arg1 error
	if args.Error(1) != nil {
		arg1 = args.Error(1)
	}
	return arg0, arg1
}

func (m *journeyRepositoryMock) GetLastJourneyByGroupID(groupID uint) (*entity.Journey, error) {
	args := m.Called(groupID)
	var arg0 *entity.Journey
	if args.Get(0) != nil {
		arg0 = args.Get(0).(*entity.Journey)
	}
	var arg1 error
	if args.Error(1) != nil {
		arg1 = args.Error(1)
	}
	return arg0, arg1
}

func (m *journeyRepositoryMock) CreateJourney(journey *entity.Journey) (uint, error) {
	args := m.Called(journey)
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

func (m *journeyRepositoryMock) GetJourneysByStatus(status enum.JourneyStatus) ([]*entity.Journey, error) {
	args := m.Called(status)
	var arg0 []*entity.Journey
	if args.Get(0) != nil {
		arg0 = args.Get(0).([]*entity.Journey)
	}
	var arg1 error
	if args.Error(1) != nil {
		arg1 = args.Error(1)
	}
	return arg0, arg1
}

func (m *journeyRepositoryMock) CleanJourneys() error {
	args := m.Called()
	var arg0 error
	if args.Error(0) != nil {
		arg0 = args.Error(0)
	}
	return arg0
}

func (m *journeyRepositoryMock) UpdateJourneyStatus(journeyID uint, newStatus enum.JourneyStatus) error {
	args := m.Called(journeyID, newStatus)
	var arg0 error
	if args.Error(0) != nil {
		arg0 = args.Error(0)
	}
	return arg0
}

func (m *journeyRepositoryMock) GetPendingJourneysWherePassengersLessOrEqualTo(maxPassengers uint) ([]*entity.Journey, error) {
	args := m.Called(maxPassengers)
	var arg0 []*entity.Journey
	if args.Get(0) != nil {
		arg0 = args.Get(0).([]*entity.Journey)
	}
	var arg1 error
	if args.Error(1) != nil {
		arg1 = args.Error(1)
	}
	return arg0, arg1
}

func (m *journeyRepositoryMock) AssignCarToJourney(journeyID uint, carID uint) error {
	args := m.Called(journeyID, carID)
	var arg0 error
	if args.Error(0) != nil {
		arg0 = args.Error(0)
	}
	return arg0
}

func (m *journeyRepositoryMock) GetJourneysByCarID(carID uint) ([]*entity.Journey, error) {
	args := m.Called(carID)
	var arg0 []*entity.Journey
	if args.Get(0) != nil {
		arg0 = args.Get(0).([]*entity.Journey)
	}
	var arg1 error
	if args.Error(1) != nil {
		arg1 = args.Error(1)
	}
	return arg0, arg1
}
