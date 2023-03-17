package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/docs"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/infrastructure/app"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/interface-adapters/controller"
)

type ginRouter struct {
	router *gin.Engine
}

func NewGinRouter(router *gin.Engine) Router {
	return &ginRouter{
		router: router,
	}
}

var ginRouterName = "GIN_ROUTER"

func (r *ginRouter) Setup(appController *controller.AppController) error {
	app.Logger.Debug(ginRouterName, "Setting up gin Router...")

	gin.SetMode(gin.ReleaseMode)

	r.router.HandleMethodNotAllowed = true

	// Default
	r.router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "404 Page not found"})
	})

	r.router.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "405 Method not allowed"})
	})

	// Status
	r.router.GET("/status", appController.HealthCheckController.GetStatus)

	// Cars
	r.router.POST("/cars", MiddlewareApplicationJsonContentType, appController.CarController.CreateCar)
	r.router.PUT("/cars", MiddlewareApplicationJsonContentType, appController.CarController.DeleteCarsAndInsert)
	r.router.GET("/cars/:id", appController.CarController.GetCarByID)

	// Journeys
	r.router.POST("/journey", MiddlewareApplicationJsonContentType, appController.JourneyController.CreateJourney)
	r.router.GET("/journey/:id", appController.JourneyController.GetJourneyByID)
	r.router.POST("/dropoff", MiddlewareWwwFormUrlEncodedContentType, appController.JourneyController.DropoffJourney)
	r.router.POST("/locate", MiddlewareWwwFormUrlEncodedContentType, appController.JourneyController.LocateGroup)

	// Consumers
	r.router.POST("/consumers/new-pending-journey", MiddlewareApplicationJsonContentType, appController.PendingJourneysConsumer.NewPendingJourneys)
	r.router.POST("/consumers/new-car-available", MiddlewareApplicationJsonContentType, appController.CarsAvailableConsumer.NewCarAvailable)

	r.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	app.Logger.Debug(ginRouterName, "Gin Router Setup finished successfully...")
	return nil
}

func (r *ginRouter) Serve(port string) {
	app.Logger.Info(ginRouterName, "Gin HTTP server running on port %v\n", port)

	err := r.router.Run(":" + port)

	if err != nil {
		panic("Failed to start Gin server, error: " + err.Error())
	}
}

func MiddlewareApplicationJsonContentType(c *gin.Context) {
	contentType := c.Request.Header.Get("Content-Type")
	if contentType != "application/json" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "El Content-Type de la solicitud debe ser application/json",
		})
	}
	c.Next()
}

func MiddlewareWwwFormUrlEncodedContentType(c *gin.Context) {
	contentType := c.Request.Header.Get("Content-Type")
	if contentType != "application/x-www-form-urlencoded" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "El Content-Type de la solicitud debe ser application/x-www-form-urlencoded",
		})
	}
	c.Next()
}
