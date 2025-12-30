package main

import (
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"

	productApp "github.com/leonardo849/product_supermarket/internal/application/product"
	userApp "github.com/leonardo849/product_supermarket/internal/application/user"
	"github.com/leonardo849/product_supermarket/internal/config"
	domainError "github.com/leonardo849/product_supermarket/internal/domain/error"
	domainUser "github.com/leonardo849/product_supermarket/internal/domain/user"
	httpServer "github.com/leonardo849/product_supermarket/internal/infrastructure/http"
	productHandler "github.com/leonardo849/product_supermarket/internal/infrastructure/http/handlers/product"
	userHandler "github.com/leonardo849/product_supermarket/internal/infrastructure/http/handlers/user"
	"github.com/leonardo849/product_supermarket/internal/infrastructure/messaging/rabbitmq"
	userConsumer "github.com/leonardo849/product_supermarket/internal/infrastructure/messaging/rabbitmq/consumer/users"
	"github.com/leonardo849/product_supermarket/internal/infrastructure/persistence/postgres"
	"github.com/leonardo849/product_supermarket/internal/infrastructure/redis"
)

func main() {
	cfg := config.Load()

	db := mustPostgres(cfg.DatabaseURL)
	redisClient := mustRedis(&cfg)
	rabbitConn, rabbitChannel := mustRabbit(cfg.RabbitURL)
	defer rabbitConn.Close()

	productRepo := postgres.NewProductRepository(db)
	stockRepo := postgres.NewStockRepository(db)
	userRepo := postgres.NewUserRepository(db)
	uow := postgres.NewUnitOfWork(db)

	userCache := redis.NewUserCache(redisClient, 60*time.Minute)
	errorCache := redis.NewErrorCache(redisClient, 30*time.Minute)

	startConsumers(
		rabbitChannel,
		userRepo,
		userCache,
		errorCache,
		redisClient,
	)

	findUserByAuthId := userApp.NewCreateFindUserUseCaseByAuthId(
		userRepo,
		userCache,
	)

	productUC := productApp.NewCreateProductUseCase(
		productRepo,
		stockRepo,
		uow,
		findUserByAuthId,
		rabbitmq.NewPublisher(rabbitChannel, "email_direct"),
	)

	checkUserInErrorsUC := userApp.NewFindIfUserIsInErrors(
		errorCache,
		findUserByAuthId,
	)

	productHTTPHandler := productHandler.NewProductHandler(productUC)
	userHTTPHandler := userHandler.NewUserHandler(checkUserInErrorsUC)

	app := httpServer.SetupApp(productHTTPHandler, userHTTPHandler, errorCache)

	log.Fatal(app.Listen(":" + cfg.HTTPPort))
}

func mustPostgres(databaseURL string) *gorm.DB {
	db, err := postgres.NewConnection(databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func mustRedis(cfg *config.Config) *redis.Client {
	client, err := redis.NewClient(
		cfg.RedisUri,
		cfg.RedisPassword,
		cfg.RedisDatabase,
	)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func mustRabbit(rabbitURL string) (*amqp091.Connection, *amqp091.Channel) {
	conn, err := rabbitmq.NewConnection(rabbitURL)
	if err != nil {
		log.Fatal(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	return conn, ch
}

func startConsumers(
	ch *amqp091.Channel,
	userRepo *postgres.UserRepository,
	userCache domainUser.Cache,
	errorCache domainError.ErrorCache,
	rc *redis.Client,
) {
	userCreateUC := userApp.NewCreateUserUseCase(userRepo, userCache)
	userDeleteUC := userApp.NewDeleteUserUseCase(userRepo, userCache)

	userCreatedConsumer := userConsumer.NewUserCreatedProductConsumer(
		ch,
		"queue_product_auth",
		userCreateUC,
		userDeleteUC,
		rabbitmq.NewPublisher(ch, "product_auth_direct"),
		errorCache,
	)

	userDeletedConsumer := userConsumer.NewDeletedUserProductConsumer(
		ch,
		"queue_product_auth_deleted",
		userDeleteUC,
		errorCache,
	)

	userCreatedConsumer.Start()
	userDeletedConsumer.Start()
}
