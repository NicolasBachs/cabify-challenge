package repository

import "gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/domain/entity"

type CarRepository interface {
	GetCarByID(carID uint) (*entity.Car, error)
	CreateCar(*entity.Car) (carID uint, err error)
	DeleteCarsAndInsert([]*entity.Car) (carIDs []uint, err error)
	GetCarsWithAvailableSeatsGreaterOrEqualTo(passengers uint) ([]*entity.Car, error)
}
