package consumer

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/domain/event"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/domain/usecase"
)

type CarsAvailableConsumer struct {
	journeyUseCase usecase.JourneyUseCase
}

func NewCarAvailableConsumer(journeyUseCase usecase.JourneyUseCase) *CarsAvailableConsumer {
	return &CarsAvailableConsumer{
		journeyUseCase: journeyUseCase,
	}
}

func (c *CarsAvailableConsumer) NewCarAvailable(ctx *gin.Context) {
	var event event.NewCarAvailableEvent

	if err := ctx.ShouldBindJSON(&event); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	appErr := c.journeyUseCase.TryToAssignWaitlistedJourneysToCar(event.CarID)
	if appErr != nil {
		ctx.AbortWithStatusJSON(appErr.StatusCode(), gin.H{"error": appErr.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "processed"})
}
