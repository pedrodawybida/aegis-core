package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pedrodawybida/nexo-hub/internal/audit"
	"github.com/pedrodawybida/nexo-hub/internal/config"
	"github.com/pedrodawybida/nexo-hub/internal/proxy"
)

func getEnvOrDefault(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func main() {
	// Parse CLI flags with Environment Variable fallbacks (NEXO_ with AEGIS_ fallback)
	defaultConfig := getEnvOrDefault("NEXO_CONFIG", getEnvOrDefault("AEGIS_CONFIG", "nexo.yaml"))
	defaultPort := getEnvOrDefault("NEXO_PORT", getEnvOrDefault("AEGIS_PORT", "8080"))
	defaultLog := getEnvOrDefault("NEXO_LOG_FILE", getEnvOrDefault("AEGIS_LOG_FILE", "audit_bacen.log"))

	configPathFlag := flag.String("config", defaultConfig, "Path to nexo configuration file")
	portFlag := flag.String("port", defaultPort, "Port for Nexo Hub proxy server")
	logFileFlag := flag.String("log", defaultLog, "Path to audit log file")
	dryRunFlag := flag.Bool("dry-run", false, "Enable Dry-Run (Shadow / Audit-Only) Mode without blocking requests")
	flag.Parse()

	// 1. Load the dynamic configuration
	cfg, err := config.LoadConfig(*configPathFlag)
	if err != nil {
		log.Fatalf("Fatal: Failed to read config file '%s': %v", *configPathFlag, err)
	}

	if *dryRunFlag {
		cfg.DryRun = true
	}

	// 2. Initialize the Immutable Audit Logger
	auditLogger, err := audit.NewLogger(*logFileFlag)
	if err != nil {
		log.Fatalf("Fatal: Failed to initialize Audit Logger at '%s': %v", *logFileFlag, err)
	}

	// 3. Initialize the Core Proxy Engine
	policyMap := cfg.ToPolicyMap()
	agProxy, err := proxy.NewAegisProxy(cfg.TargetAPI, auditLogger, policyMap)
	if err != nil {
		log.Fatalf("Fatal: Failed to initialize Nexo Proxy: %v", err)
	}
	agProxy.SetDryRun(cfg.DryRun)

	listenAddr := fmt.Sprintf(":%s", *portFlag)
	server := &http.Server{
		Addr:         listenAddr,
		Handler:      agProxy,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Channel to listen for OS interrupt signals for Graceful Shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	go func() {
		fmt.Println("===================================================")
		fmt.Println("🛡️  Nexo Hub (BR Compliance Edition)")
		fmt.Printf("🎯 Protecting API: %s\n", cfg.TargetAPI)
		fmt.Printf("🤖 Agents Loaded: %d\n", len(policyMap))
		fmt.Printf("📝 Immutable audit logs active at: %s\n", *logFileFlag)
		if cfg.DryRun {
			fmt.Println("⚠️  Dry-Run (Shadow / Audit-Only) Mode ACTIVE")
		}
		fmt.Printf("🚀 Nexo Hub listening on %s...\n", listenAddr)
		fmt.Println("===================================================")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Fatal: Proxy Server crashed: %v", err)
		}
	}()

	// Wait for OS shutdown signal
	sig := <-shutdown
	log.Printf("[NEXO-SYSTEM] Shutting down Nexo Hub (Signal: %v)...", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("[NEXO-SYSTEM-ERROR] HTTP Server Graceful Shutdown failed: %v", err)
	}

	if err := auditLogger.Close(); err != nil {
		log.Printf("[NEXO-SYSTEM-ERROR] Failed to close audit logger cleanly: %v", err)
	}

	log.Println("[NEXO-SYSTEM] Nexo Hub stopped cleanly.")
}
