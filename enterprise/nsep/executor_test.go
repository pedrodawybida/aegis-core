package nsep

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/pedrodawybida/nexo-hub/internal/audit"
	"github.com/pedrodawybida/nexo-hub/internal/compliance"
)

func TestNSEPExecutor(t *testing.T) {
	// 1. Setup temporary log file
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "nsep_audit.log")
	logger, err := audit.NewLogger(logPath)
	if err != nil {
		t.Fatalf("Failed to create audit logger: %v", err)
	}
	defer logger.Close()

	// 2. Setup operation mapping
	operations := map[string]OperationTarget{
		"list_clientes": {Method: "GET", Path: "/clientes"},
		"check_saldo":   {Method: "GET", Path: "/saldo/detalhes"},
		"delete_conta":  {Method: "DELETE", Path: "/transacoes/99"},
	}

	executor := NewExecutor(operations, logger, 2*time.Second, 10)

	// 3. Setup agent policy with LGPD & BACEN_538 active
	policy := compliance.AgentPolicy{
		AgentID:     "fintech-agent",
		ActiveModes: []compliance.Mode{compliance.ModeLGPD, compliance.ModeBacen538},
	}

	t.Run("Allowed Operations Flow with Correlated Audit Log", func(t *testing.T) {
		script := `
			var r1 = nsep.call("check_saldo", { id: 10 });
			r1.status;
		`

		res, err := executor.ExecuteScript("fintech-agent", policy, script, func(target OperationTarget, params map[string]interface{}) (interface{}, error) {
			return map[string]interface{}{"status": "ok", "balance": 1500}, nil
		})

		if err != nil {
			t.Fatalf("Unexpected script error: %v", err)
		}

		if res.Status != "SUCCESS" {
			t.Errorf("Expected SUCCESS, got %s", res.Status)
		}

		if res.TotalCalls != 1 {
			t.Errorf("Expected 1 total call, got %d", res.TotalCalls)
		}

		// Verify audit log content
		logBytes, _ := os.ReadFile(logPath)
		logStr := string(logBytes)

		if !strings.Contains(logStr, res.ExecutionID) {
			t.Errorf("Expected log to contain execution_id '%s'", res.ExecutionID)
		}

		if !strings.Contains(logStr, "sequence_in_execution\":1") {
			t.Errorf("Expected sequence_in_execution: 1 in log")
		}

		if !strings.Contains(logStr, "ALLOWED") {
			t.Errorf("Expected ALLOWED status in audit log")
		}
	})

	t.Run("Blocked Operations Flow Halts Execution Immediately", func(t *testing.T) {
		script := `
			nsep.call("check_saldo", {});
			nsep.call("delete_conta", { id: 99 }); // Should block!
			nsep.call("check_saldo", {});
		`

		res, err := executor.ExecuteScript("fintech-agent", policy, script, nil)
		if err == nil {
			t.Fatalf("Expected script execution error due to compliance block, got nil")
		}

		if res.Status != "BLOCKED_OR_FAILED" {
			t.Errorf("Expected BLOCKED_OR_FAILED, got %s", res.Status)
		}

		// Verify execution halted after second call
		if res.TotalCalls != 2 {
			t.Errorf("Expected execution to halt after 2nd call, got %d calls", res.TotalCalls)
		}

		logBytes, _ := os.ReadFile(logPath)
		logStr := string(logBytes)

		if !strings.Contains(logStr, "BLOCKED_BACEN_538_MUTATION_DENIED") {
			t.Errorf("Expected log to record BLOCKED_BACEN_538_MUTATION_DENIED")
		}
	})
}
