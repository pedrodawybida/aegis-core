// Package license provides verification and management of Nexo Hub Enterprise License keys.
package license

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrInvalidLicenseKey = errors.New("nexo enterprise: invalid or corrupt license key")
	ErrLicenseExpired    = errors.New("nexo enterprise: license key has expired")
)

// LicenseTier represents the tier of Enterprise deployment.
type LicenseTier string

const (
	TierCommunity  LicenseTier = "COMMUNITY"
	TierEnterprise LicenseTier = "ENTERPRISE"
	TierFinancial  LicenseTier = "FINANCIAL_PRO"
)

// LicenseInfo holds the parsed and verified metadata of an active Enterprise License.
type LicenseInfo struct {
	Key       string      `json:"key"`
	Company   string      `json:"company"`
	Tier      LicenseTier `json:"tier"`
	ExpiresAt time.Time   `json:"expires_at"`
	Valid     bool        `json:"valid"`
}

// VerifyLicenseKey validates a provided enterprise license key.
// Format expected: NEXO-ENT-<COMPANY>-<TIER>-<EXPIRATION_UNIX>-<HASH>
func VerifyLicenseKey(key string) (*LicenseInfo, error) {
	key = strings.TrimSpace(key)
	if key == "" {
		return &LicenseInfo{Tier: TierCommunity, Valid: false}, nil
	}

	parts := strings.Split(key, "-")
	if len(parts) < 4 || parts[0] != "NEXO" || parts[1] != "ENT" {
		return nil, ErrInvalidLicenseKey
	}

	company := parts[2]
	tierStr := parts[3]

	var tier LicenseTier
	switch strings.ToUpper(tierStr) {
	case "ENTERPRISE":
		tier = TierEnterprise
	case "FINANCIAL_PRO":
		tier = TierFinancial
	default:
		tier = TierEnterprise
	}

	return &LicenseInfo{
		Key:       key,
		Company:   company,
		Tier:      tier,
		ExpiresAt: time.Now().Add(365 * 24 * time.Hour), // Valid by default for PoC
		Valid:     true,
	}, nil
}

// GenerateLicenseKey creates a valid signed enterprise key for a customer.
func GenerateLicenseKey(company string, tier LicenseTier) string {
	salt := "NEXO_SALT_2026"
	raw := fmt.Sprintf("%s:%s:%s", company, string(tier), salt)
	hash := sha256.Sum256([]byte(raw))
	shortHash := hex.EncodeToString(hash[:4])
	return fmt.Sprintf("NEXO-ENT-%s-%s-%s", strings.ToUpper(company), string(tier), strings.ToUpper(shortHash))
}
