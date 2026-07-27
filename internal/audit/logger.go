// Package audit provides an immutable JSON logging mechanism required for
// strict compliance environments such as BACEN 538/2025.
package audit

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"
)

// AuditLog represents the immutable log structure required by BACEN 538/2025.
// It traces who (Agent), what (Action), when (Timestamp), and the Result, including NSEP execution correlation.
type AuditLog struct {
	Timestamp   string `json:"timestamp"`
	AgentID     string `json:"agent_id"`
	ExecutionID string `json:"execution_id,omitempty"`
	Sequence    int    `json:"sequence_in_execution,omitempty"`
	Action      string `json:"action"`
	ToolPayload string `json:"tool_payload"`
	Result      string `json:"result"`
	IPAddress   string `json:"ip_address"`
}

// Logger handles appending structured audit records to an immutable file in a thread-safe manner.
type Logger struct {
	mu   sync.Mutex
	file *os.File
}

// NewLogger creates a new Logger instance. It opens the specified file in
// append-only mode to prevent historical mutation, ensuring audit integrity.
func NewLogger(filePath string) (*Logger, error) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &Logger{file: file}, nil
}

// LogAction constructs an AuditLog entry and appends it to the file as JSON.
// It is thread-safe and safe for concurrent invocations.
func (l *Logger) LogAction(agentID, action, payload, result, ip string) {
	l.LogExecutionAction(agentID, "", action, payload, result, ip, 0)
}

// LogExecutionAction constructs an AuditLog entry with NSEP execution_id correlation.
func (l *Logger) LogExecutionAction(agentID, executionID, action, payload, result, ip string, sequence int) {
	entry := AuditLog{
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
		AgentID:     agentID,
		ExecutionID: executionID,
		Sequence:    sequence,
		Action:      action,
		ToolPayload: payload,
		Result:      result,
		IPAddress:   ip,
	}

	logBytes, err := json.Marshal(entry)
	if err != nil {
		log.Printf("[NEXO-AUDIT-ERROR] Failed to marshal log entry: %v", err)
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if _, err := l.file.Write(append(logBytes, '\n')); err != nil {
		log.Printf("[NEXO-AUDIT-ERROR] Failed to write to audit log file: %v", err)
	}
	log.Printf("[NEXO-AUDIT] %s", string(logBytes))
}

// Close safely flushes and closes the underlying audit file.
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.file.Close()
}
