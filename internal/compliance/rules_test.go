package compliance

import "testing"

func TestEvaluateRequest(t *testing.T) {
	fintechPolicy := AgentPolicy{
		AgentID:     "ia-fintech-support",
		ActiveModes: []Mode{ModeLGPD, ModeBacen538},
	}

	healthPolicy := AgentPolicy{
		AgentID:     "ia-health-bot",
		ActiveModes: []Mode{ModeCFM},
	}

	tests := []struct {
		name           string
		policy         AgentPolicy
		method         string
		path           string
		expectedAllow  bool
		expectedReason string
	}{
		{
			name:           "LGPD Block Bulk Customers",
			policy:         fintechPolicy,
			method:         "GET",
			path:           "/clientes",
			expectedAllow:  false,
			expectedReason: "BLOCKED_LGPD_BULK_DATA_ACCESS_DENIED",
		},
		{
			name:           "BACEN 538 Block DELETE",
			policy:         fintechPolicy,
			method:         "DELETE",
			path:           "/transacoes/123",
			expectedAllow:  false,
			expectedReason: "BLOCKED_BACEN_538_MUTATION_DENIED",
		},
		{
			name:           "BACEN 538 Block PUT",
			policy:         fintechPolicy,
			method:         "PUT",
			path:           "/saldo",
			expectedAllow:  false,
			expectedReason: "BLOCKED_BACEN_538_MUTATION_DENIED",
		},
		{
			name:           "Fintech Allow Safe POST",
			policy:         fintechPolicy,
			method:         "POST",
			path:           "/pix/validar",
			expectedAllow:  true,
			expectedReason: "ALLOWED",
		},
		{
			name:           "CFM Block Medical Records",
			policy:         healthPolicy,
			method:         "GET",
			path:           "/prontuarios/999",
			expectedAllow:  false,
			expectedReason: "BLOCKED_CFM_MEDICAL_RECORDS_DENIED",
		},
		{
			name:           "CFM Allow Appointment Lookup",
			policy:         healthPolicy,
			method:         "GET",
			path:           "/consultas/horarios",
			expectedAllow:  true,
			expectedReason: "ALLOWED",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			allowed, reason := EvaluateRequest(tt.policy, tt.method, tt.path)
			if allowed != tt.expectedAllow {
				t.Errorf("Expected allowed=%v, got %v", tt.expectedAllow, allowed)
			}
			if reason != tt.expectedReason {
				t.Errorf("Expected reason=%s, got %s", tt.expectedReason, reason)
			}
		})
	}
}
