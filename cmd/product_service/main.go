package main

import (
	"log"
	"time"

	productApp "github.com/leonardo849/product_supermarket/internal/application/product"
	userApp "github.com/leonardo849/product_supermarket/internal/application/user"
	"github.com/leonardo849/product_supermarket/internal/config"
	"github.com/leonardo849/product_supermarket/internal/infrastructure/http"
	productHandler "github.com/leonardo849/product_supermarket/internal/infrastructure/http/handlers/product"
	"github.com/leonardo849/product_supermarket/internal/infrastructure/messaging/rabbitmq"
	userConsumer "github.com/leonardo849/product_supermarket/internal/infrastructure/messaging/rabbitmq/consumer/users"
	"github.com/leonardo849/product_supermarket/internal/infrastructure/persistence/postgres"
	"github.com/leonardo849/product_supermarket/internal/infrastructure/redis"
	"github.com/rabbitmq/amqp091-go"
	domainUser "github.com/leonardo849/product_supermarket/internal/domain/user"
	domainError "github.com/leonardo849/product_supermarket/internal/domain/error"
)



func main() {
	config := config.Load()
	db, err := postgres.NewConnection(config.DatabaseURL)
	if err != nil {
		log.Fatal(err.Error())
	}
	client, err := redis.NewClient(config.RedisUri, config.RedisPassword, config.RedisDatabase)
	if err != nil {
		log.Fatal(err.Error())
	}
	productRepo := postgres.NewProductRepository(db)
	stockRepo := postgres.NewStockRepository(db)
	userRepo := postgres.NewUserRepository(db)
	uow := postgres.NewUnitOfWork(db)
	conn, err := rabbitmq.NewConnection(config.RabbitURL)
	if err != nil {
		log.Fatal(err.Error())
	}
	channel, err := conn.Channel()
	if err != nil {
		log.Fatal(err.Error())
	}
	userCache := redis.NewUserCache(client, time.Minute * 60)
	runConsumers(channel, userRepo, userCache, client)
	findUserByAuthId := userApp.NewCreateFindUserUseCaseByAuthId(userRepo, userCache)
	productUc := productApp.NewCreateProductUseCase(productRepo, stockRepo, uow, findUserByAuthId, rabbitmq.NewPublisher(channel, "email_direct"))
	productHandler := productHandler.NewProductHandler(productUc)
	app := http.SetupApp(productHandler)
	app.Listen(":" + config.HTTPPort)
}

func runConsumers(ch *amqp091.Channel, userRepo *postgres.UserRepository, userCache domainUser.Cache, rc *redis.Client) {
	var errorCache domainError.ErrorCache
	errorCache = redis.NewErrorCache(rc, 30*time.Minute)
	userUc := userApp.NewCreateUserUseCase(userRepo, userCache)
	userUcDeleteUserUseCase := userApp.NewDeleteUserUseCase(userRepo, userCache)
	userCreatedProductConsumer := userConsumer.NewUserCreatedProductConsumer(ch, "queue_product_auth", userUc, userUcDeleteUserUseCase,rabbitmq.NewPublisher(ch, "product_auth_direct"), errorCache)
	userCreatedProductConsumer.Start()
	userDeletedConsumer := userConsumer.NewDeletedUserProductConsumer(ch, "queue_product_auth_deleted",  userUcDeleteUserUseCase, errorCache)
	userDeletedConsumer.Start()
}