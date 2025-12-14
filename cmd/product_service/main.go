package main

import (
	"log"

	productApp "github.com/leonardo849/product_supermarket/internal/application/product"
	"github.com/leonardo849/product_supermarket/internal/config"
	"github.com/leonardo849/product_supermarket/internal/infrastructure/http"
	productHandler "github.com/leonardo849/product_supermarket/internal/infrastructure/http/handlers/product"
	"github.com/leonardo849/product_supermarket/internal/infrastructure/persistence/postgres"
)

func main() {
	config := config.Load()
	db, err := postgres.NewConnection(config.DatabaseURL)
	if err != nil {
		log.Fatal(err.Error())
	}
	productRepo := postgres.NewProductRepository(db)
	stockRepo := postgres.NewStockRepository(db)
	uow := postgres.NewUnitOfWork(db)
	productUc := productApp.NewCreateProductUseCase(productRepo, stockRepo, uow)
	productHandler := productHandler.NewProductHandler(productUc)
	app := http.SetupApp(productHandler)
	app.Listen(":" + config.HTTPPort)
}
