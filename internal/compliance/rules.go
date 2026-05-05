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
	for _, mode := range policy.ActiveModes {
		switch mode {
		
		case ModeBacen538:
			// BACEN 538 mandates strict control over state mutations.
			if method == "DELETE" || method == "PUT" {
				return false, "BLOCKED_BACEN_538_MUTATION_DENIED"
			}
			
		case ModeLGPD:
			// LGPD enforces the principle of data minimization and prevents bulk extraction.
			if method == "GET" && strings.Contains(path, "/clientes") {
				return false, "BLOCKED_LGPD_BULK_DATA_ACCESS_DENIED"
			}
			
		case ModeCFM:
			// CFM restricts access to medical records without homologated human oversight.
			if strings.Contains(path, "/prontuarios") {
				return false, "BLOCKED_CFM_MEDICAL_RECORDS_DENIED"
			}
		}
	}
	
	return true, "ALLOWED"
}
