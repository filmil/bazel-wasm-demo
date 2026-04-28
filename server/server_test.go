package main_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/bazelbuild/rules_go/go/tools/bazel"
)

func TestServerProxyPath(t *testing.T) {
	serverBin, _ := bazel.Runfile(os.Getenv("SERVER_BIN"))
	if serverBin == "" {
		t.Skip("SERVER_BIN not set, skipping integration test")
	}

	port := 8082
	proxyPath := "/here"

	// Start the server with --proxy-path
	cmd := exec.Command(serverBin,
		"--port", fmt.Sprintf("%d", port),
		"--proxy-path", proxyPath,
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

	baseURL := fmt.Sprintf("http://localhost:%d", port)
	client := &http.Client{Timeout: 2 * time.Second}

	ready := false
	for i := 0; i < 10; i++ {
		resp, err := client.Get(baseURL + proxyPath + "/")
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				ready = true
				break
			}
		}
		time.Sleep(500 * time.Millisecond)
	}

	if !ready {
		t.Fatalf("server not ready at proxy path")
	}

	tests := []struct {
		name         string
		path         string
		expectedCode int
	}{
		{"Proxy Root", proxyPath + "/", http.StatusOK},
		{"Proxy Static", proxyPath + "/web/bootstrap.min.css", http.StatusOK},
		{"Proxy Favicon", proxyPath + "/favicon.ico", http.StatusOK},
		{"Direct Root (should fail or be handled by mux)", "/", http.StatusOK}, // Mux will handle / even without prefix
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := client.Get(baseURL + tc.path)
			if err != nil {
				t.Fatalf("failed to GET %s: %v", tc.path, err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tc.expectedCode {
				t.Errorf("expected status %d, got %d for %s", tc.expectedCode, resp.StatusCode, tc.path)
			}
		})
	}
	
	// Test X-Forwarded-Prefix header
	t.Run("X-Forwarded-Prefix", func(t *testing.T) {
		req, _ := http.NewRequest("GET", baseURL+"/", nil)
		req.Header.Set("X-Forwarded-Prefix", "/there")
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("failed to GET with header: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
		}
		
		body, _ := io.ReadAll(resp.Body)
		if !bytes.Contains(body, []byte("/there/web/")) {
			t.Errorf("expected body to contain /there/web/, but it didn't")
		}
	})
}
