package tests

import (
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
	"log"
)

func startServer() (*exec.Cmd, error) {
	// Start the server in the background
	cmd := exec.Command("go", "run", "../cmd/main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		return nil, err
	}
	return cmd, nil
}

func TestIntegration(t *testing.T) {
	// Start the server
	server, err := startServer()
	if err != nil {
		t.Fatalf("Failed to start the server: %v", err)
	}

	// Make sure to stop the server after the test is done
	defer func() {
		if server.Process != nil {
			if err := server.Process.Kill(); err != nil {
				log.Fatalf("Failed to stop server: %v", err)
			}
		}
	}()

	// Wait for a while to make sure the server is up and running
	time.Sleep(5 * time.Second)

	// Send HTTP request
	resp, err := http.Get("http://localhost:3000/resolve?domain=snapp.ir")
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Check response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	if !strings.Contains(string(body), "snapp.ir") {
		t.Errorf("Response body does not contain 'snapp.ir': %s", string(body))
	}
}
