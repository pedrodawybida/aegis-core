package proxy

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/pedrodawybida/aegis-core/internal/audit"
	"github.com/pedrodawybida/aegis-core/internal/compliance"
)

func TestAegisProxyIntegration(t *testing.T) {
	// 1. Setup mock target upstream backend server
	mockBackend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Target API Response OK"))
	}))
	defer mockBackend.Close()

	// 2. Setup audit logger
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "proxy_audit.log")
	logger, err := audit.NewLogger(logPath)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	// 3. Setup policy DB
	policyDB := map[string]compliance.AgentPolicy{
		"token-fintech": {
			AgentID:     "token-fintech",
			ActiveModes: []compliance.Mode{compliance.ModeLGPD, compliance.ModeBacen538},
		},
	}

	agProxy, err := NewAegisProxy(mockBackend.URL, logger, policyDB)
	if err != nil {
		t.Fatalf("Failed to create AegisProxy: %v", err)
	}

	proxyServer := httptest.NewServer(agProxy)
	defer proxyServer.Close()

	client := &http.Client{}

	t.Run("Health Check Endpoint", func(t *testing.T) {
		resp, err := client.Get(fmt.Sprintf("%s/_aegis/health", proxyServer.URL))
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("Missing Authorization Header", func(t *testing.T) {
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/data", proxyServer.URL), nil)
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401 Unauthorized, got %d", resp.StatusCode)
		}
	})

	t.Run("Blocked Request - LGPD Bulk Data", func(t *testing.T) {
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/clientes", proxyServer.URL), nil)
		req.Header.Set("Authorization", "Bearer token-fintech")
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusForbidden {
			t.Errorf("Expected status 403 Forbidden, got %d", resp.StatusCode)
		}
		statusHeader := resp.Header.Get("X-Aegis-Compliance-Status")
		if statusHeader != "BLOCKED_LGPD_BULK_DATA_ACCESS_DENIED" {
			t.Errorf("Expected header BLOCKED_LGPD_BULK_DATA_ACCESS_DENIED, got '%s'", statusHeader)
		}
	})

	t.Run("Allowed Request - Forward to Backend", func(t *testing.T) {
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/saldo/detalhes", proxyServer.URL), nil)
		req.Header.Set("Authorization", "Bearer token-fintech")
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200 OK, got %d", resp.StatusCode)
		}
		body, _ := io.ReadAll(resp.Body)
		if string(body) != "Target API Response OK" {
			t.Errorf("Expected body 'Target API Response OK', got '%s'", string(body))
		}
	})

	t.Run("Dry-Run Mode - Log Violation but Allow Forwarding", func(t *testing.T) {
		agProxy.SetDryRun(true)
		defer agProxy.SetDryRun(false)

		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/clientes", proxyServer.URL), nil)
		req.Header.Set("Authorization", "Bearer token-fintech")
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		defer resp.Body.Close()

		// Under Dry-Run mode, status code should be 200 OK because request was forwarded!
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200 OK under Dry-Run, got %d", resp.StatusCode)
		}

		dryRunHeader := resp.Header.Get("X-Aegis-Dry-Run")
		if dryRunHeader != "true" {
			t.Errorf("Expected X-Aegis-Dry-Run header to be 'true', got '%s'", dryRunHeader)
		}

		statusHeader := resp.Header.Get("X-Aegis-Compliance-Status")
		if statusHeader != "DRY_RUN_BLOCKED_LGPD_BULK_DATA_ACCESS_DENIED" {
			t.Errorf("Expected DRY_RUN_BLOCKED_LGPD_BULK_DATA_ACCESS_DENIED status, got '%s'", statusHeader)
		}
	})
}
