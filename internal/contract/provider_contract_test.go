package contract

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/leonardo849/product_supermarket/internal/config"
)

func TestPactProvider(t *testing.T) {
	cfg := config.Load()
	cmd := exec.Command("go", "run", "../../cmd/product_service/main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	providerBaseUrl := "http://localhost:" + cfg.HTTPPort
	time.Sleep(3 * time.Second)


	pactExecutable := filepath.Join(cfg.PactPath, "pact_verifier_cli.exe")

	verify := exec.Command(
		pactExecutable,
		"--provider-base-url", providerBaseUrl,
		"--provider-states-url", providerBaseUrl+"/_pact/provider-states",
		"--broker-base-url", cfg.PactAddress,
	)
	



	verify.Stdout = os.Stdout
	verify.Stderr = os.Stderr

	if err := verify.Run(); err != nil {
		t.Fatal(err)
	}
}
