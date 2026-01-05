

package product_test

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/leonardo849/product_supermarket/internal/infrastructure/persistence/postgres"
	"github.com/leonardo849/product_supermarket/internal/test_utils"
)

var app *fiber.App

func createUserTest() {
	db := testutils.GetDB()
	db.Create(&postgres.UserModel{
		ID: uuid.New(),
		AuthId: "69558ca84f914ff89826587f",
		Role: "MANAGER",
		AuthUpdatedAt: time.Now().Add(-5 * time.Minute),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
}

func deleteUserAndProductsAfterTest() {
	db := testutils.GetDB()
	if err := db.Exec(`DELETE FROM product_models`).Error; err != nil {
		log.Panic(err)
	}

	if err := db.Exec(`DELETE FROM user_models`).Error; err != nil {
		log.Panic(err)
	}
}

func TestMain(m *testing.M) {
	app = testutils.SetupTestApp()
	createUserTest()
	code := m.Run()
	deleteUserAndProductsAfterTest()
	os.Exit(code)
}