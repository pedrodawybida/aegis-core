package proxy

import (
	"encoding/json"
	"net/http"
)

// serveDashboardHTML renders the Nexo Hub Web UI Console for live inspection and testing.
func (p *AegisProxy) serveDashboardHTML(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(dashboardHTMLTemplate))
}

// serveAgentsAPI returns JSON of configured agents and active compliance modes.
func (p *AegisProxy) serveAgentsAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"target_api": p.targetURL.String(),
		"agents":     p.policyDB,
		"nsep": map[string]interface{}{
			"status":            "ACTIVE",
			"js_engine":         "goja",
			"mcp_endpoint":      "/_nexo/mcp",
			"context_reduction": "Up to 90% savings",
		},
	})
}

const dashboardHTMLTemplate = `<!DOCTYPE html>
<html lang="pt-BR" class="dark">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>🛡️ Nexo Hub - Compliance & NSEP Protocol Console</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700&display=swap" rel="stylesheet">
    <style>
        body { font-family: 'Inter', sans-serif; background-color: #0b0f19; color: #f3f4f6; }
        .glass { background: rgba(17, 24, 39, 0.7); backdrop-filter: blur(12px); border: 1px solid rgba(255, 255, 255, 0.08); }
        .pulse-green { box-shadow: 0 0 15px rgba(16, 185, 129, 0.4); }
    </style>
</head>
<body class="min-h-screen flex flex-col antialiased">
    <!-- Header -->
    <header class="glass sticky top-0 z-50 px-6 py-4 flex items-center justify-between border-b border-gray-800">
        <div class="flex items-center space-x-3">
            <span class="text-3xl">🛡️</span>
            <div>
                <h1 class="text-xl font-bold bg-gradient-to-r from-blue-400 to-indigo-400 bg-clip-text text-transparent">Nexo Hub</h1>
                <p class="text-xs text-gray-400">Home of NSEP (Nexo Secure Execution Protocol) & BACEN 538 Compliance Engine</p>
            </div>
        </div>
        <div class="flex items-center space-x-4">
            <span class="inline-flex items-center px-3 py-1 rounded-full text-xs font-medium bg-emerald-950 text-emerald-400 border border-emerald-800 pulse-green">
                <span class="w-2 h-2 mr-2 rounded-full bg-emerald-400 animate-ping"></span> Live Proxy & NSEP Active
            </span>
            <a href="https://github.com/pedrodawybida/nexo-hub" target="_blank" class="text-xs text-gray-400 hover:text-white transition">GitHub Repo ↗</a>
        </div>
    </header>

    <!-- Main Content -->
    <main class="flex-1 max-w-7xl w-full mx-auto p-6 space-y-8">

        <!-- Stat Cards -->
        <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
            <div class="glass p-5 rounded-2xl">
                <p class="text-xs font-semibold uppercase text-gray-400">API Protegida</p>
                <p class="text-lg font-bold text-blue-400 mt-1 truncate" id="target-api">http://localhost:9000</p>
                <p class="text-xs text-gray-500 mt-2">Target Backend Protocol</p>
            </div>

            <div class="glass p-5 rounded-2xl">
                <p class="text-xs font-semibold uppercase text-gray-400">NSEP Protocol</p>
                <p class="text-xl font-bold text-emerald-400 mt-1">90% Context Savings</p>
                <p class="text-xs text-emerald-500/80 mt-2">⚡ Sandbox Goja (0% CGO)</p>
            </div>

            <div class="glass p-5 rounded-2xl">
                <p class="text-xs font-semibold uppercase text-gray-400">Transporte MCP</p>
                <p class="text-xl font-bold text-indigo-400 mt-1">/_nexo/mcp</p>
                <p class="text-xs text-gray-500 mt-2">Claude, ChatGPT & Cursor Ready</p>
            </div>

            <div class="glass p-5 rounded-2xl">
                <p class="text-xs font-semibold uppercase text-gray-400">Conformidade BACEN</p>
                <p class="text-2xl font-bold text-emerald-400 mt-1">100% OK</p>
                <p class="text-xs text-gray-500 mt-2">CMN 5.274 & LGPD Audit</p>
            </div>
        </div>

        <!-- Interactive Testing Console -->
        <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
            <div class="lg:col-span-1 glass p-6 rounded-2xl space-y-4">
                <h2 class="text-lg font-semibold flex items-center gap-2">
                    <span>⚡ Simulação de Tool-Calling & NSEP</span>
                </h2>
                <p class="text-xs text-gray-400">Simule chamadas de um Agente de IA para testar as regras de compliance e execução do NSEP.</p>
                
                <div>
                    <label class="block text-xs font-medium text-gray-300 mb-1">Selecione o Agente de IA</label>
                    <select id="agent-select" class="w-full bg-gray-900 border border-gray-700 rounded-lg p-2.5 text-sm text-white focus:ring-2 focus:ring-blue-500">
                        <option value="ia-fintech-support">ia-fintech-support (LGPD + BACEN 538)</option>
                        <option value="ia-health-bot">ia-health-bot (CFM)</option>
                    </select>
                </div>

                <div>
                    <label class="block text-xs font-medium text-gray-300 mb-1">Método HTTP</label>
                    <select id="method-select" class="w-full bg-gray-900 border border-gray-700 rounded-lg p-2.5 text-sm text-white focus:ring-2 focus:ring-blue-500">
                        <option value="POST">POST (Permitido em /pix/validar)</option>
                        <option value="GET">GET (Bloqueado em /clientes)</option>
                        <option value="DELETE">DELETE (Bloqueado em /transacoes/99)</option>
                    </select>
                </div>

                <div>
                    <label class="block text-xs font-medium text-gray-300 mb-1">Rota da API</label>
                    <input type="text" id="route-input" value="/pix/validar" class="w-full bg-gray-900 border border-gray-700 rounded-lg p-2.5 text-sm text-white focus:ring-2 focus:ring-blue-500" />
                </div>

                <button onclick="sendSimulatedRequest()" class="w-full py-3 px-4 bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-500 hover:to-indigo-500 text-white font-semibold rounded-xl shadow-lg transition duration-200">
                    Enviar Requisição via Nexo Hub
                </button>
            </div>

            <!-- Log & Response Panel -->
            <div class="lg:col-span-2 glass p-6 rounded-2xl flex flex-col space-y-4">
                <div class="flex justify-between items-center">
                    <h2 class="text-lg font-semibold">📜 Replay de Execução & Logs de Auditoria</h2>
                    <span class="text-xs text-gray-400">audit_bacen.log</span>
                </div>

                <div id="response-box" class="flex-1 bg-gray-950 p-4 rounded-xl font-mono text-xs overflow-auto border border-gray-800 text-gray-300 min-h-[250px]">
                    // Aguardando simulação...
                </div>
            </div>
        </div>
    </main>

    <script>
        async function sendSimulatedRequest() {
            var agent = document.getElementById('agent-select').value;
            var method = document.getElementById('method-select').value;
            var route = document.getElementById('route-input').value;
            var box = document.getElementById('response-box');

            box.innerHTML = '<span class="text-yellow-400">⏳ Enviando requisição através do Nexo Hub...</span>';

            try {
                var res = await fetch(route, {
                    method: method,
                    headers: {
                        'Authorization': 'Bearer ' + agent,
                        'Content-Type': 'application/json'
                    }
                });

                var statusHeader = res.headers.get('X-Nexo-Compliance-Status') || res.headers.get('X-Aegis-Compliance-Status') || 'UNKNOWN';
                var body = await res.json().catch(function() { return { message: 'Sem corpo na resposta' }; });

                var statusBadge = res.status === 200 ? '<span class="text-emerald-400 font-bold">[APROVADO]</span>' : '<span class="text-red-400 font-bold">[BLOQUEADO PELO NEXO HUB]</span>';

                box.innerHTML = statusBadge + ' HTTP Status: ' + res.status + '\n' +
                    'Status de Conformidade: ' + statusHeader + '\n' +
                    'Identidade do Agente: ' + agent + '\n\n' +
                    'Corpo da Resposta:\n' + JSON.stringify(body, null, 2);
            } catch (err) {
                box.innerHTML = '<span class="text-red-400">Erro na chamada: ' + err.message + '</span>';
            }
        }
    </script>
</body>
</html>`
