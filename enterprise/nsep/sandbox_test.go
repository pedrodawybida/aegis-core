package nsep

import (
	"errors"
	"testing"
	"time"
)

func TestNSEPSandbox(t *testing.T) {
	t.Run("Basic JS Execution", func(t *testing.T) {
		sandbox := NewSandbox(2*time.Second, 10, func(opID string, params map[string]interface{}) (interface{}, error) {
			return nil, nil
		})

		res, err := sandbox.Execute("var a = 10; var b = 20; a + b;")
		if err != nil {
			t.Fatalf("Unexpected execution error: %v", err)
		}

		if val, ok := res.(int64); !ok || val != 30 {
			t.Errorf("Expected 30, got %v (%T)", res, res)
		}
	})

	t.Run("nsep.call Binding Interception", func(t *testing.T) {
		calledOp := ""
		calledParams := make(map[string]interface{})

		mockHandler := func(opID string, params map[string]interface{}) (interface{}, error) {
			calledOp = opID
			calledParams = params
			return map[string]interface{}{
				"status": "success",
				"id":     101,
			}, nil
		}

		sandbox := NewSandbox(2*time.Second, 10, mockHandler)

		script := `
			var res = nsep.call("list_clientes", { limit: 50 });
			res.status;
		`

		res, err := sandbox.Execute(script)
		if err != nil {
			t.Fatalf("Execution failed: %v", err)
		}

		if calledOp != "list_clientes" {
			t.Errorf("Expected operation 'list_clientes', got '%s'", calledOp)
		}

		if limit, ok := calledParams["limit"].(int64); !ok || limit != 50 {
			t.Errorf("Expected limit 50 in params, got %v", calledParams["limit"])
		}

		if str, ok := res.(string); !ok || str != "success" {
			t.Errorf("Expected 'success', got %v", res)
		}
	})

	t.Run("Execution Timeout Protection", func(t *testing.T) {
		sandbox := NewSandbox(100*time.Millisecond, 10, func(opID string, params map[string]interface{}) (interface{}, error) {
			return nil, nil
		})

		// Infinite loop script
		script := `while(true) {}`

		_, err := sandbox.Execute(script)
		if err == nil {
			t.Fatalf("Expected timeout error, got nil")
		}

		if !errors.Is(err, ErrExecutionTimeout) {
			t.Errorf("Expected ErrExecutionTimeout, got %v", err)
		}
	})

	t.Run("Max Calls Limit Enforcement", func(t *testing.T) {
		sandbox := NewSandbox(2*time.Second, 3, func(opID string, params map[string]interface{}) (interface{}, error) {
			return "ok", nil
		})

		// Script making 5 calls when limit is 3
		script := `
			for (var i = 0; i < 5; i++) {
				nsep.call("ping", {});
			}
		`

		_, err := sandbox.Execute(script)
		if err == nil {
			t.Fatalf("Expected call limit error, got nil")
		}
	})

	t.Run("Environment Isolation", func(t *testing.T) {
		sandbox := NewSandbox(2*time.Second, 10, func(opID string, params map[string]interface{}) (interface{}, error) {
			return nil, nil
		})

		script := `typeof process === 'undefined' && typeof require === 'undefined'`
		res, err := sandbox.Execute(script)
		if err != nil {
			t.Fatalf("Execution failed: %v", err)
		}

		if isIsolated, ok := res.(bool); !ok || !isIsolated {
			t.Errorf("Expected environment to be fully isolated without process/require globals")
		}
	})
}
