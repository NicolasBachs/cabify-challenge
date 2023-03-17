package mapper

import (
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/domain/entity"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/interface-adapters/model"
)

func JourneyModelToEntity(journeyModel *model.Journey) *entity.Journey {
	var carAssigned *entity.Car

	if journeyModel.CarAssigned != nil {
		carAssigned = CarModelToEntity(journeyModel.CarAssigned)
	}

	return &entity.Journey{
		ID:            journeyModel.ID,
		GroupID:       journeyModel.GroupID,
		Passengers:    journeyModel.Passengers,
		Status:        journeyModel.Status,
		CreationDate:  journeyModel.CreationDate,
		UpdateDate:    journeyModel.UpdateDate,
		DeleteDate:    journeyModel.DeleteDate,
		CarAssignedID: journeyModel.CarAssignedID,
		CarAssigned:   carAssigned,
	}
}

func JourneyEntityToModel(journeyEntity *entity.Journey) *model.Journey {
	var carAssigned *model.Car

	if journeyEntity.CarAssigned != nil {
		carAssigned = CarEntityToModel(journeyEntity.CarAssigned)
	}

	return &model.Journey{
		ID:            journeyEntity.ID,
		GroupID:       journeyEntity.GroupID,
		Passengers:    journeyEntity.Passengers,
		Status:        journeyEntity.Status,
		CreationDate:  journeyEntity.CreationDate,
		UpdateDate:    journeyEntity.UpdateDate,
		DeleteDate:    journeyEntity.DeleteDate,
		CarAssignedID: journeyEntity.CarAssignedID,
		CarAssigned:   carAssigned,
	}
}

func ListJourneyEntityToModel(journeyEntities []*entity.Journey) []*model.Journey {
	var journeyModels []*model.Journey

	for _, journeyEntity := range journeyEntities {
		journeyModels = append(journeyModels, JourneyEntityToModel(journeyEntity))
	}

	return journeyModels
}

func ListJourneyModelToEntity(journeyModels []*model.Journey) []*entity.Journey {
	var journeyEntities []*entity.Journey

	for _, journeyModel := range journeyModels {
		journeyEntities = append(journeyEntities, JourneyModelToEntity(journeyModel))
	}

	return journeyEntities
}
