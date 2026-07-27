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

	"github.com/pedrodawybida/aegis-core/internal/audit"
	"github.com/pedrodawybida/aegis-core/internal/config"
	"github.com/pedrodawybida/aegis-core/internal/proxy"
)

func getEnvOrDefault(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func main() {
	// Parse CLI flags with Environment Variable fallbacks
	configPathFlag := flag.String("config", getEnvOrDefault("AEGIS_CONFIG", "aegis.yaml"), "Path to aegis configuration file")
	portFlag := flag.String("port", getEnvOrDefault("AEGIS_PORT", "8080"), "Port for Aegis Core proxy server")
	logFileFlag := flag.String("log", getEnvOrDefault("AEGIS_LOG_FILE", "audit_bacen.log"), "Path to audit log file")
	flag.Parse()

	// 1. Load the dynamic configuration
	cfg, err := config.LoadConfig(*configPathFlag)
	if err != nil {
		log.Fatalf("Fatal: Failed to read config file '%s': %v", *configPathFlag, err)
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
		log.Fatalf("Fatal: Failed to initialize Aegis Proxy: %v", err)
	}

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
		fmt.Println("🛡️  Aegis Core (BR Compliance Edition)")
		fmt.Printf("🎯 Protecting API: %s\n", cfg.TargetAPI)
		fmt.Printf("🤖 Agents Loaded: %d\n", len(policyMap))
		fmt.Printf("📝 Immutable audit logs active at: %s\n", *logFileFlag)
		fmt.Printf("🚀 Aegis listening on %s...\n", listenAddr)
		fmt.Println("===================================================")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Fatal: Proxy Server crashed: %v", err)
		}
	}()

	// Wait for OS shutdown signal
	sig := <-shutdown
	log.Printf("[AEGIS-SYSTEM] Shutting down Aegis Core (Signal: %v)...", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("[AEGIS-SYSTEM-ERROR] HTTP Server Graceful Shutdown failed: %v", err)
	}

	if err := auditLogger.Close(); err != nil {
		log.Printf("[AEGIS-SYSTEM-ERROR] Failed to close audit logger cleanly: %v", err)
	}

	log.Println("[AEGIS-SYSTEM] Aegis Core stopped cleanly.")
}
