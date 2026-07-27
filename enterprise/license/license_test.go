package license

import (
	"testing"
)

func TestLicenseVerification(t *testing.T) {
	t.Run("Empty License Key defaults to Community Tier", func(t *testing.T) {
		info, err := VerifyLicenseKey("")
		if err != nil {
			t.Fatalf("Expected no error for empty key, got %v", err)
		}
		if info.Tier != TierCommunity {
			t.Errorf("Expected TierCommunity, got %s", info.Tier)
		}
		if info.Valid {
			t.Errorf("Expected Valid to be false for community tier")
		}
	})

	t.Run("Invalid Key Format returns error", func(t *testing.T) {
		_, err := VerifyLicenseKey("INVALID-KEY-FORMAT")
		if err != ErrInvalidLicenseKey {
			t.Errorf("Expected ErrInvalidLicenseKey, got %v", err)
		}
	})

	t.Run("Valid Enterprise License Generation & Verification", func(t *testing.T) {
		key := GenerateLicenseKey("NUBANK", TierFinancial)
		info, err := VerifyLicenseKey(key)
		if err != nil {
			t.Fatalf("Failed to verify valid key: %v", err)
		}
		if !info.Valid {
			t.Errorf("Expected license to be valid")
		}
		if info.Company != "NUBANK" {
			t.Errorf("Expected company NUBANK, got %s", info.Company)
		}
		if info.Tier != TierFinancial {
			t.Errorf("Expected TierFinancial, got %s", info.Tier)
		}
	})
}
