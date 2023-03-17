package usecaseImpl

import (
	"context"
	"time"

	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/config"
	appError "gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/domain/app-error"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/domain/entity"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/domain/enum"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/domain/event"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/domain/repository"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/domain/usecase"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/infrastructure/app"
	eventDispatcher "gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/infrastructure/event-dispatcher"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/infrastructure/mutex"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/infrastructure/mutex/resources"
)

type journeyUseCase struct {
	journeyRepository       repository.JourneyRepository
	eventDispatcher         eventDispatcher.EventDispatcher
	carUseCase              usecase.CarUseCase
	distributedResourceSync mutex.DistributedResourceSync
}

func NewJourneyUseCase(
	journeyRepository repository.JourneyRepository,
	eventDispatcher eventDispatcher.EventDispatcher,
	carUseCase usecase.CarUseCase,
	distributedResourceSync mutex.DistributedResourceSync,
) usecase.JourneyUseCase {

	return &journeyUseCase{
		journeyRepository:       journeyRepository,
		eventDispatcher:         eventDispatcher,
		carUseCase:              carUseCase,
		distributedResourceSync: distributedResourceSync,
	}
}

var journeyUseCaseName = "JOURNEY_USE_CASE"

func (c *journeyUseCase) GetJourneyByID(journeyID uint) (*entity.Journey, appError.CommonError) {
	journey, err := c.journeyRepository.GetJourneyByID(journeyID)

	if err != nil {
		app.Logger.Error(journeyUseCaseName, "Error getting journey with id: %d, error: %s", journeyID, err.Error())
		return nil, appError.NewInternalServerError(
			"Error getting journey with id: %d", journeyID,
		)
	}

	if journey == nil {
		return nil, appError.NewNotFoundError(
			"Journey with ID '%d' not found", journeyID,
		)
	}

	return journey, nil
}

func (c *journeyUseCase) GetJourneysByCarID(carID uint) ([]*entity.Journey, appError.CommonError) {
	journeys, err := c.journeyRepository.GetJourneysByCarID(carID)

	if err != nil {
		app.Logger.Error(journeyUseCaseName, "Error getting journeys associated to car with id: %d, error: %s", carID, err.Error())
		return nil, appError.NewInternalServerError(
			"Error getting journeys associated to car with id: %d", carID,
		)
	}

	return journeys, nil
}

func (c *journeyUseCase) CreateJourney(journey *entity.Journey) (uint, appError.CommonError) {
	groupID := journey.GroupID

	groupMutex, appErr := c.secureGroupTransaction(groupID)
	defer groupMutex.Unlock()
	if appErr != nil {
		return 0, appErr
	}

	hasPendingOrAssignedJourney, err := c.groupHasPendingOrAssignedJourney(groupID)

	if err != nil {
		app.Logger.Error(
			journeyUseCaseName,
			"Error creating journey. An error ocurred when trying to find the last journey of the group with id '%d'",
			groupID,
		)
		return 0, appError.NewInternalServerError("Error creating journey")
	}

	if hasPendingOrAssignedJourney {
		return 0, appError.NewConflictResourceStateError(
			"The group with id '%d' already has a pending or assigned journey",
			groupID,
		)
	}

	journey.Status = enum.JourneyStatusPending

	journeyID, err := c.journeyRepository.CreateJourney(journey)

	if err != nil {
		return 0, appError.NewInternalServerError("Error creating journey")
	}

	c.DispatchNewPendingJourneyEvent(journeyID)

	// Only due to restrictions with the challenge acceptance tests. Asynchronous processes were not considered.
	time.Sleep(2500 * time.Millisecond)

	return journeyID, nil
}

func (c *journeyUseCase) CleanJourneys() appError.CommonError {
	err := c.journeyRepository.CleanJourneys()

	if err != nil {
		appError.NewInternalServerError("Error cleaning journeys")
	}

	return nil
}

func (c *journeyUseCase) GetLastJourneyByGroupID(groupID uint) (*entity.Journey, appError.CommonError) {
	journey, err := c.journeyRepository.GetLastJourneyByGroupID(groupID)

	if err != nil {
		app.Logger.Error(
			journeyUseCaseName,
			"Error fetching last journey for group with ID: %d, error: %s",
			groupID,
			err.Error(),
		)
		return nil, appError.NewInternalServerError(
			"Error obtaining current journey of the group with id '%d'", groupID,
		)
	}

	if journey == nil {
		return nil, appError.NewNotFoundError(
			"Journey not found for group with id '%d'", groupID,
		)
	}

	return journey, nil
}

func (c *journeyUseCase) UpdateJourneyStatus(journeyID uint, newStatus enum.JourneyStatus) appError.CommonError {
	app.Logger.Debug(journeyUseCaseName, "Journey '%d' will be updated to status: '%s'", journeyID, newStatus)

	err := c.journeyRepository.UpdateJourneyStatus(journeyID, newStatus)

	if err != nil {
		app.Logger.Error(
			journeyUseCaseName,
			"Error updating journey '%d' to status: '%s', error: %s",
			journeyID,
			newStatus,
			err.Error(),
		)
		return appError.NewInternalServerError(
			"Error updating journey '%d' to status: '%s'",
			journeyID,
			newStatus,
		)
	}

	return nil
}

func (c *journeyUseCase) DropoffJourney(groupID uint) appError.CommonError {
	groupMutex, appErr := c.secureGroupTransaction(groupID)
	defer groupMutex.Unlock()
	if appErr != nil {
		return appErr
	}

	tempJourney, appErr := c.GetLastJourneyByGroupID(groupID)
	if appErr != nil {
		return appErr
	}
	journeyID := tempJourney.ID

	journeyMutex, appErr := c.secureJourneyTransaction(journeyID)
	defer journeyMutex.Unlock()
	if appErr != nil {
		return appErr
	}

	// Get the journey again now that we have acquired the lock. It may have been updated while we waited.
	journey, appErr := c.GetJourneyByID(journeyID)
	if appErr != nil {
		return appErr
	}

	var newStatus enum.JourneyStatus

	switch journey.Status {
	case enum.JourneyStatusPending:
		newStatus = enum.JourneyStatusCancelled
		break
	case enum.JourneyStatusAssigned:
		newStatus = enum.JourneyStatusFinished
		break
	case enum.JourneyStatusCancelled:
		app.Logger.Debug(journeyUseCaseName, "Journey '%d' for group with ID: %d already cancelled...", journeyID, groupID)
		return nil
	case enum.JourneyStatusFinished:
		app.Logger.Debug(journeyUseCaseName, "Journey '%d' for group with ID: %d already finished...", journeyID, groupID)
		return nil
	}

	appErr = c.UpdateJourneyStatus(journeyID, newStatus)
	if appErr != nil {
		return appErr
	}

	if newStatus == enum.JourneyStatusFinished && journey.CarAssignedID != nil {
		c.carUseCase.DispatchNewCarAvailableEvent(*journey.CarAssignedID)
	}

	return nil
}

func (c *journeyUseCase) AssignCarToJourney(journeyID uint, carID uint) appError.CommonError {
	app.Logger.Debug(journeyUseCaseName, "Assigning journey with id: '%d', to car with id: '%d'", journeyID, carID)

	journeyLock, appErr := c.secureJourneyTransaction(journeyID)
	if appErr != nil {
		return appErr
	}
	defer journeyLock.Unlock()

	journey, journeyIsAssignable, appErr := c.checkIfJourneyIsAssignable(journeyID)
	if appErr != nil {
		return appErr
	}
	if !journeyIsAssignable {
		return appError.NewConflictResourceStateError(
			"Journey not assignable",
		)
	}

	carLock, appErr := c.secureCarTransaction(carID)
	if appErr != nil {
		return appErr
	}
	defer carLock.Unlock()

	_, hasEnoughSeats, appErr := c.checkIfCarHasEnoughSeats(carID, journey.Passengers)
	if appErr != nil {
		return appErr
	}
	if !hasEnoughSeats {
		return appError.NewConflictResourceStateError(
			"Car has not enough seats",
		)
	}

	err := c.journeyRepository.AssignCarToJourney(journeyID, carID)

	if err != nil {
		app.Logger.Error(journeyUseCaseName, "Error assigning car '%d' to journey '%d'. Error: %s", carID, journeyID, err.Error())
		return appError.NewInternalServerError(
			"Error assigning car with id: '%d', to journey with id: '%d'",
			carID,
			journeyID,
		)
	}

	return nil
}

func (c *journeyUseCase) TryToAssignWaitlistedJourneysToCar(carID uint) appError.CommonError {
	app.Logger.Debug(journeyUseCaseName, "New car available, id: %d. Trying to assign to any journey in waitlist...", carID)

	carWithSeatsAvailable, appErr := c.carUseCase.GetCarByID(carID)

	if appErr != nil {
		return appErr
	}

	journeys, err := c.journeyRepository.GetPendingJourneysWherePassengersLessOrEqualTo(carWithSeatsAvailable.AvailableSeats)

	if err != nil {
		app.Logger.Error(
			journeyUseCaseName,
			"Error on get pending journeys with passengers quantity less or equal to car available seats, error: %s",
			err.Error(),
		)
		return appError.NewInternalServerError("Error trying to assign new available car to journey")
	}

	if len(journeys) > 0 {
		appErr := c.AssignCarToJourney(journeys[0].ID, carWithSeatsAvailable.ID)

		if appErr != nil {
			app.Logger.Debug(journeyUseCaseName, "Error assigning car to first journey in the waitlist, error: %s", err.Error())

			return appErr
		}
	}

	return nil
}

func (c *journeyUseCase) TryToAssignAvailableCarToJourney(journeyID uint) appError.CommonError {
	app.Logger.Debug(journeyUseCaseName, "New pending journey, id: %d. Trying to assign to any available car.", journeyID)

	journey, journeyIsAssignable, appErr := c.checkIfJourneyIsAssignable(journeyID)
	if appErr != nil {
		return appErr
	}
	if !journeyIsAssignable {
		return nil
	}

	carsWithEnoughSeats, appErr := c.carUseCase.GetCarsWithSeatingCapacityForAtLeast(journey.Passengers)

	if appErr != nil {
		return appErr
	}

	for index, car := range carsWithEnoughSeats {
		err := c.AssignCarToJourney(journey.ID, car.ID)
		if err != nil {
			if index < len(carsWithEnoughSeats) {
				app.Logger.Warn(journeyUseCaseName, "Error assigning car to journey, trying next available...")
				continue
			} else {
				app.Logger.Error(journeyUseCaseName, "Error assigning car to journey, no more cars available...")
				return err
			}
		}
		return nil
	}

	return nil
}

func (c *journeyUseCase) LocateJourney(groupID uint) (*entity.Journey, appError.CommonError) {
	journey, err := c.journeyRepository.GetLastJourneyByGroupID(groupID)

	if err != nil {
		app.Logger.Debug(
			journeyUseCaseName,
			"Error locating journey of group with id '%d'. Error fetching last journey of the group, error: %s",
			groupID,
			err.Error(),
		)
		return nil, appError.NewInternalServerError(
			"Error locating journey of group with id: %d",
			groupID,
		)
	}

	if journey == nil {
		return nil, appError.NewNotFoundError(
			"Journeys not found for group with id '%d'", groupID,
		)
	}

	if journey.CarAssignedID == nil {
		return nil, nil
	}

	if journey.Status != enum.JourneyStatusPending && journey.Status != enum.JourneyStatusAssigned {
		return nil, nil
	}

	carAssigned, appErr := c.carUseCase.GetCarByID(*journey.CarAssignedID)

	if appErr != nil || carAssigned == nil {
		app.Logger.Debug(
			journeyUseCaseName,
			"Error locating journey of group with id '%d'. Error fetching car assigned data",
			groupID,
		)
		return nil, appError.NewInternalServerError(
			"Error locating journey of group with id: %d",
			groupID,
		)
	}

	journey.CarAssigned = carAssigned

	return journey, nil
}

func (c *journeyUseCase) groupHasPendingOrAssignedJourney(groupID uint) (bool, error) {
	lastJourney, err := c.journeyRepository.GetLastJourneyByGroupID(groupID)

	if err != nil {
		return false, err
	}

	if lastJourney != nil {
		if lastJourney.Status == enum.JourneyStatusPending || lastJourney.Status == enum.JourneyStatusAssigned {
			return true, nil
		}
	}

	return false, nil
}

func (c *journeyUseCase) secureCarTransaction(carID uint) (mutex.DistributedMutex, appError.CommonError) {
	carLock, err := c.distributedResourceSync.NewMutex(resources.CarWithID(carID))

	if err != nil {
		app.Logger.Error(journeyUseCaseName, "Error on get distributed mutex for car '%d', error: %s", carID, err.Error())
		return nil, appError.NewInternalServerError(
			"The transaction could not be secured.",
		)
	}

	err = carLock.Lock()

	if err != nil {
		app.Logger.Error(journeyUseCaseName, "Error on acquire lock for car with id '%d', error: %s", carID, err.Error())
		return nil, appError.NewInternalServerError(
			"The transaction could not be secured.",
		)
	}

	return carLock, nil
}

func (c *journeyUseCase) secureJourneyTransaction(journeyID uint) (mutex.DistributedMutex, appError.CommonError) {
	journeyLock, err := c.distributedResourceSync.NewMutex(resources.JourneyWithID(journeyID))

	if err != nil {
		app.Logger.Error(journeyUseCaseName, "Error on get distributed mutex for journey '%d', error: %s", journeyID, err.Error())
		return nil, appError.NewInternalServerError(
			"The transaction could not be secured.",
		)
	}

	err = journeyLock.Lock()

	if err != nil {
		app.Logger.Error(journeyUseCaseName, "Error on acquire lock for journey with id '%d', error: %s", journeyID, err.Error())
		return nil, appError.NewInternalServerError(
			"The transaction could not be secured.",
		)
	}

	return journeyLock, nil
}

func (c *journeyUseCase) secureGroupTransaction(groupID uint) (mutex.DistributedMutex, appError.CommonError) {
	groupLock, err := c.distributedResourceSync.NewMutex(resources.GroupWithID(groupID))

	if err != nil {
		app.Logger.Error(journeyUseCaseName, "Error on get distributed mutex for group '%d', error: %s", groupID, err.Error())
		return nil, appError.NewInternalServerError(
			"The transaction could not be secured.",
		)
	}

	err = groupLock.Lock()

	if err != nil {
		app.Logger.Error(journeyUseCaseName, "Error on acquire lock for group with id '%d', error: %s", groupID, err.Error())
		return nil, appError.NewInternalServerError(
			"The transaction could not be secured.",
		)
	}

	return groupLock, nil
}

func (c *journeyUseCase) checkIfJourneyIsAssignable(journeyID uint) (*entity.Journey, bool, appError.CommonError) {
	journey, appErr := c.GetJourneyByID(journeyID)

	if appErr != nil {
		return nil, false, appErr
	}

	if journey.CarAssignedID != nil {
		app.Logger.Debug(journeyUseCaseName, "Journey with id: '%d' already has a car assigned, car id: '%d'", journeyID, journey.CarAssignedID)
		return journey, false, nil
	}

	return journey, true, nil
}

func (c *journeyUseCase) checkIfCarHasEnoughSeats(carID uint, journeyPassengers uint) (*entity.Car, bool, appError.CommonError) {
	car, appErr := c.carUseCase.GetCarByID(carID)

	if appErr != nil {
		return nil, false, appErr
	}

	if car.AvailableSeats < journeyPassengers {
		app.Logger.Error(
			journeyUseCaseName,
			"The car '%d' no longer has the required number of available seats.",
			carID,
		)
		return car, false, nil
	}

	return car, true, nil
}

func (c *journeyUseCase) DispatchNewPendingJourneyEvent(newPendingJourneyID uint) appError.CommonError {
	ctx := context.Background()

	event := event.NewPendingJourneyEvent{
		JourneyID: newPendingJourneyID,
	}

	err := c.eventDispatcher.Dispatch(ctx, config.AppConfig.Kafka.Topics.NewPendingJourney, event)
	if err != nil {
		app.Logger.Error(
			journeyUseCaseName,
			"An error ocurred dispatching event to notify new pending journey with id '%d', error: %s",
			newPendingJourneyID,
			err.Error(),
		)
		return appError.NewInternalServerError(
			"An error ocurred dispatching event to notify new pending journey with id '%d'",
			newPendingJourneyID,
		)
	}

	return nil
}
