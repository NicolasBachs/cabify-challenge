package repository

import (
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/domain/entity"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/domain/enum"
)

type JourneyRepository interface {
	GetJourneyByID(journeyID uint) (*entity.Journey, error)
	GetLastJourneyByGroupID(groupID uint) (*entity.Journey, error)
	CreateJourney(journey *entity.Journey) (journeyID uint, err error)
	GetJourneysByStatus(status enum.JourneyStatus) ([]*entity.Journey, error)
	CleanJourneys() (err error)
	UpdateJourneyStatus(journeyID uint, newStatus enum.JourneyStatus) error
	GetPendingJourneysWherePassengersLessOrEqualTo(maxPassengers uint) ([]*entity.Journey, error)
	AssignCarToJourney(journeyID uint, carID uint) error
	GetJourneysByCarID(carID uint) ([]*entity.Journey, error)
}
