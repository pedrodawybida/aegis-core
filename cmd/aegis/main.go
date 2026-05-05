// Package main is the entrypoint for the Aegis Core application.
// It initializes the configuration, the audit logger, and starts the reverse proxy server.
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pedrodawybida/aegis-core/internal/audit"
	"github.com/pedrodawybida/aegis-core/internal/config"
	"github.com/pedrodawybida/aegis-core/internal/proxy"
)

func main() {
	// 1. Load the dynamic configuration
	cfg, err := config.LoadConfig("aegis.yaml")
	if err != nil {
		log.Fatalf("Fatal: Failed to read aegis.yaml: %v", err)
	}

	// 2. Initialize the Immutable Audit Logger
	auditLogger, err := audit.NewLogger("audit_bacen.log")
	if err != nil {
		log.Fatalf("Fatal: Failed to initialize Audit Logger: %v", err)
	}
	defer auditLogger.Close()

	// 3. Initialize the Core Proxy Engine
	policyMap := cfg.ToPolicyMap()
	agProxy, err := proxy.NewAegisProxy(cfg.TargetAPI, auditLogger, policyMap)
	if err != nil {
		log.Fatalf("Fatal: Failed to initialize Aegis Proxy: %v", err)
	}

	fmt.Println("===================================================")
	fmt.Println("🛡️  Aegis Core (BR Compliance Edition)")
	fmt.Printf("🎯 Protecting API: %s\n", cfg.TargetAPI)
	fmt.Printf("🤖 Agents Loaded: %d\n", len(policyMap))
	fmt.Println("📝 Immutable logs active at: audit_bacen.log")
	fmt.Println("🚀 Aegis listening on :8080...")
	fmt.Println("===================================================")

	if err := http.ListenAndServe(":8080", agProxy); err != nil {
		log.Fatalf("Fatal: Proxy Server crashed: %v", err)
	}
}
