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
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/test/mocks"
)

type CarUseCaseTestSuite struct {
	suite.Suite
}

func (s *CarUseCaseTestSuite) BeforeTest(suiteName, testName string) {
	config.AppConfig = &config.Config{}
	config.AppConfig.App.Environment = "test"
	config.AppConfig.Kafka.Topics.NewPendingJourney = "topic-name-1"
	config.AppConfig.Kafka.Topics.NewCarAvailable = "topic-name-2"
	app.Logger = app.NewAppLogger()
}

func TestCarUseCaseSuite(t *testing.T) {
	suite.Run(t, new(CarUseCaseTestSuite))
}

func (s *CarUseCaseTestSuite) Test_NewCarUseCase_Success() {
	// Given
	mockCarRepository := mocks.NewCarRepositoryMock()
	mockJourneyUseCase := mocks.NewJourneyUseCaseMock()
	mockEventDispatcher := mocks.NewEventDispatcherMock()
	mockDistributedMutex := mocks.NewDistributedMutexMock()
	mockDistributedSyncResource := mocks.NewDistributedResourceSyncMock(mockDistributedMutex)

	// Assert
	s.NotNil(NewCarUseCase(mockEventDispatcher, mockCarRepository, mockJourneyUseCase, mockDistributedSyncResource))
}

func (s *CarUseCaseTestSuite) Test_GetCarByID_Success() {
	// Given
	mockCarRepository := mocks.NewCarRepositoryMock()
	mockJourneyUseCase := mocks.NewJourneyUseCaseMock()
	mockEventDispatcher := mocks.NewEventDispatcherMock()
	mockDistributedMutex := mocks.NewDistributedMutexMock()
	mockDistributedSyncResource := mocks.NewDistributedResourceSyncMock(mockDistributedMutex)

	carUseCase := NewCarUseCase(mockEventDispatcher, mockCarRepository, mockJourneyUseCase, mockDistributedSyncResource)
	carID := uint(1)
	expectedJourneys := []*entity.Journey{
		{
			ID:            1,
			Passengers:    2,
			Status:        enum.JourneyStatusAssigned,
			CarAssignedID: &carID,
			CreationDate:  time.Now(),
			UpdateDate:    time.Now(),
		},
	}
	expectedCar := &entity.Car{
		ID:             1,
		MaxSeats:       4,
		AvailableSeats: 2,
		Journeys:       expectedJourneys,
	}

	// When
	mockCarRepository.On("GetCarByID", carID).Return(expectedCar, nil)
	mockJourneyUseCase.On("GetJourneysByCarID", carID).Return(expectedJourneys, nil)

	actualCar, appErr := carUseCase.GetCarByID(uint(carID))

	// Assert
	s.Nil(appErr)
	s.Equal(expectedCar, actualCar)
}

func (s *CarUseCaseTestSuite) Test_CreateCar_Success() {
	// Given
	mockCarRepository := mocks.NewCarRepositoryMock()
	mockJourneyUseCase := mocks.NewJourneyUseCaseMock()
	mockEventDispatcher := mocks.NewEventDispatcherMock()
	mockDistributedMutex := mocks.NewDistributedMutexMock()
	mockDistributedSyncResource := mocks.NewDistributedResourceSyncMock(mockDistributedMutex)

	carUseCase := NewCarUseCase(mockEventDispatcher, mockCarRepository, mockJourneyUseCase, mockDistributedSyncResource)
	expectedCarID := uint(1)
	car := &entity.Car{
		ID:       expectedCarID,
		MaxSeats: 4,
	}

	// When
	mockCarRepository.On("CreateCar", car).Return(expectedCarID, nil)
	actualCarID, err := carUseCase.CreateCar(car)

	// Assert
	s.Nil(err)
	s.Equal(expectedCarID, actualCarID)
}

func (s *CarUseCaseTestSuite) Test_DeleteCarsAndInsert_Success() {
	// Given
	mockCarRepository := mocks.NewCarRepositoryMock()
	mockJourneyUseCase := mocks.NewJourneyUseCaseMock()
	mockEventDispatcher := mocks.NewEventDispatcherMock()
	mockDistributedMutex := mocks.NewDistributedMutexMock()
	mockDistributedSyncResource := mocks.NewDistributedResourceSyncMock(mockDistributedMutex)

	carUseCase := NewCarUseCase(mockEventDispatcher, mockCarRepository, mockJourneyUseCase, mockDistributedSyncResource)
	expectedCarIDs := []uint{1, 2}
	cars := []*entity.Car{
		{
			ID:       1,
			MaxSeats: 4,
		},
		{
			ID:       2,
			MaxSeats: 6,
		},
	}

	// When
	mockJourneyUseCase.On("CleanJourneys").Return(nil)
	mockCarRepository.On("DeleteCarsAndInsert", cars).Return(expectedCarIDs, nil)

	actualCarIDs, err := carUseCase.DeleteCarsAndInsert(cars)

	// Assert
	s.Nil(err)
	s.Equal(expectedCarIDs, actualCarIDs)
}

func (s *CarUseCaseTestSuite) Test_GetCarsWithSeatingCapacityForAtLeast_Success() {
	// Given
	mockCarRepository := mocks.NewCarRepositoryMock()
	mockJourneyUseCase := mocks.NewJourneyUseCaseMock()
	mockEventDispatcher := mocks.NewEventDispatcherMock()
	mockDistributedMutex := mocks.NewDistributedMutexMock()
	mockDistributedSyncResource := mocks.NewDistributedResourceSyncMock(mockDistributedMutex)

	passengers := uint(2)
	carUseCase := NewCarUseCase(mockEventDispatcher, mockCarRepository, mockJourneyUseCase, mockDistributedSyncResource)

	expectedCars := []*entity.Car{
		{
			ID:             1,
			MaxSeats:       4,
			AvailableSeats: 2,
			CreationDate:   time.Now(),
			UpdateDate:     time.Now(),
		},
		{
			ID:             2,
			MaxSeats:       6,
			AvailableSeats: 2,
			CreationDate:   time.Now(),
			UpdateDate:     time.Now(),
		},
	}

	// When
	mockCarRepository.On("GetCarsWithAvailableSeatsGreaterOrEqualTo", passengers).Return(expectedCars, nil)
	actualCars, err := carUseCase.GetCarsWithSeatingCapacityForAtLeast(passengers)

	// Assert
	s.Nil(err)
	s.Equal(expectedCars, actualCars)
}

func (s *CarUseCaseTestSuite) Test_DispatchNewCarAvailableEvent_Success() {
	// Given
	mockCarRepository := mocks.NewCarRepositoryMock()
	mockJourneyUseCase := mocks.NewJourneyUseCaseMock()
	mockEventDispatcher := mocks.NewEventDispatcherMock()
	mockDistributedMutex := mocks.NewDistributedMutexMock()
	mockDistributedSyncResource := mocks.NewDistributedResourceSyncMock(mockDistributedMutex)

	newCarAvailableID := uint(1)
	event := event.NewCarAvailableEvent{
		CarID: newCarAvailableID,
	}
	carUseCase := NewCarUseCase(mockEventDispatcher, mockCarRepository, mockJourneyUseCase, mockDistributedSyncResource)

	// When
	mockEventDispatcher.On(
		"Dispatch",
		context.Background(),
		config.AppConfig.Kafka.Topics.NewCarAvailable,
		event,
	).Return(nil)

	appErr := carUseCase.DispatchNewCarAvailableEvent(newCarAvailableID)

	// Assert
	s.Nil(appErr)
}
