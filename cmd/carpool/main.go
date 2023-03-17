package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/jmoiron/sqlx"
	goredislib "github.com/redis/go-redis/v9"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/config"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/domain/usecase"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/infrastructure/app"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/infrastructure/db"
	eventDispatcher "gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/infrastructure/event-dispatcher"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/infrastructure/event-dispatcher/kafka"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/infrastructure/mutex"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/infrastructure/router"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/interface-adapters/controller"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/interface-adapters/controller/consumer"
	repositoryImpl "gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/interface-adapters/repository"
	usecaseImpl "gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/usecase"
)

func main() {
	// Configs
	err := config.LoadConfig()

	if err != nil {
		log.Fatal("Error loading application configs, err: ", err)
	}

	// DB Handler
	dbHandler := getDBHandler()

	// Redsync
	redsyncInstance := getRedsyncInstance()

	// Logger
	app.Logger = app.NewAppLogger()

	// Event dispatcher
	eventDispatcher := getEventDispatcher()
	defer eventDispatcher.Close()

	// Router
	ginRouter := gin.Default()
	router := router.NewGinRouter(ginRouter)

	err = router.Setup(
		setupAppController(
			dbHandler.GetDB(), eventDispatcher, mutex.NewRedisSync(redsyncInstance),
		),
	)

	if err != nil {
		log.Fatal("Error setting up router, err: ", err)
	}

	port := config.AppConfig.App.Port

	router.Serve(port)
}

func setupAppController(
	db *sqlx.DB,
	eventDispatcher eventDispatcher.EventDispatcher,
	distributedResourceSync mutex.DistributedResourceSync,
) *controller.AppController {
	// Repositories
	journeyRepository := repositoryImpl.NewJourneyRepository(db)
	carRepository := repositoryImpl.NewCarRepository(db)

	// Usecases
	var journeyUseCase usecase.JourneyUseCase
	carUseCase := usecaseImpl.NewCarUseCase(eventDispatcher, carRepository, nil, distributedResourceSync)

	journeyUseCase = usecaseImpl.NewJourneyUseCase(journeyRepository, eventDispatcher, carUseCase, distributedResourceSync)
	carUseCase.SetJourneyUseCase(journeyUseCase)

	// Controllers
	healthCheckController := controller.NewHealthCheckController()
	carController := controller.NewCarController(carUseCase)
	journeyController := controller.NewJourneyController(journeyUseCase)

	// Consumers
	pendingJourneysConsumer := consumer.NewPendingJourneysConsumer(journeyUseCase)
	carsAvailableConsumer := consumer.NewCarAvailableConsumer(journeyUseCase)

	// AppController
	appController := controller.NewAppController(
		healthCheckController,
		carController,
		journeyController,
		pendingJourneysConsumer,
		carsAvailableConsumer,
	)

	return appController
}

func getDBHandler() *db.PostgresDBHandler {
	return db.NewPostgresDBHandler(config.AppConfig.Database.ConnectionString)
}

func getEventDispatcher() eventDispatcher.EventDispatcher {
	brokerUrl := fmt.Sprintf("%s:%s", config.AppConfig.Kafka.Host, config.AppConfig.Kafka.Port)
	eventDispatcher, err := kafka.NewKafkaEventDispatcher(brokerUrl)

	if err != nil {
		log.Panicf("Error initializing kafka event dispatcher, error: %s", err.Error())
	}

	return eventDispatcher
}

func getRedsyncInstance() *redsync.Redsync {
	redisAddr := fmt.Sprintf("%s:%s", config.AppConfig.Redis.Host, config.AppConfig.Redis.Port)

	client := goredislib.NewClient(&goredislib.Options{
		Addr: redisAddr,
	})
	pool := goredis.NewPool(client)
	return redsync.New(pool)
}
