// Package config handles the parsing and validation of the Aegis YAML configuration.
package config

import (
	"os"

	"github.com/pedrodawybida/nexo-hub/internal/compliance"
	"gopkg.in/yaml.v3"
)

// AgentConfig represents the configuration of a single agent in the YAML file.
type AgentConfig struct {
	ID    string            `yaml:"id"`
	Modes []compliance.Mode `yaml:"modes"`
}

// Config represents the root structure of aegis.yaml.
type Config struct {
	TargetAPI string        `yaml:"target_api"`
	DryRun    bool          `yaml:"dry_run"`
	Agents    []AgentConfig `yaml:"agents"`
}

// LoadConfig reads and unmarshals the YAML configuration from the given file path.
// It allows environment variable overrides and validates required fields.
func LoadConfig(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	// Environment variable override for TargetAPI
	if envTarget := os.Getenv("AEGIS_TARGET_API"); envTarget != "" {
		cfg.TargetAPI = envTarget
	}

	// Environment variable override for DryRun
	if envDryRun := os.Getenv("AEGIS_DRY_RUN"); envDryRun == "true" || envDryRun == "1" {
		cfg.DryRun = true
	}

	if cfg.TargetAPI == "" {
		cfg.TargetAPI = "http://localhost:9000"
	}

	return &cfg, nil
}

// ToPolicyMap converts the slice of agent configurations into an O(1) lookup map
// for high-performance policy enforcement in the proxy layer.
func (c *Config) ToPolicyMap() map[string]compliance.AgentPolicy {
	policyMap := make(map[string]compliance.AgentPolicy)
	for _, agent := range c.Agents {
		policyMap[agent.ID] = compliance.AgentPolicy{
			AgentID:     agent.ID,
			ActiveModes: agent.Modes,
		}
	}
	return policyMap
}
