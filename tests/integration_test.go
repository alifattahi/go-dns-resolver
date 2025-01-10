package tests

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestIntegration(t *testing.T) {
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
	body, err := io.ReadAll(resp.Body) // Use io.ReadAll instead of ioutil.ReadAll
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	if !strings.Contains(string(body), "snapp.ir") {
		t.Errorf("Response body does not contain 'snapp.ir': %s", string(body))
	}
}
