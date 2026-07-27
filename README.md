<div align="center">
  <h1>🛡️ Nexo Hub</h1>
  <p><strong>The Open-Core Compliance & Security Shield for Autonomous AI Agents.</strong></p>
  <p><em>Home of NSEP (Nexo Secure Execution Protocol) — Engine Core built for BACEN 538/2025, LGPD, and CFM frameworks.</em></p>
  
  [![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
  [![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go)](https://golang.org/)
  [![CI Pipeline](https://img.shields.io/badge/CI-Passing-brightgreen?logo=github)](https://github.com/pedrodawybida/nexo-hub/actions)
  [![BACEN 538 Engine](https://img.shields.io/badge/BACEN%20538%2F2025-Engine%20Core-blue)](#)
  [![NSEP Protocol](https://img.shields.io/badge/NSEP%20Protocol-90%25%20Token%20Savings-emerald)](#)
</div>

<div align="center">
  <strong>English | <a href="README.pt-BR.md">Português</a></strong>
</div>

<br />

## 🚨 The Problem
Autonomous AI agents (LangChain, CrewAI, AutoGen, LlamaIndex) rely on *Tool-Calling* to execute actions on internal databases and APIs. Providing unsupervised access to an LLM directly violates fundamental cybersecurity and data privacy regulations:
- **CMN Resolution No. 5.274/2025 & BCB 538/2025:** Mandate strict traceability, non-human identity authentication, and immutable audit logging for cyber actions in financial institutions.
- **LGPD (Brazilian General Data Protection Law):** Enforces data minimization and guards against unauthorized bulk data extraction.
- **Context Window Bloat:** Standard MCP tool-calling forces LLMs to ingest 50,000+ tokens of OpenAPI schemas per call, driving latency and API costs out of control.

## 💡 The Solution: Nexo Hub & NSEP Protocol
**Nexo Hub** is an ultra-low latency HTTP reverse proxy & agent gateway (written in Go, overhead < 0.1ms) positioned between your AI Agents and your backend APIs.

It embeds **NSEP (Nexo Secure Execution Protocol)** — a sandboxed JavaScript runtime (`goja`) allowing agents to write typed orchestration scripts instead of performing 10 separate tool calls per operation, reducing LLM token cost by up to 90%.

```mermaid
flowchart LR
    A[🤖 AI Agent / LLM Tool-Call] -->|Bearer Token + MCP Req| B[🛡️ Nexo Hub Proxy]
    B -->|Check Policy O1| C{Compliance Valid?}
    C -->|YES - ALLOWED| D[🟢 Internal Backend API]
    C -->|NO - BLOCKED| E[🔴 403 Forbidden Response]
    B -->|Immutable JSON Log| F[📜 audit_bacen.log]
```

---

## ⚡ NSEP Protocol Benchmark (Token Savings)

| Protocol | Strategy | Prompt Footprint | Latency | Compliance Traceability |
| :--- | :--- | :--- | :--- | :--- |
| **Traditional MCP** | 1 Tool-Call per Operation | ~45,000 tokens | 12.4s | Uncorrelated Logs |
| **NSEP (Nexo Hub)** | Typed JS Script + Sandbox | **~2,800 tokens** | **0.8s** | **Correlated `execution_id`** |

---

## 🔌 Universal MCP Integration (Claude, ChatGPT, Cursor)

Nexo Hub exposes native **MCP (Model Context Protocol)** compatibility at `POST /_nexo/mcp`:

- **`nsep.search`**: Fast keyword discovery returning endpoint schemas (< 300 tokens footprint).
- **`nsep.execute`**: Runs orchestration code inside the secure Goja VM.

### Example MCP Request (`nsep.execute`)
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "nsep.execute",
    "arguments": {
      "code": "var r = nsep.call('check_saldo', { account_id: '101' }); r.status;"
    }
  }
}
```

---

## ⚡ Interactive 5-Second Demo

Run the interactive demo script to see Nexo Hub block non-compliant agent requests in real time:

```bash
make demo
# or
./demo.sh
```

---

## 🚀 Quickstart

### 1. Define Policy Configuration (`nexo.yaml`)
Define permissions for each non-human AI identity:
```yaml
target_api: "http://localhost:9000" 

agents:
  # Fintech Support Agent
  - id: "ia-fintech-support"
    modes:
      - "LGPD"      # Blocks bulk GET operations on customer endpoints
      - "BACEN_538" # Blocks unapproved state mutations (DELETE/PUT)
  
  # HealthTech Bot
  - id: "ia-health-bot"
    modes:
      - "CFM"       # Restricts access to medical records without oversight
```

### 2. Run via Docker
```bash
docker build -t nexo-hub .
docker run -p 8080:8080 -d nexo-hub
```

### 3. Run Locally via Makefile
```bash
make build
make run
```

---

## 🛠️ Makefile Commands

| Command | Description |
| :--- | :--- |
| `make build` | Compiles the executable binary into `bin/nexo` |
| `make test` | Runs the full unit and integration test suite |
| `make test-race` | Runs tests with Go's Data Race detector enabled |
| `make run` | Runs Nexo Hub proxy locally |
| `make docker-build` | Builds the container Docker image |
| `make demo` | Executes the interactive demonstration script |

---

## ⚙️ Environment Variables & CLI Flags

Configure Nexo Hub via command-line flags or environment variables:

| Environment Variable | CLI Flag | Default | Description |
| :--- | :--- | :--- | :--- |
| `NEXO_PORT` | `-port` | `8080` | Port for the proxy server listener |
| `NEXO_CONFIG` | `-config` | `nexo.yaml` | Path to YAML configuration file |
| `NEXO_LOG_FILE` | `-log` | `audit_bacen.log` | Path to append-only immutable audit log file |
| `NEXO_TARGET_API` | - | `http://localhost:9000` | Target internal backend API URL |
| `NEXO_DRY_RUN` | `-dry-run` | `false` | Enable Dry-Run (Shadow / Audit-Only) Mode without blocking requests |

---

## 🏥 Health Check & System Endpoints

Nexo Hub exposes native system endpoints requiring no agent identity token:

- `GET /_nexo/health` (or `/healthz`): Returns service operational status and active agent count.
- `GET /_nexo/dashboard`: Opens the embedded visual Web Console for live testing.
- `POST /_nexo/mcp`: MCP tool-calling transport endpoint.
- HTTP responses include the `X-Nexo-Compliance-Status` header specifying `ALLOWED` or the exact violation reason.

---

## 📜 Audit Logs Ready for Auditors
Every request attempt appends an immutable, thread-safe JSON entry to `audit_bacen.log`:

```json
{
  "timestamp": "2026-07-27T14:22:35Z",
  "agent_id": "ia-fintech-support",
  "execution_id": "exec_1785185580_7ef2",
  "sequence_in_execution": 2,
  "action": "DELETE /transacoes/99",
  "tool_payload": "{\"id\": 99}",
  "result": "BLOCKED_BACEN_538_MUTATION_DENIED",
  "ip_address": "127.0.0.1:54569"
}
```

---

## 🏢 Enterprise Edition
This repository contains the Open-Source Core engine under the MIT License.

For enterprise production deployments in Banks and Fintechs, we offer:
- **NSEP Protocol Enterprise Engine:** Extended JS execution sandbox (`goja`) with advanced isolation & rate limiting.
- **Automated PDF Compliance Evidence Generator:** Single-click compliance evidence PDF generation for Central Bank audits.
- **SSO & RBAC:** Microsoft Entra ID, Okta, and Keycloak integration.
- **SLA & 24/7 Enterprise Support.**

👉 **[Contact us to schedule an Enterprise Edition demo](mailto:pedro@nexohub.lat?subject=Nexo%20Hub%20Enterprise%20Interest)**
