//go:build integration
// +build integration

package user_test

import (
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/leonardo849/product_supermarket/internal/test_utils"
)

var app *fiber.App

func TestMain(m *testing.M) {
	app = testutils.SetupTestApp()

	code := m.Run()
	os.Exit(code)
}
