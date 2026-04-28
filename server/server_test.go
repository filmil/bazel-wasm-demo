package main_test

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/bazelbuild/rules_go/go/tools/bazel"
)

func TestServerResources(t *testing.T) {
	// Locate the server binary. Bazel sets TEST_WORKSPACE and TEST_SRCDIR.
	serverBin, _ := bazel.Runfile(os.Getenv("SERVER_BIN"))
	if serverBin == "" {
		t.Skip("SERVER_BIN not set, skipping integration test")
	}

	// Choose an ephemeral port
	port := 8081

	// Start the server
	cmd := exec.Command(serverBin, 
		"--port", fmt.Sprintf("%d", port),
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start server: %v", err)
	}
	defer func() {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
	}()

	// Wait for server to be ready
	baseURL := fmt.Sprintf("http://localhost:%d", port)
	client := &http.Client{Timeout: 2 * time.Second}
	
	ready := false
	for i := 0; i < 10; i++ {
		resp, err := client.Get(baseURL + "/")
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode != http.StatusNotFound {
				ready = true
				break
			}
		}
		time.Sleep(500 * time.Millisecond)
	}

	if !ready {
		t.Logf("server might not be fully ready or root returned 404")
	}

	// Define test cases for resources
	tests := []struct {
		name         string
		path         string
		expectedCode int
		minSize      int64
	}{
		{"Root HTML", "/", http.StatusOK, 100},
		{"WASM Binary", "/web/app.wasm", http.StatusOK, 1024 * 1024}, // Expecting >1MB for Go WASM
		{"Icon", "/web/icon.png", http.StatusOK, 100},
		{"Favicon", "/favicon.ico", http.StatusOK, 50},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := client.Get(baseURL + tc.path)
			if err != nil {
				t.Fatalf("failed to GET %s: %v", tc.path, err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tc.expectedCode {
				t.Errorf("expected status %d, got %d", tc.expectedCode, resp.StatusCode)
			}

			body, _ := io.ReadAll(resp.Body)
			if int64(len(body)) < tc.minSize {
				t.Errorf("expected body size > %d, got %d", tc.minSize, len(body))
			}
		})
	}
}
