// Package nsep provides the execution pipeline connecting the JS Sandbox to Nexo Hub Compliance.
package nsep

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/pedrodawybida/nexo-hub/internal/audit"
	"github.com/pedrodawybida/nexo-hub/internal/compliance"
)

// OperationTarget defines the target HTTP method and path for an operation ID.
type OperationTarget struct {
	Method string `json:"method"`
	Path   string `json:"path"`
}

// ExecutionResult holds the output and compliance metadata of a completed NSEP execution.
type ExecutionResult struct {
	ExecutionID   string      `json:"execution_id"`
	AgentID       string      `json:"agent_id"`
	TotalCalls    int         `json:"total_calls"`
	Output        interface{} `json:"output"`
	Status        string      `json:"status"`
	BlockedReason string      `json:"blocked_reason,omitempty"`
}

// Executor coordinates sandboxed script execution with compliance checks and audit logs.
type Executor struct {
	operations map[string]OperationTarget
	logger     *audit.Logger
	timeout    time.Duration
	maxCalls   int
}

// NewExecutor creates a new NSEP Executor instance.
func NewExecutor(operations map[string]OperationTarget, logger *audit.Logger, timeout time.Duration, maxCalls int) *Executor {
	if operations == nil {
		operations = make(map[string]OperationTarget)
	}
	return &Executor{
		operations: operations,
		logger:     logger,
		timeout:    timeout,
		maxCalls:   maxCalls,
	}
}

// GenerateExecutionID creates a unique correlation ID for an execution session.
func GenerateExecutionID() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("exec_%d_%04x", time.Now().Unix(), r.Intn(0xffff))
}

// ExecuteScript executes a JS script under the identity of an agent policy.
func (e *Executor) ExecuteScript(agentID string, policy compliance.AgentPolicy, script string, targetHandler func(target OperationTarget, params map[string]interface{}) (interface{}, error)) (*ExecutionResult, error) {
	executionID := GenerateExecutionID()
	var mu sync.Mutex
	sequence := 0

	// Interceptor handler for nsep.call()
	interceptor := func(opID string, params map[string]interface{}) (interface{}, error) {
		mu.Lock()
		sequence++
		seq := sequence
		mu.Unlock()

		target, exists := e.operations[opID]
		if !exists {
			// Fallback: infer operation if not found in map
			target = OperationTarget{Method: "POST", Path: "/" + opID}
		}

		payloadBytes, _ := json.Marshal(params)
		payload := string(payloadBytes)
		action := fmt.Sprintf("%s %s", target.Method, target.Path)

		// 1. Compliance Policy Evaluation
		allowed, reason := compliance.EvaluateRequest(policy, target.Method, target.Path)

		// 2. Log Action with execution_id and sequence correlation
		if e.logger != nil {
			e.logger.LogExecutionAction(agentID, executionID, action, payload, reason, "nsep-sandbox", seq)
		}

		// 3. Block execution instantly if compliance fails
		if !allowed {
			return nil, fmt.Errorf("Aegis Compliance Violation: %s", reason)
		}

		// 4. Forward to target handler if allowed
		if targetHandler != nil {
			return targetHandler(target, params)
		}

		return map[string]interface{}{"status": "ok", "operation": opID, "http_code": http.StatusOK}, nil
	}

	sandbox := NewSandbox(e.timeout, e.maxCalls, interceptor)
	output, err := sandbox.Execute(script)

	res := &ExecutionResult{
		ExecutionID: executionID,
		AgentID:     agentID,
		TotalCalls:  sequence,
		Output:      output,
		Status:      "SUCCESS",
	}

	if err != nil {
		res.Status = "BLOCKED_OR_FAILED"
		res.BlockedReason = err.Error()
		return res, err
	}

	return res, nil
}
