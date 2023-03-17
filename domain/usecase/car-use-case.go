package usecase

import (
	appError "gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/domain/app-error"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/domain/entity"
)

type CarUseCase interface {
	GetCarByID(carID uint) (*entity.Car, appError.CommonError)
	CreateCar(car *entity.Car) (carID uint, err appError.CommonError)
	DeleteCarsAndInsert(cars []*entity.Car) (carIDs []uint, err appError.CommonError)
	GetCarsWithSeatingCapacityForAtLeast(passengers uint) ([]*entity.Car, appError.CommonError)
	SetJourneyUseCase(journeyUseCase JourneyUseCase)
	DispatchNewCarAvailableEvent(newCarAvailableID uint) appError.CommonError
}
