package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/domain/entity"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/domain/usecase"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/interface-adapters/response"
)

type JourneyController struct {
	journeyUseCase usecase.JourneyUseCase
}

func NewJourneyController(journeyUseCase usecase.JourneyUseCase) *JourneyController {
	return &JourneyController{
		journeyUseCase: journeyUseCase,
	}
}

type createJourneyRequest struct {
	GroupID uint `json:"id" binding:"required"`
	People  uint `json:"people" binding:"required"`
}

// GetJourneyByID
//
//	@Summary		Get journey by ID
//	@Description	Get journey with the specified id
//	@Tags			Journeys
//	@Param			id	path	int	true	"Journey ID"
//	@Produce		json
//	@Success		200	{object}	entity.Journey
//	@Failure		400	{object}	response.ErrorResponse
//	@Failure		404	{object}	response.ErrorResponse
//	@Failure		500	{object}	response.ErrorResponse
//	@Router			/journey/{id} [get]
func (c *JourneyController) GetJourneyByID(ctx *gin.Context) {
	journeyIDParam := ctx.Param("id")

	journeyID, err := strconv.Atoi(journeyIDParam)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorResponse{Message: err.Error(), Status: http.StatusBadRequest})
		return
	}

	journey, appErr := c.journeyUseCase.GetJourneyByID(uint(journeyID))
	if appErr != nil {
		ctx.AbortWithStatusJSON(appErr.StatusCode(), response.ErrorResponse{Message: appErr.Error(), Status: appErr.StatusCode()})
		return
	}

	ctx.JSON(http.StatusOK, journey)
}

// CreateJourney
//
//	@Summary		Create journey
//	@Description	Create a new journey with the provided data
//	@Tags			Journeys
//	@Accept			json
//	@Produce		json
//	@Param			request	body	controller.createJourneyRequest	true	"Journey to create"
//	@Success		202
//	@Failure		400	{object}	response.ErrorResponse
//	@Failure		500	{object}	response.ErrorResponse
//	@Router			/journey [post]
func (c *JourneyController) CreateJourney(ctx *gin.Context) {
	var req createJourneyRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorResponse{Message: err.Error(), Status: http.StatusBadRequest})
		return
	}

	var journey *entity.Journey

	journey = &entity.Journey{
		Passengers: req.People,
		GroupID:    req.GroupID,
	}

	_, appErr := c.journeyUseCase.CreateJourney(journey)
	if appErr != nil {
		ctx.AbortWithStatusJSON(appErr.StatusCode(), response.ErrorResponse{Message: appErr.Error(), Status: appErr.StatusCode()})
		return
	}

	ctx.Status(http.StatusAccepted)
}

// GetJourneyByID
//
//	@Summary		Drop off car
//	@Description	Finish or cancell journey of the group with the specified ID
//	@Tags			Journeys
//	@Accept			x-www-form-urlencoded
//	@Param			ID	formData	int	true	"Group ID"
//	@Produce		json
//	@Success		204
//	@Failure		400	{object}	response.ErrorResponse
//	@Failure		404	{object}	response.ErrorResponse
//	@Failure		500	{object}	response.ErrorResponse
//	@Router			/dropoff [post]
func (c *JourneyController) DropoffJourney(ctx *gin.Context) {
	journeyID, err := strconv.Atoi(ctx.PostForm("ID"))

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorResponse{Message: err.Error(), Status: http.StatusBadRequest})
		return
	}

	appErr := c.journeyUseCase.DropoffJourney(uint(journeyID))
	if appErr != nil {
		ctx.AbortWithStatusJSON(appErr.StatusCode(), response.ErrorResponse{Message: appErr.Error(), Status: appErr.StatusCode()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// GetJourneyByID
//
//	@Summary		Locate car
//	@Description	Get car assigned to the group with the specified ID
//	@Tags			Journeys
//	@Accept			x-www-form-urlencoded
//	@Param			ID	formData	int	true	"Group ID"
//	@Produce		json
//	@Success		200	{object}	entity.Car
//	@Success		204
//	@Failure		400	{object}	response.ErrorResponse
//	@Failure		404
//	@Failure		500	{object}	response.ErrorResponse
//	@Router			/locate [post]
func (c *JourneyController) LocateGroup(ctx *gin.Context) {
	groupID, err := strconv.Atoi(ctx.PostForm("ID"))

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorResponse{Message: err.Error(), Status: http.StatusBadRequest})
		return
	}

	journey, appErr := c.journeyUseCase.LocateJourney(uint(groupID))

	if appErr != nil {
		if appErr.StatusCode() == http.StatusNotFound {
			ctx.Status(http.StatusNotFound)
			return
		}
		ctx.AbortWithStatusJSON(appErr.StatusCode(), response.ErrorResponse{Message: appErr.Error(), Status: appErr.StatusCode()})
		return
	}

	if journey == nil {
		ctx.Status(http.StatusNoContent)
		return
	}

	ctx.JSON(http.StatusOK, journey.CarAssigned)
}
