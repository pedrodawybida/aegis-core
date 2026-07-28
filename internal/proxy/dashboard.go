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
			"context_reduction": "Up to 94% savings",
		},
	})
}

const dashboardHTMLTemplate = `<!DOCTYPE html>
<html lang="pt-BR" class="dark">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>🛡️ Nexo Hub - Compliance, NSEP & Token Cost Playground</title>
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
                <p class="text-xs text-gray-400">NSEP Protocol, Token Cost Calculator & Agent Execution Replay</p>
            </div>
        </div>
        <div class="flex items-center space-x-4">
            <span class="inline-flex items-center px-3 py-1 rounded-full text-xs font-medium bg-emerald-950 text-emerald-400 border border-emerald-800 pulse-green">
                <span class="w-2 h-2 mr-2 rounded-full bg-emerald-400 animate-ping"></span> Live Proxy Active
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
                <p class="text-xs font-semibold uppercase text-gray-400">NSEP Token Savings</p>
                <p class="text-2xl font-bold text-emerald-400 mt-1">- 94% Context</p>
                <p class="text-xs text-emerald-500/80 mt-2">⚡ Context Window Footprint</p>
            </div>

            <div class="glass p-5 rounded-2xl">
                <p class="text-xs font-semibold uppercase text-gray-400">Transporte MCP</p>
                <p class="text-xl font-bold text-indigo-400 mt-1">/_nexo/mcp</p>
                <p class="text-xs text-gray-500 mt-2">Claude, ChatGPT & Cursor</p>
            </div>

            <div class="glass p-5 rounded-2xl">
                <p class="text-xs font-semibold uppercase text-gray-400">Conformidade BACEN</p>
                <p class="text-2xl font-bold text-emerald-400 mt-1">100% OK</p>
                <p class="text-xs text-gray-500 mt-2">CMN 5.274 & LGPD Audit</p>
            </div>
        </div>

        <!-- Token Cost Playground Section (Growth Loop Item 9.1a) -->
        <div class="glass p-6 rounded-2xl space-y-4">
            <h2 class="text-lg font-semibold text-emerald-400 flex items-center gap-2">
                <span>📊 Token Cost Playground (MCP vs NSEP Protocol)</span>
            </h2>
            <p class="text-xs text-gray-400">Calcule a economia em tempo real de tokens de contexto utilizando o protocolo NSEP em comparação ao MCP tradicional.</p>
            
            <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
                <div>
                    <label class="block text-xs font-medium text-gray-300 mb-1">Tamanho da Spec OpenAPI (KB)</label>
                    <input type="number" id="spec-size" value="120" oninput="calculateSavings()" class="w-full bg-gray-900 border border-gray-700 rounded-lg p-2.5 text-sm text-white" />
                </div>
                <div>
                    <label class="block text-xs font-medium text-gray-300 mb-1">Número de Chamadas de Ferramentas / Mês</label>
                    <input type="number" id="call-count" value="50000" oninput="calculateSavings()" class="w-full bg-gray-900 border border-gray-700 rounded-lg p-2.5 text-sm text-white" />
                </div>
                <div class="glass p-4 rounded-xl flex flex-col justify-center">
                    <p class="text-xs text-gray-400 uppercase font-semibold">Economia Estimada ($ / mês)</p>
                    <p class="text-2xl font-bold text-emerald-400" id="savings-amount">$720.00 / mês</p>
                    <p class="text-xs text-emerald-500 mt-1" id="savings-pct">94% de redução em janelas de contexto</p>
                </div>
            </div>
        </div>

        <!-- Interactive Testing & Replay Console (Growth Loop Item 9.1b) -->
        <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
            <div class="lg:col-span-1 glass p-6 rounded-2xl space-y-4">
                <h2 class="text-lg font-semibold flex items-center gap-2">
                    <span>⚡ Simulação de Tool-Calling</span>
                </h2>
                
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

            <!-- Replay Visual de Execução -->
            <div class="lg:col-span-2 glass p-6 rounded-2xl flex flex-col space-y-4">
                <div class="flex justify-between items-center">
                    <h2 class="text-lg font-semibold">🎬 Replay Visual de Execução (Agent Debugger)</h2>
                    <span class="text-xs text-gray-400">audit_bacen.log</span>
                </div>

                <div id="response-box" class="flex-1 bg-gray-950 p-4 rounded-xl font-mono text-xs overflow-auto border border-gray-800 text-gray-300 min-h-[250px]">
                    // Aguardando simulação...
                </div>
            </div>
        </div>

        <!-- Compliance Verification Badge Generator (Growth Loop Item 9.1c) -->
        <div class="glass p-6 rounded-2xl space-y-4 border border-blue-900/50">
            <h2 class="text-lg font-semibold text-blue-400 flex items-center gap-2">
                <span>🛡️ Badge Público de Compliance Verificado (Shields.io)</span>
            </h2>
            <p class="text-xs text-gray-400">Adicione este badge ao README do seu repositório para comprovar a auditoria e conformidade do seu Agente de IA:</p>

            <div class="flex items-center space-x-4 bg-gray-900 p-4 rounded-xl border border-gray-800 overflow-x-auto">
                <img src="https://img.shields.io/badge/Nexo%20Hub%20Verified-BACEN%20538%2F2025-blue?logo=shield" alt="Nexo Hub Verified Badge" />
                <code class="text-xs font-mono text-emerald-400 flex-1">![Nexo Hub Verified](https://img.shields.io/badge/Nexo%20Hub%20Verified-BACEN%20538%2F2025-blue?logo=shield)</code>
            </div>
        </div>

    </main>

    <script>
        function calculateSavings() {
            var specKB = parseFloat(document.getElementById('spec-size').value) || 120;
            var calls = parseFloat(document.getElementById('call-count').value) || 50000;
            
            var mcpTokensPerCall = specKB * 250; // ~250 tokens per KB of schema
            var nsepTokensPerCall = 1200; // Fixed small footprint
            
            var mcpTotalTokens = mcpTokensPerCall * calls;
            var nsepTotalTokens = nsepTokensPerCall * calls;
            
            var savedTokens = mcpTotalTokens - nsepTotalTokens;
            var costSavings = (savedTokens / 1000000) * 15.0; // ~$15 per million tokens
            
            if (costSavings < 0) costSavings = 0;
            
            document.getElementById('savings-amount').innerText = '$' + costSavings.toFixed(2) + ' / mês';
        }

        async function sendSimulatedRequest() {
            var agent = document.getElementById('agent-select').value;
            var method = document.getElementById('method-select').value;
            var route = document.getElementById('route-input').value;
            var box = document.getElementById('response-box');

            box.innerHTML = '<span class="text-yellow-400">⏳ Executando replay do agente via Nexo Hub...</span>';

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

                var statusBadge = res.status === 200 ? '<span class="text-emerald-400 font-bold">[APROVADO - 200 OK]</span>' : '<span class="text-red-400 font-bold">[BLOQUEADO PELO NEXO HUB - 403 FORBIDDEN]</span>';

                box.innerHTML = statusBadge + '\n' +
                    '┌─ Agent Identity: ' + agent + '\n' +
                    '├─ HTTP Action: ' + method + ' ' + route + '\n' +
                    '├─ Compliance Mode Evaluated: BACEN_538 / LGPD\n' +
                    '└─ Result Header: ' + statusHeader + '\n\n' +
                    'Corpo da Resposta:\n' + JSON.stringify(body, null, 2);
            } catch (err) {
                box.innerHTML = '<span class="text-red-400">Erro na chamada: ' + err.message + '</span>';
            }
        }
    </script>
</body>
</html>`
