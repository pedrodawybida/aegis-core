// Package proxy provides the high-performance HTTP reverse proxy engine
// that enforces compliance policies before forwarding traffic.
package proxy

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/pedrodawybida/aegis-core/internal/audit"
	"github.com/pedrodawybida/aegis-core/internal/compliance"
)

// AegisProxy is a specialized HTTP handler that wraps an httputil.ReverseProxy.
// It injects identity verification, compliance checking, and immutable audit logging.
type AegisProxy struct {
	targetURL *url.URL
	reverse   *httputil.ReverseProxy
	logger    *audit.Logger
	policyDB  map[string]compliance.AgentPolicy
	dryRun    bool
}

// NewAegisProxy initializes a new proxy instance pointing to the target URL.
// It requires an audit logger and a policy map for evaluation.
func NewAegisProxy(target string, logger *audit.Logger, policyDB map[string]compliance.AgentPolicy) (*AegisProxy, error) {
	tURL, err := url.Parse(target)
	if err != nil {
		return nil, err
	}

	rp := httputil.NewSingleHostReverseProxy(tURL)
	rp.Director = func(req *http.Request) {
		req.URL.Scheme = tURL.Scheme
		req.URL.Host = tURL.Host
		req.Host = tURL.Host
	}

	return &AegisProxy{
		targetURL: tURL,
		reverse:   rp,
		logger:    logger,
		policyDB:  policyDB,
		dryRun:    false,
	}, nil
}

// SetDryRun enables or disables Dry-Run (Shadow / Audit-Only) Mode.
func (p *AegisProxy) SetDryRun(dryRun bool) {
	p.dryRun = dryRun
}

// ServeHTTP implements the http.Handler interface. It executes the core security pipeline:
// Identity Check -> Policy Evaluation -> Audit Logging -> Forwarding.
func (p *AegisProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Handle Aegis System Endpoints
	if r.URL.Path == "/_aegis/health" || r.URL.Path == "/healthz" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"status":"ok","service":"aegis-core","target_api":"%s","active_agents":%d,"dry_run":%t}`, p.targetURL.String(), len(p.policyDB), p.dryRun)))
		return
	}

	if r.URL.Path == "/_aegis/dashboard" || r.URL.Path == "/_aegis/console" {
		p.serveDashboardHTML(w, r)
		return
	}

	if r.URL.Path == "/_aegis/api/agents" {
		p.serveAgentsAPI(w, r)
		return
	}

	w.Header().Set("X-Aegis-Proxy", "aegis-core/v1.0")

	// Extract Agent Token from Authorization header or X-Aegis-Agent-Id
	agentToken := r.Header.Get("X-Aegis-Agent-Id")
	if agentToken == "" {
		authHeader := r.Header.Get("Authorization")
		if strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
			agentToken = strings.TrimSpace(authHeader[7:])
		} else {
			agentToken = authHeader
		}
	}

	ip := r.RemoteAddr
	action := fmt.Sprintf("%s %s", r.Method, r.URL.Path)

	var payload string
	if r.Body != nil {
		bodyBytes, _ := io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		payload = string(bodyBytes)
	}

	// 1. Identity Check
	if agentToken == "" {
		p.logger.LogAction("UNKNOWN", action, payload, "BLOCKED_NO_IDENTITY", ip)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Aegis-Compliance-Status", "BLOCKED_NO_IDENTITY")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "Aegis: Missing Identity Token. Provide Bearer token or X-Aegis-Agent-Id"}`))
		return
	}

	// 2. Policy Lookup
	policy, exists := p.policyDB[agentToken]
	if !exists {
		p.logger.LogAction(agentToken, action, payload, "BLOCKED_AGENT_NOT_FOUND", ip)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Aegis-Compliance-Status", "BLOCKED_AGENT_NOT_FOUND")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(fmt.Sprintf(`{"error": "Aegis: Agent identity '%s' not recognized"}`, agentToken)))
		return
	}

	// 3. Compliance Verification
	allowed, reason := compliance.EvaluateRequest(policy, r.Method, r.URL.Path)

	if p.dryRun && !allowed {
		dryRunReason := fmt.Sprintf("DRY_RUN_%s", reason)
		p.logger.LogAction(agentToken, action, payload, dryRunReason, ip)
		w.Header().Set("X-Aegis-Compliance-Status", dryRunReason)
		w.Header().Set("X-Aegis-Dry-Run", "true")
		p.reverse.ServeHTTP(w, r)
		return
	}

	// 4. Immutable Audit Log
	p.logger.LogAction(agentToken, action, payload, reason, ip)

	w.Header().Set("X-Aegis-Compliance-Status", reason)

	if !allowed {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(fmt.Sprintf(`{"error": "Aegis Compliance Violation: %s"}`, reason)))
		return
	}

	// 5. Safe Forwarding
	p.reverse.ServeHTTP(w, r)
}
