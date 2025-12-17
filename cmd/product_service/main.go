package main

import (
	"log"

	productApp "github.com/leonardo849/product_supermarket/internal/application/product"
	userApp "github.com/leonardo849/product_supermarket/internal/application/user"
	"github.com/leonardo849/product_supermarket/internal/config"
	"github.com/leonardo849/product_supermarket/internal/infrastructure/http"
	productHandler "github.com/leonardo849/product_supermarket/internal/infrastructure/http/handlers/product"
	"github.com/leonardo849/product_supermarket/internal/infrastructure/messaging/rabbitmq"
	userConsumer "github.com/leonardo849/product_supermarket/internal/infrastructure/messaging/rabbitmq/consumer/users"
	"github.com/leonardo849/product_supermarket/internal/infrastructure/persistence/postgres"
	"github.com/rabbitmq/amqp091-go"
)



func main() {
	config := config.Load()
	db, err := postgres.NewConnection(config.DatabaseURL)
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
	runConsumers(channel, userRepo)
	productUc := productApp.NewCreateProductUseCase(productRepo, stockRepo, uow, userRepo, rabbitmq.NewPublisher(channel, "email_direct"))
	productHandler := productHandler.NewProductHandler(productUc)
	app := http.SetupApp(productHandler)
	app.Listen(":" + config.HTTPPort)
}

func runConsumers(ch *amqp091.Channel, userRepo *postgres.UserRepository) {
	userUc := userApp.NewCreateUserUseCase(userRepo)
	userCreatedProductConsumer := userConsumer.NewUserCreatedProductConsumer(ch, "queue_product_auth", userUc, rabbitmq.NewPublisher(ch, "product_auth_direct"))
	userCreatedProductConsumer.Start()
}