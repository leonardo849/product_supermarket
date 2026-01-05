//go:build integration
// +build integration

package product_test

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"log"
	"github.com/stretchr/testify/assert"
)

func TestCreateProduct(t *testing.T) {
	body := `
		{
		"name": "generical product",
		"price": 2332,
		"category": "FOOD",
		"description": "that product is a generical product",
		"stock": {
			"initial_stock": 50,
			"minimum_stock": 44
		}
		}
		`

	req := httptest.NewRequest("POST", "/product/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer fake-token")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t,200, resp.StatusCode)
	bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	log.Println("response body:", string(bodyBytes))
}