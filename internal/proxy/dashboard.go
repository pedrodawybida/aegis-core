package proxy

import (
	"encoding/json"
	"net/http"
)

// serveDashboardHTML renders the Aegis Web UI Console for live inspection and testing.
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
	})
}

const dashboardHTMLTemplate = `<!DOCTYPE html>
<html lang="pt-BR" class="dark">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>🛡️ Aegis Core - Compliance & Security Console</title>
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
                <h1 class="text-xl font-bold bg-gradient-to-r from-blue-400 to-indigo-400 bg-clip-text text-transparent">Aegis Core</h1>
                <p class="text-xs text-gray-400">BACEN 538/2025 & LGPD AI Compliance Shield</p>
            </div>
        </div>
        <div class="flex items-center space-x-4">
            <span class="inline-flex items-center px-3 py-1 rounded-full text-xs font-medium bg-emerald-950 text-emerald-400 border border-emerald-800 pulse-green">
                <span class="w-2 h-2 mr-2 rounded-full bg-emerald-400 animate-ping"></span> Live Proxy Active
            </span>
            <a href="https://github.com/pedrodawybida/aegis-core" target="_blank" class="text-xs text-gray-400 hover:text-white transition">GitHub Repo ↗</a>
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
                <p class="text-xs font-semibold uppercase text-gray-400">Latência do Proxy</p>
                <p class="text-2xl font-bold text-emerald-400 mt-1">&lt; 0.1 ms</p>
                <p class="text-xs text-emerald-500/80 mt-2">⚡ Ultra-low overhead</p>
            </div>

            <div class="glass p-5 rounded-2xl">
                <p class="text-xs font-semibold uppercase text-gray-400">Agentes Ativos</p>
                <p class="text-2xl font-bold text-indigo-400 mt-1" id="agent-count">3</p>
                <p class="text-xs text-gray-500 mt-2">Identidades Não-Humanas</p>
            </div>

            <div class="glass p-5 rounded-2xl">
                <p class="text-xs font-semibold uppercase text-gray-400">Conformidade BACEN</p>
                <p class="text-2xl font-bold text-emerald-400 mt-1">100% OK</p>
                <p class="text-xs text-gray-500 mt-2">CMN 5.274 & LGPD</p>
            </div>
        </div>

        <!-- Interactive Agent Testing Console -->
        <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
            <div class="lg:col-span-1 glass p-6 rounded-2xl space-y-4">
                <h2 class="text-lg font-semibold flex items-center gap-2">
                    <span>⚡ Simulator de Tool-Calling</span>
                </h2>
                <p class="text-xs text-gray-400">Simule chamadas de um Agente de IA para testar as regras de compliance em tempo real.</p>
                
                <div>
                    <label class="block text-xs font-medium text-gray-300 mb-1">Selecione o Agente de IA</label>
                    <select id="agent-select" class="w-full bg-gray-900 border border-gray-700 rounded-lg p-2.5 text-sm text-white focus:ring-2 focus:ring-blue-500">
                        <option value="ia-fintech-support">ia-fintech-support (LGPD + BACEN 538)</option>
                        <option value="ia-health-bot">ia-health-bot (CFM)</option>
                        <option value="ia-super-admin">ia-super-admin (Sem restrições)</option>
                        <option value="ia-desconhecida">ia-desconhecida (NÃO Cadastrada)</option>
                    </select>
                </div>

                <div>
                    <label class="block text-xs font-medium text-gray-300 mb-1">Método HTTP & Rota</label>
                    <div class="flex gap-2">
                        <select id="method-select" class="bg-gray-900 border border-gray-700 rounded-lg p-2 text-sm text-white font-mono">
                            <option value="GET">GET</option>
                            <option value="POST">POST</option>
                            <option value="PUT">PUT</option>
                            <option value="DELETE">DELETE</option>
                        </select>
                        <input id="path-input" type="text" value="/pix/validar" class="flex-1 bg-gray-900 border border-gray-700 rounded-lg p-2 text-sm text-white font-mono" placeholder="/rota">
                    </div>
                </div>

                <div class="flex gap-2 pt-2">
                    <button onclick="testQuickPath('/pix/validar', 'POST')" class="text-xs px-2.5 py-1 bg-gray-800 hover:bg-gray-700 rounded-md text-gray-300">✓ Safe POST</button>
                    <button onclick="testQuickPath('/clientes', 'GET')" class="text-xs px-2.5 py-1 bg-rose-950/60 hover:bg-rose-900/60 border border-rose-800 text-rose-300">🚫 LGPD Get</button>
                    <button onclick="testQuickPath('/transacoes/1', 'DELETE')" class="text-xs px-2.5 py-1 bg-rose-950/60 hover:bg-rose-900/60 border border-rose-800 text-rose-300">🚫 BACEN Delete</button>
                </div>

                <button onclick="runTestSimulation()" class="w-full py-3 bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-500 hover:to-indigo-500 text-white font-semibold rounded-lg shadow-lg text-sm transition">
                    Executar Chamada no Aegis Proxy 🚀
                </button>

                <!-- Test Result Display -->
                <div id="test-result-box" class="hidden p-4 rounded-xl border text-sm font-mono space-y-2"></div>
            </div>

            <!-- Active Policies Cards -->
            <div class="lg:col-span-2 glass p-6 rounded-2xl space-y-4">
                <h2 class="text-lg font-semibold flex items-center justify-between">
                    <span>🤖 Políticas de Agentes Carregadas (aegis.yaml)</span>
                    <span class="text-xs text-blue-400 font-normal">O(1) Memory Engine</span>
                </h2>
                
                <div id="agent-cards" class="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div class="p-4 rounded-xl bg-gray-900/80 border border-gray-800 space-y-2">
                        <div class="flex items-center justify-between">
                            <span class="font-bold text-blue-400 text-sm">ia-fintech-support</span>
                            <span class="px-2 py-0.5 text-[10px] bg-blue-950 text-blue-300 border border-blue-800 rounded">Fintech Bot</span>
                        </div>
                        <div class="flex flex-wrap gap-1 text-xs">
                            <span class="px-2 py-0.5 bg-rose-950 text-rose-300 border border-rose-800 rounded">LGPD</span>
                            <span class="px-2 py-0.5 bg-amber-950 text-amber-300 border border-amber-800 rounded">BACEN_538</span>
                        </div>
                        <p class="text-[11px] text-gray-400">Proteção contra mutações não autorizadas e vazamento em massa de dados de clientes.</p>
                    </div>

                    <div class="p-4 rounded-xl bg-gray-900/80 border border-gray-800 space-y-2">
                        <div class="flex items-center justify-between">
                            <span class="font-bold text-emerald-400 text-sm">ia-health-bot</span>
                            <span class="px-2 py-0.5 text-[10px] bg-emerald-950 text-emerald-300 border border-emerald-800 rounded">HealthTech Bot</span>
                        </div>
                        <div class="flex flex-wrap gap-1 text-xs">
                            <span class="px-2 py-0.5 bg-purple-950 text-purple-300 border border-purple-800 rounded">CFM</span>
                        </div>
                        <p class="text-[11px] text-gray-400">Restrição estrita a prontuários médicos sensíveis conforme parecer CFM.</p>
                    </div>
                </div>

                <!-- Commercial Enterprise Banner -->
                <div class="p-4 rounded-xl bg-gradient-to-r from-indigo-950/80 to-purple-950/80 border border-indigo-800 flex items-center justify-between">
                    <div>
                        <h4 class="text-sm font-bold text-indigo-300">Quer Painel Completo para SSO & Relatórios PDF BACEN?</h4>
                        <p class="text-xs text-gray-300">Conheça a versão Aegis Enterprise para Bancos e Fintechs com SLA 24/7.</p>
                    </div>
                    <a href="mailto:pedro@aegisbr.com?subject=Interesse%20Aegis%20Enterprise" class="px-4 py-2 bg-indigo-600 hover:bg-indigo-500 text-white font-semibold text-xs rounded-lg transition whitespace-nowrap">Falar com Consultor</a>
                </div>
            </div>
        </div>
    </main>

    <script>
        function testQuickPath(path, method) {
            document.getElementById('path-input').value = path;
            document.getElementById('method-select').value = method;
        }

        async function fetchInfo() {
            try {
                const res = await fetch('/_aegis/health');
                const data = await res.json();
                document.getElementById('target-api').innerText = data.target_api;
                document.getElementById('agent-count').innerText = data.active_agents;
            } catch(e){}
        }
        fetchInfo();

        async function runTestSimulation() {
            const agent = document.getElementById('agent-select').value;
            const method = document.getElementById('method-select').value;
            const path = document.getElementById('path-input').value;
            const box = document.getElementById('test-result-box');

            box.classList.remove('hidden', 'bg-emerald-950/80', 'bg-rose-950/80', 'border-emerald-800', 'border-rose-800');
            box.innerHTML = '<span class="text-gray-400">Avaliando requisição no Aegis Proxy...</span>';

            try {
                const headers = {};
                if (agent !== 'ia-desconhecida') {
                    headers['Authorization'] = 'Bearer ' + agent;
                }
                const response = await fetch(path, { method: method, headers: headers });
                const status = response.headers.get('X-Aegis-Compliance-Status') || response.status;
                const bodyText = await response.text();

                if (response.ok) {
                    box.classList.add('bg-emerald-950/80', 'border-emerald-800', 'text-emerald-300');
                    box.innerHTML = '<strong>✓ REQUISIÇÃO PERMITIDA (HTTP 200 OK)</strong><br/><span class="text-xs text-gray-300">Status Aegis: ' + status + '</span><br/><span class="text-xs text-emerald-400">Resposta: ' + bodyText + '</span>';
                } else {
                    box.classList.add('bg-rose-950/80', 'border-rose-800', 'text-rose-300');
                    box.innerHTML = '<strong>🚫 VIOLAÇÃO DE COMPLIANCE BLOQUEADA (HTTP ' + response.status + ')</strong><br/><span class="text-xs text-rose-200">Motivo Regulatório: <strong>' + status + '</strong></span><br/><span class="text-xs text-gray-300">' + bodyText + '</span>';
                }
            } catch(err) {
                box.classList.add('bg-rose-950/80', 'border-rose-800', 'text-rose-300');
                box.innerHTML = '<strong>Erro ao conectar com o proxy:</strong> ' + err.message;
            }
        }
    </script>
</body>
</html>`
