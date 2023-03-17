package controller

import "gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/interface-adapters/controller/consumer"

type AppController struct {
	HealthCheckController   *HealthCheckController
	CarController           *CarController
	JourneyController       *JourneyController
	PendingJourneysConsumer *consumer.PendingJourneysConsumer
	CarsAvailableConsumer   *consumer.CarsAvailableConsumer
}

func NewAppController(
	HealthCheckController *HealthCheckController,
	CarController *CarController,
	JourneyController *JourneyController,
	PendingJourneysConsumer *consumer.PendingJourneysConsumer,
	CarsAvailableConsumer *consumer.CarsAvailableConsumer,
) *AppController {
	return &AppController{
		HealthCheckController:   HealthCheckController,
		CarController:           CarController,
		JourneyController:       JourneyController,
		PendingJourneysConsumer: PendingJourneysConsumer,
		CarsAvailableConsumer:   CarsAvailableConsumer,
	}
}
