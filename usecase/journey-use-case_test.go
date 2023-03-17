package usecaseImpl

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/config"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/domain/entity"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/domain/enum"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/domain/event"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/infrastructure/app"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/infrastructure/mutex/resources"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/test/mocks"
)

type JourneyUseCaseTestSuite struct {
	suite.Suite
}

func (s *JourneyUseCaseTestSuite) BeforeTest(suiteName, testName string) {
	config.AppConfig = &config.Config{}
	config.AppConfig.App.Environment = "test"
	config.AppConfig.Kafka.Topics.NewPendingJourney = "topic-name-1"
	config.AppConfig.Kafka.Topics.NewCarAvailable = "topic-name-2"
	app.Logger = app.NewAppLogger()
}

func TestJourneyUseCaseSuite(t *testing.T) {
	suite.Run(t, new(JourneyUseCaseTestSuite))
}

func (s *JourneyUseCaseTestSuite) Test_GetJourneyByID_Success() {
	// Given
	mockJourneyRepository := mocks.NewJourneyRepositoryMock()
	mockCarUseCase := mocks.NewCarUseCaseMock()
	mockEventDispatcher := mocks.NewEventDispatcherMock()
	mockDistributedMutex := mocks.NewDistributedMutexMock()
	mockDistributedSyncResource := mocks.NewDistributedResourceSyncMock(mockDistributedMutex)

	journeyUseCase := NewJourneyUseCase(mockJourneyRepository, mockEventDispatcher, mockCarUseCase, mockDistributedSyncResource)
	journeyID := uint(1)
	carAssignedID := uint(1)

	expectedJourney := &entity.Journey{
		ID:            1,
		GroupID:       10,
		Passengers:    2,
		Status:        enum.JourneyStatusAssigned,
		CreationDate:  time.Now(),
		UpdateDate:    time.Now(),
		CarAssignedID: &carAssignedID,
	}

	// When
	mockJourneyRepository.On("GetJourneyByID", journeyID).Return(expectedJourney, nil)
	actualJourney, appErr := journeyUseCase.GetJourneyByID(journeyID)

	// Assert
	s.Nil(appErr)
	s.Equal(expectedJourney, actualJourney)
}

func (s *JourneyUseCaseTestSuite) Test_GetJourneysByCarID_Success() {
	// Given
	mockJourneyRepository := mocks.NewJourneyRepositoryMock()
	mockCarUseCase := mocks.NewCarUseCaseMock()
	mockEventDispatcher := mocks.NewEventDispatcherMock()
	mockDistributedMutex := mocks.NewDistributedMutexMock()
	mockDistributedSyncResource := mocks.NewDistributedResourceSyncMock(mockDistributedMutex)

	journeyUseCase := NewJourneyUseCase(mockJourneyRepository, mockEventDispatcher, mockCarUseCase, mockDistributedSyncResource)
	carID := uint(1)
	expectedJourneys := []*entity.Journey{
		{
			ID:            1,
			GroupID:       10,
			Passengers:    2,
			CarAssignedID: &carID,
			Status:        enum.JourneyStatusAssigned,
		},
	}

	// When
	mockJourneyRepository.On("GetJourneysByCarID", carID).Return(expectedJourneys, nil)
	actualJourneys, appErr := journeyUseCase.GetJourneysByCarID(carID)

	// Assert
	s.Nil(appErr)
	s.Equal(expectedJourneys, actualJourneys)
}

func (s *JourneyUseCaseTestSuite) Test_CreateJourney_Success() {
	// Given
	mockJourneyRepository := mocks.NewJourneyRepositoryMock()
	mockCarUseCase := mocks.NewCarUseCaseMock()
	mockEventDispatcher := mocks.NewEventDispatcherMock()
	mockDistributedMutex := mocks.NewDistributedMutexMock()
	mockDistributedSyncResource := mocks.NewDistributedResourceSyncMock(mockDistributedMutex)

	journeyUseCase := NewJourneyUseCase(mockJourneyRepository, mockEventDispatcher, mockCarUseCase, mockDistributedSyncResource)
	expectedJourneyID := uint(5)
	groupID := uint(1)
	lastCarAssignedID := uint(7)
	lastJourneyOfGroup := &entity.Journey{
		ID:            1,
		GroupID:       groupID,
		Passengers:    2,
		Status:        enum.JourneyStatusFinished,
		CarAssignedID: &lastCarAssignedID,
		CreationDate:  time.Now(),
		UpdateDate:    time.Now(),
	}
	event := event.NewPendingJourneyEvent{
		JourneyID: uint(expectedJourneyID),
	}
	newJourney := &entity.Journey{
		GroupID:    groupID,
		Passengers: 2,
	}

	// When
	mockDistributedSyncResource.On("NewMutex", resources.GroupWithID(groupID)).Return(mockDistributedMutex, nil)
	mockDistributedMutex.On("Lock").Return(nil)
	mockDistributedMutex.On("Unlock").Return(nil)
	mockJourneyRepository.On("GetLastJourneyByGroupID", groupID).Return(lastJourneyOfGroup, nil)
	mockJourneyRepository.On("CreateJourney", newJourney).Return(expectedJourneyID, nil)
	mockEventDispatcher.On(
		"Dispatch",
		context.Background(),
		config.AppConfig.Kafka.Topics.NewPendingJourney,
		event,
	).Return(nil)

	actualJourneyID, appErr := journeyUseCase.CreateJourney(newJourney)

	// Assert
	s.Nil(appErr)
	s.Equal(expectedJourneyID, actualJourneyID)
}

func (s *JourneyUseCaseTestSuite) Test_CleanJourneys_Success() {
	// Given
	mockJourneyRepository := mocks.NewJourneyRepositoryMock()
	mockCarUseCase := mocks.NewCarUseCaseMock()
	mockEventDispatcher := mocks.NewEventDispatcherMock()
	mockDistributedMutex := mocks.NewDistributedMutexMock()
	mockDistributedSyncResource := mocks.NewDistributedResourceSyncMock(mockDistributedMutex)

	journeyUseCase := NewJourneyUseCase(mockJourneyRepository, mockEventDispatcher, mockCarUseCase, mockDistributedSyncResource)

	// When
	mockJourneyRepository.On("CleanJourneys").Return(nil)
	appErr := journeyUseCase.CleanJourneys()

	// Assert
	s.Nil(appErr)
}

func (s *JourneyUseCaseTestSuite) Test_DropoffJourney_Success() {
	// Given
	mockJourneyRepository := mocks.NewJourneyRepositoryMock()
	mockCarUseCase := mocks.NewCarUseCaseMock()
	mockEventDispatcher := mocks.NewEventDispatcherMock()
	mockDistributedMutex := mocks.NewDistributedMutexMock()
	mockDistributedSyncResource := mocks.NewDistributedResourceSyncMock(mockDistributedMutex)

	journeyUseCase := NewJourneyUseCase(mockJourneyRepository, mockEventDispatcher, mockCarUseCase, mockDistributedSyncResource)
	groupID := uint(1)
	newStatus := enum.JourneyStatusFinished
	lastCarAssignedID := uint(7)
	lastJourneyID := uint(1)
	lastJourneyOfGroup := &entity.Journey{
		ID:            lastJourneyID,
		GroupID:       groupID,
		Passengers:    2,
		Status:        enum.JourneyStatusFinished,
		CarAssignedID: &lastCarAssignedID,
		CreationDate:  time.Now(),
		UpdateDate:    time.Now(),
	}

	// When
	mockDistributedSyncResource.On("NewMutex", resources.GroupWithID(groupID)).Return(mockDistributedMutex, nil)
	mockDistributedMutex.On("Lock").Return(nil)
	mockDistributedMutex.On("Unlock").Return(nil)
	mockJourneyRepository.On("GetLastJourneyByGroupID", groupID).Return(lastJourneyOfGroup, nil)
	mockDistributedSyncResource.On("NewMutex", resources.JourneyWithID(lastJourneyID)).Return(mockDistributedMutex, nil)
	mockJourneyRepository.On("GetJourneyByID", lastJourneyID).Return(lastJourneyOfGroup, nil) // assuming it didn't change
	mockJourneyRepository.On("UpdateJourneyStatus", lastJourneyID, newStatus).Return(nil)
	mockCarUseCase.On("DispatchNewCarAvailableEvent", lastCarAssignedID).Return(nil) // only called if new status is finished and has car assigned

	appErr := journeyUseCase.DropoffJourney(groupID)

	// Assert
	s.Nil(appErr)

}

func (s *JourneyUseCaseTestSuite) Test_AssignCarToJourney_Success() {
	// Given
	mockJourneyRepository := mocks.NewJourneyRepositoryMock()
	mockCarUseCase := mocks.NewCarUseCaseMock()
	mockEventDispatcher := mocks.NewEventDispatcherMock()
	mockDistributedMutex := mocks.NewDistributedMutexMock()
	mockDistributedSyncResource := mocks.NewDistributedResourceSyncMock(mockDistributedMutex)

	journeyUseCase := NewJourneyUseCase(mockJourneyRepository, mockEventDispatcher, mockCarUseCase, mockDistributedSyncResource)

	journeyID := uint(1)
	carID := uint(1)
	journey := &entity.Journey{
		ID:            1,
		GroupID:       10,
		Passengers:    2,
		Status:        enum.JourneyStatusFinished,
		CarAssignedID: nil,
		CreationDate:  time.Now(),
		UpdateDate:    time.Now(),
	}
	car := &entity.Car{
		ID:             1,
		MaxSeats:       4,
		AvailableSeats: 2,
		CreationDate:   time.Now(),
		UpdateDate:     time.Now(),
	}

	// When
	mockDistributedSyncResource.On("NewMutex", resources.JourneyWithID(journeyID)).Return(mockDistributedMutex, nil)
	mockDistributedMutex.On("Lock").Return(nil)
	mockDistributedMutex.On("Unlock").Return(nil)
	mockJourneyRepository.On("GetJourneyByID", journeyID).Return(journey, nil)
	mockDistributedSyncResource.On("NewMutex", resources.CarWithID(journeyID)).Return(mockDistributedMutex, nil)
	mockCarUseCase.On("GetCarByID", carID).Return(car, nil)
	mockJourneyRepository.On("AssignCarToJourney", journeyID, carID).Return(nil)

	appErr := journeyUseCase.AssignCarToJourney(journeyID, carID)

	// Assert
	s.Nil(appErr)
}

func (s *JourneyUseCaseTestSuite) Test_TryToAssignWaitlistedJourneysToCar_Success() {
	// Given
	mockJourneyRepository := mocks.NewJourneyRepositoryMock()
	mockCarUseCase := mocks.NewCarUseCaseMock()
	mockEventDispatcher := mocks.NewEventDispatcherMock()
	mockDistributedMutex := mocks.NewDistributedMutexMock()
	mockDistributedSyncResource := mocks.NewDistributedResourceSyncMock(mockDistributedMutex)

	journeyUseCase := NewJourneyUseCase(mockJourneyRepository, mockEventDispatcher, mockCarUseCase, mockDistributedSyncResource)
	carID := uint(1)
	journeyID := uint(1)
	journeys := []*entity.Journey{
		{
			ID:            journeyID,
			GroupID:       10,
			Passengers:    2,
			Status:        enum.JourneyStatusFinished,
			CarAssignedID: nil,
			CreationDate:  time.Now(),
			UpdateDate:    time.Now(),
		},
	}
	car := &entity.Car{
		ID:             carID,
		MaxSeats:       4,
		AvailableSeats: 2,
		CreationDate:   time.Now(),
		UpdateDate:     time.Now(),
	}

	// When
	mockCarUseCase.On("GetCarByID", carID).Return(car, nil)
	mockJourneyRepository.On("GetPendingJourneysWherePassengersLessOrEqualTo", car.AvailableSeats).Return(journeys, nil)
	mockDistributedSyncResource.On("NewMutex", resources.JourneyWithID(journeyID)).Return(mockDistributedMutex, nil)
	mockDistributedMutex.On("Lock").Return(nil)
	mockDistributedMutex.On("Unlock").Return(nil)
	mockJourneyRepository.On("GetJourneyByID", journeyID).Return(journeys[0], nil)
	mockDistributedSyncResource.On("NewMutex", resources.CarWithID(journeyID)).Return(mockDistributedMutex, nil)
	mockCarUseCase.On("GetCarByID", carID).Return(car, nil)
	mockJourneyRepository.On("AssignCarToJourney", journeyID, carID).Return(nil)

	appErr := journeyUseCase.TryToAssignWaitlistedJourneysToCar(carID)

	s.Nil(appErr)
}

func (s *JourneyUseCaseTestSuite) Test_TryToAssignAvailableCarToJourney_Success() {
	// Given
	mockJourneyRepository := mocks.NewJourneyRepositoryMock()
	mockCarUseCase := mocks.NewCarUseCaseMock()
	mockEventDispatcher := mocks.NewEventDispatcherMock()
	mockDistributedMutex := mocks.NewDistributedMutexMock()
	mockDistributedSyncResource := mocks.NewDistributedResourceSyncMock(mockDistributedMutex)

	journeyUseCase := NewJourneyUseCase(mockJourneyRepository, mockEventDispatcher, mockCarUseCase, mockDistributedSyncResource)
	carID := uint(1)
	journeyID := uint(1)
	journey := &entity.Journey{
		ID:            journeyID,
		GroupID:       10,
		Passengers:    2,
		Status:        enum.JourneyStatusFinished,
		CarAssignedID: nil,
		CreationDate:  time.Now(),
		UpdateDate:    time.Now(),
	}
	cars := []*entity.Car{
		{
			ID:             carID,
			MaxSeats:       4,
			AvailableSeats: 2,
			CreationDate:   time.Now(),
			UpdateDate:     time.Now(),
		},
	}

	// When
	mockJourneyRepository.On("GetJourneyByID", journeyID).Return(journey, nil)
	mockCarUseCase.On("GetCarsWithSeatingCapacityForAtLeast", journey.Passengers).Return(cars, nil)
	mockDistributedSyncResource.On("NewMutex", resources.JourneyWithID(journeyID)).Return(mockDistributedMutex, nil)
	mockDistributedMutex.On("Lock").Return(nil)
	mockDistributedMutex.On("Unlock").Return(nil)
	mockJourneyRepository.On("GetJourneyByID", journeyID).Return(journey, nil)
	mockDistributedSyncResource.On("NewMutex", resources.CarWithID(journeyID)).Return(mockDistributedMutex, nil)
	mockCarUseCase.On("GetCarByID", carID).Return(cars[0], nil)
	mockJourneyRepository.On("AssignCarToJourney", journeyID, carID).Return(nil)

	journeyUseCase.TryToAssignAvailableCarToJourney(journeyID)
}

func (s *JourneyUseCaseTestSuite) Test_LocateJourney_Success() {
	// Given
	mockJourneyRepository := mocks.NewJourneyRepositoryMock()
	mockCarUseCase := mocks.NewCarUseCaseMock()
	mockEventDispatcher := mocks.NewEventDispatcherMock()
	mockDistributedMutex := mocks.NewDistributedMutexMock()
	mockDistributedSyncResource := mocks.NewDistributedResourceSyncMock(mockDistributedMutex)

	journeyUseCase := NewJourneyUseCase(mockJourneyRepository, mockEventDispatcher, mockCarUseCase, mockDistributedSyncResource)
	groupID := uint(1)

	journeyID := uint(1)
	carAssignedID := uint(1)
	journey := &entity.Journey{
		ID:            journeyID,
		GroupID:       10,
		Passengers:    2,
		Status:        enum.JourneyStatusFinished,
		CarAssignedID: &carAssignedID,
		CreationDate:  time.Now(),
		UpdateDate:    time.Now(),
	}
	carAssigned := &entity.Car{
		ID:             carAssignedID,
		MaxSeats:       4,
		AvailableSeats: 2,
		CreationDate:   time.Now(),
		UpdateDate:     time.Now(),
	}

	// When
	mockJourneyRepository.On("GetLastJourneyByGroupID", groupID).Return(journey, nil)
	mockCarUseCase.On("GetCarByID", carAssignedID).Return(carAssigned, nil)
	journeyUseCase.LocateJourney(groupID)
}
