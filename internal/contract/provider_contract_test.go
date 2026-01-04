package contract

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/leonardo849/product_supermarket/internal/config"
	"github.com/leonardo849/product_supermarket/internal/infrastructure/bootstrap"
	"github.com/pact-foundation/pact-go/v2/models"
	"github.com/pact-foundation/pact-go/v2/provider"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func fiberToHTTPServer(app *fiber.App) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := r.Clone(r.Context())
		req.RequestURI = ""

		resp, err := app.Test(req, -1)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		for k, vv := range resp.Header {
			for _, v := range vv {
				w.Header().Add(k, v)
			}
		}

		w.WriteHeader(resp.StatusCode)
		_, _ = io.Copy(w, resp.Body)
	}))
}


func waitForProviderReady(t *testing.T, baseURL string) {
	deadline := time.Now().Add(10 * time.Second)

	for time.Now().Before(deadline) {
		resp, err := http.Get(baseURL + "/health")
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			return
		}
		time.Sleep(100 * time.Millisecond)
	}

	t.Fatal("provider not ready")
}


func closeDb(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	sqlDB.Close()
}



func TestPactProviderV4(t *testing.T) {
	cfg := config.Load()

	container, err := bootstrap.BuildApp(&cfg, false)
	require.NoError(t, err)

	server := fiberToHTTPServer(container.App)
	defer server.Close()

	defer closeDb(container.DB)
	defer container.Redis.Close()

	waitForProviderReady(t, server.URL)
	var providerVersion string = time.Now().Format("20060102-150405")
	if cfg.GitCommit != "none" {
		providerVersion = cfg.GitCommit
	}
	
	verifier := provider.NewVerifier()
	err = verifier.VerifyProvider(t, provider.VerifyRequest{
	Provider:        "product_service",   
	ProviderBaseURL: server.URL,

	BrokerURL: cfg.PactAddress,
	BrokerUsername: cfg.PactUsername,
	BrokerPassword: cfg.PactPassword,
	PublishVerificationResults: true,
	ProviderVersion: providerVersion,    
	ProviderBranch:            cfg.GitBranch,
	StateHandlers: models.StateHandlers{
		"user exists and permissions were evaluated": func(
			setup bool,
			state models.ProviderState,
		) (models.ProviderStateResponse, error) {
			return nil, nil
		},
	},
})


	require.NoError(t, err)
}

