// Package nsep provides the Discovery Layer for filtering target API operations down to minimal token footprints.
package nsep

import (
	"strings"

	"github.com/pedrodawybida/nexo-hub/internal/config"
)

// OperationMetadata represents a lightweight operation descriptor returned to AI Agents during discovery.
type OperationMetadata struct {
	ID             string            `json:"id"`
	Method         string            `json:"method"`
	Path           string            `json:"path"`
	Description    string            `json:"description,omitempty"`
	ComplianceMode string            `json:"compliance_mode,omitempty"`
	Params         map[string]string `json:"params,omitempty"`
}

// DiscoveryEngine indexes and searches target API operations to minimize token context usage.
type DiscoveryEngine struct {
	catalog []OperationMetadata
}

// NewDiscoveryEngine creates a new DiscoveryEngine with the given operation catalog.
func NewDiscoveryEngine(catalog []OperationMetadata) *DiscoveryEngine {
	if catalog == nil {
		catalog = make([]OperationMetadata, 0)
	}
	return &DiscoveryEngine{catalog: catalog}
}

// LoadFromConfig populates the discovery catalog from nexo.yaml agent policy definitions.
func (d *DiscoveryEngine) LoadFromConfig(cfg *config.Config) {
	if cfg == nil {
		return
	}

	for _, agent := range cfg.Agents {
		for _, mode := range agent.Modes {
			modeStr := string(mode)

			switch modeStr {
			case "LGPD":
				d.AddOperation(OperationMetadata{
					ID:             "list_clientes",
					Method:         "GET",
					Path:           "/clientes",
					Description:    "List customer data records",
					ComplianceMode: "LGPD",
					Params:         map[string]string{"limit": "number"},
				})
			case "BACEN_538":
				d.AddOperation(OperationMetadata{
					ID:             "delete_conta",
					Method:         "DELETE",
					Path:           "/contas/{id}",
					Description:    "Delete account mutation",
					ComplianceMode: "BACEN_538",
					Params:         map[string]string{"id": "string"},
				})
			case "CFM":
				d.AddOperation(OperationMetadata{
					ID:             "get_prontuario",
					Method:         "GET",
					Path:           "/prontuarios/{id}",
					Description:    "Retrieve patient medical record",
					ComplianceMode: "CFM",
					Params:         map[string]string{"id": "string"},
				})
			}
		}
	}
}

// AddOperation appends an operation metadata record to the catalog if not already present.
func (d *DiscoveryEngine) AddOperation(op OperationMetadata) {
	for _, existing := range d.catalog {
		if existing.ID == op.ID {
			return
		}
	}
	d.catalog = append(d.catalog, op)
}

// Search filters the operation catalog against the user or agent query string.
// Returns a compact list of matching operations to minimize prompt token footprint.
func (d *DiscoveryEngine) Search(query string, limit int) []OperationMetadata {
	if limit <= 0 {
		limit = 10
	}

	query = strings.ToLower(strings.TrimSpace(query))
	if query == "" {
		if len(d.catalog) > limit {
			return d.catalog[:limit]
		}
		return d.catalog
	}

	var results []OperationMetadata
	for _, op := range d.catalog {
		matchScore := 0

		if strings.Contains(strings.ToLower(op.ID), query) {
			matchScore += 3
		}
		if strings.Contains(strings.ToLower(op.Path), query) {
			matchScore += 2
		}
		if strings.Contains(strings.ToLower(op.Description), query) {
			matchScore += 1
		}
		if strings.Contains(strings.ToLower(op.ComplianceMode), query) {
			matchScore += 1
		}

		if matchScore > 0 {
			results = append(results, op)
		}

		if len(results) >= limit {
			break
		}
	}

	return results
}
