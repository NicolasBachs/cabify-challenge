package consumer

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/domain/event"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/domain/usecase"
)

type PendingJourneysConsumer struct {
	journeyUseCase usecase.JourneyUseCase
}

func NewPendingJourneysConsumer(journeyUseCase usecase.JourneyUseCase) *PendingJourneysConsumer {
	return &PendingJourneysConsumer{
		journeyUseCase: journeyUseCase,
	}
}

func (c *PendingJourneysConsumer) NewPendingJourneys(ctx *gin.Context) {
	var event event.NewPendingJourneyEvent

	if err := ctx.ShouldBindJSON(&event); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	appErr := c.journeyUseCase.TryToAssignAvailableCarToJourney(event.JourneyID)
	if appErr != nil {
		ctx.AbortWithStatusJSON(appErr.StatusCode(), gin.H{"error": appErr.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "processed"})
}
