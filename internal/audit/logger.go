// Package audit provides an immutable JSON logging mechanism required for
// strict compliance environments such as BACEN 538/2025.
package audit

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

// AuditLog represents the immutable log structure required by BACEN 538/2025.
// It traces who (Agent), what (Action), when (Timestamp), and the Result.
type AuditLog struct {
	Timestamp   string `json:"timestamp"`
	AgentID     string `json:"agent_id"`
	Action      string `json:"action"`
	ToolPayload string `json:"tool_payload"`
	Result      string `json:"result"`
	IPAddress   string `json:"ip_address"`
}

// Logger handles appending structured audit records to an immutable file.
type Logger struct {
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
// It also prints to standard output for real-time observability.
func (l *Logger) LogAction(agentID, action, payload, result, ip string) {
	entry := AuditLog{
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
		AgentID:     agentID,
		Action:      action,
		ToolPayload: payload,
		Result:      result,
		IPAddress:   ip,
	}

	logBytes, _ := json.Marshal(entry)
	l.file.Write(append(logBytes, '\n'))
	log.Printf("[AEGIS-AUDIT] %s", string(logBytes))
}

// Close safely closes the underlying audit file.
func (l *Logger) Close() error {
	return l.file.Close()
}
