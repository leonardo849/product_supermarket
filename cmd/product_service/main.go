package main

import (
	"log"
	"github.com/leonardo849/product_supermarket/internal/config"
	"github.com/leonardo849/product_supermarket/internal/infrastructure/bootstrap"
)



func main() {
	cfg := config.Load()

	container, err := bootstrap.BuildApp(&cfg, cfg.PactMode == "false")
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(container.App.Listen(":" + cfg.HTTPPort))
}

// func BuildApp(cfg *config.Config, startConsumers bool) (*AppContainer, error) {
// 	db, err := postgres.NewConnection(cfg.DatabaseURL)
// 	if err != nil {
// 		return nil, err
// 	}

// 	redisClient, err := redis.NewClient(
// 		cfg.RedisUri,
// 		cfg.RedisPassword,
// 		cfg.RedisDatabase,
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var rabbitConn *amqp091.Connection
// 	var rabbitCh *amqp091.Channel

// 	if startConsumers {
// 		rabbitConn, rabbitCh = mustRabbit(cfg.RabbitURL)
// 	}

// 	productRepo := postgres.NewProductRepository(db)
// 	stockRepo := postgres.NewStockRepository(db)
// 	userRepo := postgres.NewUserRepository(db)
// 	uow := postgres.NewUnitOfWork(db)

// 	userCache := redis.NewUserCache(redisClient, 60*time.Minute)
// 	errorCache := redis.NewErrorCache(redisClient, 30*time.Minute)

// 	findUserByAuthId := userApp.NewCreateFindUserUseCaseByAuthId(
// 		userRepo,
// 		userCache,
// 	)

// 	productUC := productApp.NewCreateProductUseCase(
// 		productRepo,
// 		stockRepo,
// 		uow,
// 		findUserByAuthId,
// 		rabbitmq.NewPublisher(rabbitCh, "email_direct"),
// 	)

// 	checkUserInErrorsUC := userApp.NewFindIfUserIsInErrors(
// 		errorCache,
// 		findUserByAuthId,
// 	)

// 	productHTTPHandler := productHandler.NewProductHandler(productUC)
// 	userHTTPHandler := userHandler.NewUserHandler(checkUserInErrorsUC)

// 	app := httpServer.SetupApp(productHTTPHandler, userHTTPHandler, errorCache)

// 	if startConsumers {
// 		startConsumersFn(
// 			rabbitCh, userRepo, userCache, errorCache, redisClient,
// 		)
// 	}

// 	return &AppContainer{
// 		App:        app,
// 		DB:         db,
// 		Redis:      redisClient,
// 		RabbitConn: rabbitConn,
// 		RabbitCh:   rabbitCh,
// 	}, nil
// }


// func mustPostgres(databaseURL string) *gorm.DB {
// 	db, err := postgres.NewConnection(databaseURL)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return db
// }

// func mustRedis(cfg *config.Config) *redis.Client {
// 	client, err := redis.NewClient(
// 		cfg.RedisUri,
// 		cfg.RedisPassword,
// 		cfg.RedisDatabase,
// 	)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return client
// }

// func mustRabbit(rabbitURL string) (*amqp091.Connection, *amqp091.Channel) {
// 	conn, err := rabbitmq.NewConnection(rabbitURL)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	ch, err := conn.Channel()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	return conn, ch
// }

// func startConsumersFn(
// 	ch *amqp091.Channel,
// 	userRepo *postgres.UserRepository,
// 	userCache domainUser.Cache,
// 	errorCache domainError.ErrorCache,
// 	rc *redis.Client,
// ) {
// 	userCreateUC := userApp.NewCreateUserUseCase(userRepo, userCache)
// 	userDeleteUC := userApp.NewDeleteUserUseCase(userRepo, userCache)

// 	userCreatedConsumer := userConsumer.NewUserCreatedProductConsumer(
// 		ch,
// 		"queue_product_auth",
// 		userCreateUC,
// 		userDeleteUC,
// 		rabbitmq.NewPublisher(ch, "product_auth_direct"),
// 		errorCache,
// 	)

// 	userDeletedConsumer := userConsumer.NewDeletedUserProductConsumer(
// 		ch,
// 		"queue_product_auth_deleted",
// 		userDeleteUC,
// 		errorCache,
// 	)

// 	userCreatedConsumer.Start()
// 	userDeletedConsumer.Start()
// }
