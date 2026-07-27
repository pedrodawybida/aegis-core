package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "aegis.yaml")

	yamlData := `
target_api: "http://api.internal:8000"
agents:
  - id: "agent-1"
    modes:
      - "LGPD"
  - id: "agent-2"
    modes:
      - "BACEN_538"
`
	if err := os.WriteFile(configPath, []byte(yamlData), 0644); err != nil {
		t.Fatalf("Failed to create temporary config: %v", err)
	}

	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.TargetAPI != "http://api.internal:8000" {
		t.Errorf("Expected TargetAPI 'http://api.internal:8000', got '%s'", cfg.TargetAPI)
	}

	policyMap := cfg.ToPolicyMap()
	if len(policyMap) != 2 {
		t.Errorf("Expected 2 policies in policyMap, got %d", len(policyMap))
	}

	if p1, ok := policyMap["agent-1"]; !ok || len(p1.ActiveModes) != 1 || p1.ActiveModes[0] != "LGPD" {
		t.Errorf("Agent-1 policy mismatch: %+v", p1)
	}
}
