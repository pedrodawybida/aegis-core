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
	}, nil
}

// ServeHTTP implements the http.Handler interface. It executes the core security pipeline:
// Identity Check -> Policy Evaluation -> Audit Logging -> Forwarding.
func (p *AegisProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	agentToken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
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
		http.Error(w, `{"error": "Aegis: Missing Identity Token"}`, http.StatusUnauthorized)
		return
	}

	// 2. Policy Lookup
	policy, exists := p.policyDB[agentToken]
	if !exists {
		p.logger.LogAction(agentToken, action, payload, "BLOCKED_AGENT_NOT_FOUND", ip)
		http.Error(w, `{"error": "Aegis: Agent identity not recognized"}`, http.StatusForbidden)
		return
	}

	// 3. Compliance Verification
	allowed, reason := compliance.EvaluateRequest(policy, r.Method, r.URL.Path)
	
	// 4. Immutable Audit
	p.logger.LogAction(agentToken, action, payload, reason, ip)

	if !allowed {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(fmt.Sprintf(`{"error": "Aegis Compliance Violation: %s"}`, reason)))
		return
	}

	// 5. Safe Forwarding
	p.reverse.ServeHTTP(w, r)
}
