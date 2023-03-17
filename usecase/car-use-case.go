package usecaseImpl

import (
	"context"

	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/config"
	appError "gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/domain/app-error"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/domain/entity"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/domain/event"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/domain/repository"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/domain/usecase"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/infrastructure/app"
	eventDispatcher "gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/infrastructure/event-dispatcher"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/infrastructure/mutex"
)

type carUseCase struct {
	eventDispatcher         eventDispatcher.EventDispatcher
	carRepository           repository.CarRepository
	journeyUseCase          usecase.JourneyUseCase
	distributedResourceSync mutex.DistributedResourceSync
}

func NewCarUseCase(
	eventDispatcher eventDispatcher.EventDispatcher,
	carRepository repository.CarRepository,
	journeyUseCase usecase.JourneyUseCase,
	distributedResourceSync mutex.DistributedResourceSync,
) usecase.CarUseCase {

	return &carUseCase{
		eventDispatcher:         eventDispatcher,
		carRepository:           carRepository,
		journeyUseCase:          journeyUseCase,
		distributedResourceSync: distributedResourceSync,
	}
}

var carUseCaseName = "CAR_USE_CASE"

func (c *carUseCase) SetJourneyUseCase(journeyUseCase usecase.JourneyUseCase) {
	c.journeyUseCase = journeyUseCase
}

func (c *carUseCase) GetCarByID(carID uint) (*entity.Car, appError.CommonError) {
	car, err := c.carRepository.GetCarByID(carID)

	if err != nil {
		app.Logger.Debug(carUseCaseName, "Error getting car with id: %d, error: %s", carID, err.Error())
		return nil, appError.NewInternalServerError(
			"Error getting car with id: %d", carID,
		)
	}

	if car == nil {
		return nil, appError.NewNotFoundError(
			"Car with ID '%d' not found", carID,
		)
	}

	journeys, err := c.journeyUseCase.GetJourneysByCarID(carID)
	if err != nil {
		app.Logger.Debug(carUseCaseName, "Error getting journeys associated to carID: %d, error: %s", carID, err.Error())
		return nil, appError.NewInternalServerError(
			"Error getting car with id: %d", carID,
		)
	}

	car.Journeys = journeys

	return car, nil
}

func (c *carUseCase) CreateCar(car *entity.Car) (uint, appError.CommonError) {
	carID, err := c.carRepository.CreateCar(car)

	if err != nil {
		return 0, appError.NewInternalServerError(
			"Error creating car",
		)
	}

	return carID, nil
}

func (c *carUseCase) DeleteCarsAndInsert(cars []*entity.Car) ([]uint, appError.CommonError) {
	appErr := c.journeyUseCase.CleanJourneys()

	if appErr != nil {
		app.Logger.Error(carUseCaseName, "Error cleaning journeys to delete cars and insert new ones, error: ", appErr.Error())
		return nil, appError.NewCommonError(
			appErr.StatusCode(),
			"Error on reset cars, error: %s",
			appErr.Error(),
		)
	}

	carIDs, err := c.carRepository.DeleteCarsAndInsert(cars)

	if err != nil {
		app.Logger.Error(carUseCaseName, "Error on delete cars, error: ", err.Error())
		return nil, appError.NewInternalServerError(
			"Error on reset cars",
		)
	}

	return carIDs, nil
}

func (c *carUseCase) GetCarsWithSeatingCapacityForAtLeast(passengers uint) ([]*entity.Car, appError.CommonError) {
	cars, err := c.carRepository.GetCarsWithAvailableSeatsGreaterOrEqualTo(passengers)

	if err != nil {
		return nil, appError.NewInternalServerError(
			"Error getting cars with seating capacity for at least %d passengers", passengers,
		)
	}

	return cars, nil
}

func (c *carUseCase) DispatchNewCarAvailableEvent(newCarAvailableID uint) appError.CommonError {
	ctx := context.Background()

	event := event.NewCarAvailableEvent{
		CarID: newCarAvailableID,
	}

	err := c.eventDispatcher.Dispatch(ctx, config.AppConfig.Kafka.Topics.NewCarAvailable, event)

	if err != nil {
		app.Logger.Error(
			journeyUseCaseName,
			"An error ocurred dispatching event to notify car '%d' is available, error: %s",
			newCarAvailableID,
			err.Error(),
		)
		return appError.NewInternalServerError(
			"An error ocurred dispatching event to notify car '%d' is available",
			newCarAvailableID,
		)
	}

	return nil
}
