//go:build integration
// +build integration

package user_test

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	// "github.com/leonardo849/product_supermarket/internal/test_utils"
)

func TestIfUserIsInError(t *testing.T) {
	req := httptest.NewRequest("GET", "/user/"+"695a93841f613156b11515f9"+"/permissions/errors", nil)
	req.Header.Set("Authorization", "Bearer fake-token")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

