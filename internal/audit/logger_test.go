package audit

import (
	"bufio"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func TestConcurrentLogger(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "audit_test.log")

	logger, err := NewLogger(logPath)
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}

	var wg sync.WaitGroup
	goroutines := 50
	logsPerGoroutine := 20

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(agentID int) {
			defer wg.Done()
			for j := 0; j < logsPerGoroutine; j++ {
				logger.LogAction("agent-test", "GET /test", `{"param":"val"}`, "ALLOWED", "127.0.0.1")
			}
		}(i)
	}

	wg.Wait()
	if err := logger.Close(); err != nil {
		t.Fatalf("Failed to close logger: %v", err)
	}

	// Verify line count
	file, err := os.Open(logPath)
	if err != nil {
		t.Fatalf("Failed to open log file for verification: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}

	expectedCount := goroutines * logsPerGoroutine
	if lineCount != expectedCount {
		t.Errorf("Expected %d log lines, got %d", expectedCount, lineCount)
	}
}
