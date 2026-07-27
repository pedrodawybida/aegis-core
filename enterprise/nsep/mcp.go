// Package nsep provides the MCP (Model Context Protocol) compatibility layer for Nexo Hub.
package nsep

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pedrodawybida/nexo-hub/internal/audit"
	"github.com/pedrodawybida/nexo-hub/internal/compliance"
)

// MCPRequest represents an incoming JSON-RPC 2.0 request from an MCP client (Claude, Cursor, etc.).
type MCPRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// MCPResponse represents a JSON-RPC 2.0 response payload.
type MCPResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *MCPError   `json:"error,omitempty"`
}

// MCPError represents a standard JSON-RPC 2.0 error object.
type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// MCPToolCallParams holds the arguments when invoking tools/call.
type MCPToolCallParams struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}

// MCPServer manages the MCP protocol endpoint and tool dispatches.
type MCPServer struct {
	executor  *Executor
	discovery *DiscoveryEngine
}

// NewMCPServer initializes a new MCP Server instance wrapping the NSEP Executor and Discovery Engine.
func NewMCPServer(executor *Executor, discovery *DiscoveryEngine) *MCPServer {
	return &MCPServer{
		executor:  executor,
		discovery: discovery,
	}
}

// GetToolListReturns static schemas for nsep.search and nsep.execute tools.
func (m *MCPServer) GetToolList() map[string]interface{} {
	return map[string]interface{}{
		"tools": []map[string]interface{}{
			{
				"name":        "nsep.search",
				"description": "Search target API operations by intent/keyword to retrieve operation IDs with minimal prompt context footprint.",
				"inputSchema": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"query": map[string]interface{}{
							"type":        "string",
							"description": "Keyword or domain intent to search for (e.g. 'clientes', 'saldo', 'BACEN_538')",
						},
					},
				},
			},
			{
				"name":        "nsep.execute",
				"description": "Execute JavaScript orchestration code in an isolated Nexo Hub sandbox. Internal calls via nsep.call() pass through BACEN/LGPD compliance pipeline with a correlated execution_id.",
				"inputSchema": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"code": map[string]interface{}{
							"type":        "string",
							"description": "JavaScript orchestration function body calling nsep.call(opID, params)",
						},
					},
					"required": []string{"code"},
				},
			},
		},
	}
}

// HandleHTTP handles incoming HTTP POST JSON-RPC requests from MCP clients.
func (m *MCPServer) HandleHTTP(agentID string, policy compliance.AgentPolicy, w http.ResponseWriter, r *http.Request, targetHandler func(target OperationTarget, params map[string]interface{}) (interface{}, error)) {
	w.Header().Set("Content-Type", "application/json")

	var req MCPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(MCPResponse{
			JSONRPC: "2.0",
			ID:      nil,
			Error:   &MCPError{Code: -32700, Message: "Parse error: invalid JSON"},
		})
		return
	}

	resp := MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
	}

	switch req.Method {
	case "tools/list":
		resp.Result = m.GetToolList()

	case "tools/call":
		var callParams MCPToolCallParams
		if err := json.Unmarshal(req.Params, &callParams); err != nil {
			resp.Error = &MCPError{Code: -32602, Message: "Invalid params"}
			break
		}

		switch callParams.Name {
		case "nsep.search":
			query := m.parseStringArg(callParams.Arguments, "query")
			results := m.discovery.Search(query, 10)
			resp.Result = map[string]interface{}{
				"content": []map[string]interface{}{
					{
						"type": "text",
						"text": m.toJSONString(results),
					},
				},
			}

		case "nsep.execute":
			code := m.parseStringArg(callParams.Arguments, "code")
			if code == "" {
				resp.Error = &MCPError{Code: -32602, Message: "Missing required 'code' argument"}
				break
			}

			execRes, err := m.executor.ExecuteScript(agentID, policy, code, targetHandler)
			if err != nil {
				resp.Result = map[string]interface{}{
					"isError": true,
					"content": []map[string]interface{}{
						{
							"type": "text",
							"text": fmt.Sprintf("NSEP Execution Blocked/Failed: %v", err),
						},
					},
				}
				break
			}

			resp.Result = map[string]interface{}{
				"content": []map[string]interface{}{
					{
						"type": "text",
						"text": m.toJSONString(execRes),
					},
				},
			}

		default:
			resp.Error = &MCPError{Code: -32601, Message: fmt.Sprintf("Tool '%s' not found", callParams.Name)}
		}

	default:
		resp.Error = &MCPError{Code: -32601, Message: fmt.Sprintf("Method '%s' not supported", req.Method)}
	}

	json.NewEncoder(w).Encode(resp)
}

func (m *MCPServer) parseStringArg(raw json.RawMessage, key string) string {
	if len(raw) == 0 {
		return ""
	}
	var rawMap map[string]interface{}
	if err := json.Unmarshal(raw, &rawMap); err == nil {
		if val, ok := rawMap[key].(string); ok {
			return val
		}
	}
	var strVal string
	if err := json.Unmarshal(raw, &strVal); err == nil {
		var innerMap map[string]interface{}
		if err := json.Unmarshal([]byte(strVal), &innerMap); err == nil {
			if val, ok := innerMap[key].(string); ok {
				return val
			}
		}
		return strVal
	}
	return ""
}

func (m *MCPServer) toJSONString(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}

// BuildDefaultMCPServer helper creates an MCP Server pre-configured with discovery engine and executor.
func BuildDefaultMCPServer(logger *audit.Logger, operations map[string]OperationTarget) *MCPServer {
	disc := NewDiscoveryEngine(nil)
	exec := NewExecutor(operations, logger, 5*time.Second, 50)
	return NewMCPServer(exec, disc)
}
