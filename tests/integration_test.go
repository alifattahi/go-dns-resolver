package tests

import (
	"context"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func waitForServer(url string, timeout time.Duration) error {
	start := time.Now()
	for {
		if time.Since(start) > timeout {
			return context.DeadlineExceeded
		}
		resp, err := http.Get(url)
		if err == nil && resp.StatusCode == http.StatusOK {
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func startServer(ctx context.Context) (*exec.Cmd, error) {
	cmd := exec.CommandContext(ctx, "go", "run", "../cmd/main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return cmd, nil
}

func TestIntegration(t *testing.T) {
	// Context with timeout for the server process
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// Start the server
	server, err := startServer(ctx)
	if err != nil {
		t.Fatalf("Failed to start the server: %v", err)
	}

	// Ensure the server process is cleaned up
	defer func() {
		if server.Process != nil {
			if err := server.Process.Kill(); err != nil {
				t.Errorf("Failed to stop server: %v", err)
			}
		}
	}()

	// Wait for the server to become ready
	err = waitForServer("http://localhost:3000/ready", 10*time.Second)
	if err != nil {
		t.Fatalf("Server did not become ready in time: %v", err)
	}

	// Send HTTP request
	resp, err := http.Get("http://localhost:3000/resolve?domain=snapp.ir")
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
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
