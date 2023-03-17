package usecase

import (
	appError "gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/domain/app-error"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/domain/entity"
)

type JourneyUseCase interface {
	GetJourneyByID(journeyID uint) (*entity.Journey, appError.CommonError)
	GetJourneysByCarID(carID uint) ([]*entity.Journey, appError.CommonError)
	CreateJourney(journey *entity.Journey) (journeyID uint, err appError.CommonError)
	CleanJourneys() (err appError.CommonError)
	DropoffJourney(groupID uint) (err appError.CommonError)
	AssignCarToJourney(journeyID uint, carID uint) (err appError.CommonError)
	TryToAssignWaitlistedJourneysToCar(carID uint) appError.CommonError
	TryToAssignAvailableCarToJourney(journeyID uint) appError.CommonError
	LocateJourney(groupID uint) (*entity.Journey, appError.CommonError)
	DispatchNewPendingJourneyEvent(newPendingJourneyID uint) appError.CommonError
}
