package mapper

import (
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/domain/entity"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/interface-adapters/model"
)

func CarModelToEntity(carModel *model.Car) *entity.Car {
	return &entity.Car{
		ID:           carModel.ID,
		MaxSeats:     carModel.MaxSeats,
		CreationDate: carModel.CreationDate,
		UpdateDate:   carModel.UpdateDate,
		DeleteDate:   carModel.DeleteDate,
	}
}

func CarEntityToModel(carEntity *entity.Car) *model.Car {
	return &model.Car{
		ID:           carEntity.ID,
		MaxSeats:     carEntity.MaxSeats,
		CreationDate: carEntity.CreationDate,
		UpdateDate:   carEntity.UpdateDate,
		DeleteDate:   carEntity.DeleteDate,
	}
}

func ListCarEntityToModel(carEntities []*entity.Car) []*model.Car {
	var carModels []*model.Car

	for _, carEntity := range carEntities {
		carModels = append(carModels, CarEntityToModel(carEntity))
	}

	return carModels
}

func ListCarModelToEntity(carModels []*model.Car) []*entity.Car {
	var carEntities []*entity.Car

	for _, carModel := range carModels {
		carEntities = append(carEntities, CarModelToEntity(carModel))
	}

	return carEntities
}
