// Package compliance defines the regulatory templates and enforcement logic.
package compliance

import "strings"

// Mode represents a specific regulatory compliance standard.
type Mode string

const (
	// ModeBacen538 represents the Brazilian Central Bank Cybersecurity Resolution 538/2025.
	ModeBacen538 Mode = "BACEN_538"

	// ModeLGPD represents the Brazilian General Data Protection Law.
	ModeLGPD Mode = "LGPD"

	// ModeCFM represents the Federal Council of Medicine guidelines.
	ModeCFM Mode = "CFM"
)

// AgentPolicy defines the regulatory modes enforced on a specific Non-Human Identity (Agent).
type AgentPolicy struct {
	AgentID     string
	ActiveModes []Mode
}

// EvaluateRequest checks the incoming HTTP request against the active compliance modes
// of the given Agent. It returns false and a violation reason if the request is non-compliant.
func EvaluateRequest(policy AgentPolicy, method string, path string) (bool, string) {
	// Sanitize path (strip query params and normalize casing)
	cleanPath := strings.ToLower(strings.Split(path, "?")[0])
	cleanMethod := strings.ToUpper(method)

	for _, mode := range policy.ActiveModes {
		switch mode {

		case ModeBacen538:
			// BACEN 538 mandates strict control over unapproved state mutations (DELETE / PUT / PATCH).
			if cleanMethod == "DELETE" || cleanMethod == "PUT" {
				return false, "BLOCKED_BACEN_538_MUTATION_DENIED"
			}

		case ModeLGPD:
			// LGPD enforces data minimization and prevents bulk data extraction on sensitive endpoints.
			if cleanMethod == "GET" && (strings.Contains(cleanPath, "/clientes") || strings.Contains(cleanPath, "/customers") || strings.Contains(cleanPath, "/users")) {
				return false, "BLOCKED_LGPD_BULK_DATA_ACCESS_DENIED"
			}

		case ModeCFM:
			// CFM restricts access to medical records without homologated human oversight.
			if strings.Contains(cleanPath, "/prontuarios") || strings.Contains(cleanPath, "/medical-records") {
				return false, "BLOCKED_CFM_MEDICAL_RECORDS_DENIED"
			}
		}
	}

	return true, "ALLOWED"
}
