package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/domain/entity"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/domain/usecase"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/interface-adapters/response"
)

type CarController struct {
	carUseCase usecase.CarUseCase
}

func NewCarController(carUseCase usecase.CarUseCase) *CarController {
	return &CarController{
		carUseCase: carUseCase,
	}
}

type createCarRequest struct {
	ID       uint `json:"id" binding:"required"`
	MaxSeats uint `json:"seats" binding:"required"`
}

// GetCarByID
//
//	@Summary		Get car by ID
//	@Description	Get car with the specified id
//	@Tags			Cars
//	@Param			id	path	int	true	"Car ID"
//	@Produce		json
//	@Success		200	{object}	entity.Car
//	@Failure		400	{object}	response.ErrorResponse
//	@Failure		404	{object}	response.ErrorResponse
//	@Failure		500	{object}	response.ErrorResponse
//	@Router			/cars/{id} [get]
func (c *CarController) GetCarByID(ctx *gin.Context) {
	carIDParam := ctx.Param("id")

	carID, err := strconv.Atoi(carIDParam)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	car, appErr := c.carUseCase.GetCarByID(uint(carID))
	if appErr != nil {
		ctx.AbortWithStatusJSON(appErr.StatusCode(), gin.H{"error": appErr.Error()})
		return
	}

	ctx.JSON(http.StatusOK, car)
}

// CreateCar
//
//	@Summary		Create car
//	@Description	Create a new car with the provided data
//	@Tags			Cars
//	@Accept			json
//	@Produce		json
//	@Param			request	body		controller.createCarRequest	true	"Car to create"
//	@Success		200		{object}	response.CreateCarResponse	"Car ID"
//	@Failure		400		{object}	response.ErrorResponse
//	@Failure		500		{object}	response.ErrorResponse
//	@Router			/cars [post]
func (c *CarController) CreateCar(ctx *gin.Context) {
	var req createCarRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var car *entity.Car

	car = &entity.Car{
		MaxSeats: req.MaxSeats,
	}

	carID, appErr := c.carUseCase.CreateCar(car)
	if appErr != nil {
		ctx.AbortWithStatusJSON(appErr.StatusCode(), gin.H{"error": appErr.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response.CreateCarResponse{ID: carID})
}

// DeleteCarsAndInsert
//
//	@Summary		Put cars
//	@Description	Delete all existing data about cars and journeys and create new cars with the provided data
//	@Tags			Cars
//	@Accept			json
//	@Produce		json
//	@Param			request	body		[]controller.createCarRequest	true	"Array of cars to create"
//	@Success		200		{object}	response.PutCarsResponse
//	@Failure		400		{object}	response.ErrorResponse
//	@Failure		500		{object}	response.ErrorResponse
//	@Router			/cars [put]
func (c *CarController) DeleteCarsAndInsert(ctx *gin.Context) {
	var req []createCarRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorResponse{Message: err.Error(), Status: http.StatusBadRequest})
		return
	}

	var cars []*entity.Car

	for _, element := range req {
		car := &entity.Car{
			ID:       element.ID,
			MaxSeats: element.MaxSeats,
		}

		cars = append(cars, car)
	}

	carIDs, appErr := c.carUseCase.DeleteCarsAndInsert(cars)
	if appErr != nil {
		ctx.AbortWithStatusJSON(appErr.StatusCode(), response.ErrorResponse{Message: appErr.Error(), Status: appErr.StatusCode()})
		return
	}

	ctx.JSON(http.StatusOK, response.PutCarsResponse{IDs: carIDs})
}
