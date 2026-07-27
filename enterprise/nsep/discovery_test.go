package nsep

import (
	"testing"
	"time"

	"github.com/pedrodawybida/nexo-hub/internal/compliance"
	"github.com/pedrodawybida/nexo-hub/internal/config"
)

func TestNSEPDiscovery(t *testing.T) {
	catalog := []OperationMetadata{
		{
			ID:             "list_clientes",
			Method:         "GET",
			Path:           "/clientes",
			Description:    "List customer data records",
			ComplianceMode: "LGPD",
			Params:         map[string]string{"limit": "number"},
		},
		{
			ID:             "delete_conta",
			Method:         "DELETE",
			Path:           "/contas/{id}",
			Description:    "Delete account mutation",
			ComplianceMode: "BACEN_538",
			Params:         map[string]string{"id": "string"},
		},
		{
			ID:             "check_saldo",
			Method:         "GET",
			Path:           "/saldo/detalhes",
			Description:    "Retrieve account balance",
			ComplianceMode: "BACEN_538",
			Params:         map[string]string{"account_id": "string"},
		},
	}

	discovery := NewDiscoveryEngine(catalog)

	t.Run("Search Filters Operations by Keyword", func(t *testing.T) {
		matches := discovery.Search("clientes", 5)
		if len(matches) != 1 {
			t.Fatalf("Expected 1 match for 'clientes', got %d", len(matches))
		}
		if matches[0].ID != "list_clientes" {
			t.Errorf("Expected operation list_clientes, got %s", matches[0].ID)
		}
	})

	t.Run("Search Returns Multiple Operations Matching Compliance Mode", func(t *testing.T) {
		matches := discovery.Search("BACEN_538", 5)
		if len(matches) != 2 {
			t.Fatalf("Expected 2 matches for BACEN_538, got %d", len(matches))
		}
	})

	t.Run("LoadFromConfig Auto-populates Catalog", func(t *testing.T) {
		cfg := &config.Config{
			Agents: []config.AgentConfig{
				{
					ID:    "agent-test",
					Modes: []compliance.Mode{compliance.ModeLGPD, compliance.ModeBacen538, compliance.ModeCFM},
				},
			},
		}

		engine := NewDiscoveryEngine(nil)
		engine.LoadFromConfig(cfg)

		matches := engine.Search("", 10)
		if len(matches) < 3 {
			t.Errorf("Expected at least 3 auto-populated operations, got %d", len(matches))
		}
	})

	t.Run("nsep.search Binding Inside JS Sandbox", func(t *testing.T) {
		sandbox := NewSandbox(2*time.Second, 10, func(opID string, params map[string]interface{}) (interface{}, error) {
			return nil, nil
		})

		sandbox.SearchHandler = func(query string) interface{} {
			return discovery.Search(query, 5)
		}

		script := `
			var ops = nsep.search("clientes");
			ops.length;
		`

		res, err := sandbox.Execute(script)
		if err != nil {
			t.Fatalf("JS Execution failed: %v", err)
		}

		if lenVal, ok := res.(int64); !ok || lenVal != 1 {
			t.Errorf("Expected JS nsep.search to return length 1, got %v", res)
		}
	})
}
