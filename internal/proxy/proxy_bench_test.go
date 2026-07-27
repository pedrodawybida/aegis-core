package proxy

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/pedrodawybida/aegis-core/internal/audit"
	"github.com/pedrodawybida/aegis-core/internal/compliance"
)

// BenchmarkAegisProxy Latency Test
func BenchmarkAegisProxy(b *testing.B) {
	mockBackend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}))
	defer mockBackend.Close()

	tmpDir := b.TempDir()
	logPath := filepath.Join(tmpDir, "bench_audit.log")
	logger, _ := audit.NewLogger(logPath)
	defer logger.Close()

	policyDB := map[string]compliance.AgentPolicy{
		"bench-agent": {
			AgentID:     "bench-agent",
			ActiveModes: []compliance.Mode{compliance.ModeLGPD, compliance.ModeBacen538},
		},
	}

	agProxy, _ := NewAegisProxy(mockBackend.URL, logger, policyDB)
	req, _ := http.NewRequest("POST", "/api/v1/safe-action", nil)
	req.Header.Set("Authorization", "Bearer bench-agent")

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		agProxy.ServeHTTP(rr, req)
	}
}
