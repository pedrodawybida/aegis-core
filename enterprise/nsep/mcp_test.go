package nsep

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/pedrodawybida/nexo-hub/internal/audit"
	"github.com/pedrodawybida/nexo-hub/internal/compliance"
)

func TestMCPServer(t *testing.T) {
	// 1. Setup temporary log file & logger
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "mcp_audit.log")
	logger, err := audit.NewLogger(logPath)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	// 2. Setup operations catalog & discovery
	operations := map[string]OperationTarget{
		"list_clientes": {Method: "GET", Path: "/clientes"},
		"check_saldo":   {Method: "GET", Path: "/saldo/detalhes"},
	}

	catalog := []OperationMetadata{
		{ID: "list_clientes", Method: "GET", Path: "/clientes", Description: "List customers", ComplianceMode: "LGPD"},
		{ID: "check_saldo", Method: "GET", Path: "/saldo/detalhes", Description: "Check balance", ComplianceMode: "BACEN_538"},
	}

	discovery := NewDiscoveryEngine(catalog)
	executor := NewExecutor(operations, logger, 2*time.Second, 10)
	mcpServer := NewMCPServer(executor, discovery)

	policy := compliance.AgentPolicy{
		AgentID:     "claude-agent",
		ActiveModes: []compliance.Mode{compliance.ModeLGPD, compliance.ModeBacen538},
	}

	t.Run("MCP tools/list Request Returns nsep.search & nsep.execute", func(t *testing.T) {
		reqBody := `{"jsonrpc":"2.0","id":1,"method":"tools/list"}`
		req := httptest.NewRequest("POST", "/_nexo/mcp", bytes.NewBufferString(reqBody))
		w := httptest.NewRecorder()

		mcpServer.HandleHTTP("claude-agent", policy, w, req, nil)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200 OK, got %d", w.Code)
		}

		var resp MCPResponse
		json.NewDecoder(w.Body).Decode(&resp)

		if resp.Error != nil {
			t.Fatalf("Unexpected MCP error: %v", resp.Error)
		}

		resultMap, ok := resp.Result.(map[string]interface{})
		if !ok {
			t.Fatalf("Expected result map, got %T", resp.Result)
		}

		tools, ok := resultMap["tools"].([]interface{})
		if !ok || len(tools) != 2 {
			t.Fatalf("Expected 2 tools in tools/list, got %v", tools)
		}
	})

	t.Run("MCP tools/call for nsep.search", func(t *testing.T) {
		reqBody := `{
			"jsonrpc": "2.0",
			"id": 2,
			"method": "tools/call",
			"params": {
				"name": "nsep.search",
				"arguments": {"query": "clientes"}
			}
		}`

		req := httptest.NewRequest("POST", "/_nexo/mcp", bytes.NewBufferString(reqBody))
		w := httptest.NewRecorder()

		mcpServer.HandleHTTP("claude-agent", policy, w, req, nil)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200 OK, got %d", w.Code)
		}

		var resp MCPResponse
		json.NewDecoder(w.Body).Decode(&resp)

		if resp.Error != nil {
			t.Fatalf("Unexpected error: %v", resp.Error)
		}

		resBytes, _ := json.Marshal(resp.Result)
		if !strings.Contains(string(resBytes), "list_clientes") {
			t.Errorf("Expected search response to contain 'list_clientes', got '%s'", string(resBytes))
		}
	})

	t.Run("MCP tools/call for nsep.execute", func(t *testing.T) {
		reqBody := `{
			"jsonrpc": "2.0",
			"id": 3,
			"method": "tools/call",
			"params": {
				"name": "nsep.execute",
				"arguments": {
					"code": "var r = nsep.call('check_saldo', {account_id: 123}); r.status;"
				}
			}
		}`

		req := httptest.NewRequest("POST", "/_nexo/mcp", bytes.NewBufferString(reqBody))
		w := httptest.NewRecorder()

		mcpServer.HandleHTTP("claude-agent", policy, w, req, func(target OperationTarget, params map[string]interface{}) (interface{}, error) {
			return map[string]interface{}{"status": "ok", "balance": 5000}, nil
		})

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200 OK, got %d", w.Code)
		}

		bodyStr := w.Body.String()
		if !strings.Contains(bodyStr, "exec_") {
			t.Errorf("Expected response to contain correlated execution_id 'exec_'")
		}
	})
}
