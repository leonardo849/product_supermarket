package bootstrap

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/leonardo849/product_supermarket/internal/config"
	"github.com/leonardo849/product_supermarket/internal/infrastructure/messaging/rabbitmq"
	"github.com/leonardo849/product_supermarket/internal/infrastructure/persistence/postgres"
	"github.com/leonardo849/product_supermarket/internal/infrastructure/redis"
	"github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
	productApp "github.com/leonardo849/product_supermarket/internal/application/product"
	userApp "github.com/leonardo849/product_supermarket/internal/application/user"
	productHandler "github.com/leonardo849/product_supermarket/internal/infrastructure/http/handlers/product"
	userHandler "github.com/leonardo849/product_supermarket/internal/infrastructure/http/handlers/user"
	httpServer "github.com/leonardo849/product_supermarket/internal/infrastructure/http"
	domainError "github.com/leonardo849/product_supermarket/internal/domain/error"
	domainUser "github.com/leonardo849/product_supermarket/internal/domain/user"
	userConsumer "github.com/leonardo849/product_supermarket/internal/infrastructure/messaging/rabbitmq/consumer/users"
)

type AppContainer struct {
	App         *fiber.App
	DB          *gorm.DB
	Redis       *redis.Client
	RabbitConn  *amqp091.Connection
	RabbitCh    *amqp091.Channel
}

func BuildApp(cfg *config.Config, startConsumers bool) (*AppContainer, error) {
	db, err := postgres.NewConnection(cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	redisClient, err := redis.NewClient(
		cfg.RedisUri,
		cfg.RedisPassword,
		cfg.RedisDatabase,
	)
	if err != nil {
		return nil, err
	}
	isEnabled := false
	if cfg.RabbitOn == "true" {
		isEnabled = true
	}
	var rabbitConn *amqp091.Connection
	var rabbitCh *amqp091.Channel

	if startConsumers {
		rabbitConn, rabbitCh = mustRabbit(cfg.RabbitURL)
	}

	productRepo := postgres.NewProductRepository(db)
	stockRepo := postgres.NewStockRepository(db)
	userRepo := postgres.NewUserRepository(db)
	uow := postgres.NewUnitOfWork(db)

	userCache := redis.NewUserCache(redisClient, 60*time.Minute)
	errorCache := redis.NewErrorCache(redisClient, 30*time.Minute)

	findUserByAuthId := userApp.NewCreateFindUserUseCaseByAuthId(
		userRepo,
		userCache,
	)

	productUC := productApp.NewCreateProductUseCase(
		productRepo,
		stockRepo,
		uow,
		findUserByAuthId,
		rabbitmq.NewPublisherIfRabbitIsEnabled(rabbitCh, "email_direct", isEnabled),
	)

	checkUserInErrorsUC := userApp.NewFindIfUserIsInErrors(
		errorCache,
		findUserByAuthId,
	)

	productHTTPHandler := productHandler.NewProductHandler(productUC)
	userHTTPHandler := userHandler.NewUserHandler(checkUserInErrorsUC)

	app := httpServer.SetupApp(productHTTPHandler, userHTTPHandler, errorCache)

	if startConsumers {
		startConsumersFn(
			rabbitCh, userRepo, userCache, errorCache, isEnabled,
		)
	}

	return &AppContainer{
		App:        app,
		DB:         db,
		Redis:      redisClient,
		RabbitConn: rabbitConn,
		RabbitCh:   rabbitCh,
	}, nil
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

func startConsumersFn(
	ch *amqp091.Channel,
	userRepo *postgres.UserRepository,
	userCache domainUser.Cache,
	errorCache domainError.ErrorCache,
	// rc *redis.Client,
	isEnabled bool,
) {
	userCreateUC := userApp.NewCreateUserUseCase(userRepo, userCache)
	userDeleteUC := userApp.NewDeleteUserUseCase(userRepo, userCache)
	
	userCreatedConsumer := userConsumer.NewUserCreatedProductConsumer(
		ch,
		"queue_product_auth",
		userCreateUC,
		userDeleteUC,
		rabbitmq.NewPublisherIfRabbitIsEnabled(ch, "product_auth_direct", isEnabled),
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