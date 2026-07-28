<div align="center">
  <h1>🛡️ Nexo Hub</h1>
  <p><strong>The Open-Core Compliance & Security Shield for Autonomous AI Agents.</strong></p>
  <p><em>Home of NSEP (Nexo Secure Execution Protocol) — Engine Core built for BACEN 538/2025, LGPD, and CFM frameworks.</em></p>
  
  [![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
  [![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go)](https://golang.org/)
  [![CI Pipeline](https://img.shields.io/badge/CI-Passing-brightgreen?logo=github)](https://github.com/pedrodawybida/nexo-hub/actions)
  [![BACEN 538 Engine](https://img.shields.io/badge/BACEN%20538%2F2025-Engine%20Core-blue)](#)
  [![NSEP Protocol](https://img.shields.io/badge/NSEP%20Protocol-94%25%20Token%20Savings-emerald)](#)
  [![Nexo Hub Verified](https://img.shields.io/badge/Nexo%20Hub%20Verified-BACEN%20538%2F2025-blue?logo=shield)](#)
</div>

<div align="center">
  <strong>English | <a href="README.pt-BR.md">Português</a></strong>
</div>

<br />

## 🚀 One-Click Deploy

Deploy your own Nexo Hub instance in 1 click:

[![Deploy on Railway](https://railway.app/button.svg)](https://railway.app/template)
[![Deploy to Render](https://render.com/images/deploy-to-render-button.svg)](https://render.com/deploy)

---

## 🚨 The Problem
Autonomous AI agents (LangChain, CrewAI, AutoGen, LlamaIndex) rely on *Tool-Calling* to execute actions on internal databases and APIs. Providing unsupervised access to an LLM directly violates fundamental cybersecurity and data privacy regulations:
- **CMN Resolution No. 5.274/2025 & BCB 538/2025:** Mandate strict traceability, non-human identity authentication, and immutable audit logging for cyber actions in financial institutions.
- **LGPD (Brazilian General Data Protection Law):** Enforces data minimization and guards against unauthorized bulk data extraction.
- **Context Window Bloat:** Standard MCP tool-calling forces LLMs to ingest 50,000+ tokens of OpenAPI schemas per call, driving latency and API costs out of control.

## 💡 The Solution: Nexo Hub & NSEP Protocol
**Nexo Hub** is an ultra-low latency HTTP reverse proxy & agent gateway (written in Go, overhead < 0.1ms) positioned between your AI Agents and your backend APIs.

It embeds **NSEP (Nexo Secure Execution Protocol)** — a sandboxed JavaScript runtime (`goja`) allowing agents to write typed orchestration scripts instead of performing 10 separate tool calls per operation, reducing LLM token cost by up to 94%.

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

## ⚡ Zero-Friction Setup (`nexo-init`)

Generate `nexo.yaml` and `.env` configuration files in 1 second:

```bash
make init
# or
go run cmd/nexo-init/main.go
```

---

## 🔌 Universal MCP Integration (Claude, ChatGPT, Cursor)

Nexo Hub exposes native **MCP (Model Context Protocol)** compatibility at `POST /_nexo/mcp`:

- **`nsep.search`**: Fast keyword discovery returning endpoint schemas (< 300 tokens footprint).
- **`nsep.execute`**: Runs orchestration code inside the secure Goja VM.

---

## 📁 Pre-configured Compliance Templates (`templates/`)

Nexo Hub includes ready-to-use policy templates in the [`templates/`](file:///Users/pedrodawybida/Developer/projects/aegis-core/templates) directory:

- `templates/bacen_538.yaml`: BACEN 538/2025 & CMN 5.274 Compliance Policy.
- `templates/lgpd_data_minimization.yaml`: LGPD Bulk Data Minimization Policy.
- `templates/iso_42001_ai_governance.yaml`: ISO 42001 AI Management System Policy.
- `templates/soc2_type2_security.yaml`: SOC 2 Type II Non-Human Identity Policy.

---

## 📺 Embedded Web Console & Token Cost Playground

Open `http://localhost:8080/_nexo/dashboard` to test:
- **Live Replay Debugger:** Visual timeline tracing every tool execution and blocked attempt.
- **Token Cost Playground:** Interactive calculator measuring real-time $ context savings.
- **Badge Generator:** Copyable Shields.io compliance badges.

---

## ⚡ Interactive 5-Second Demo

Run the interactive demo script:

```bash
make demo
# or
./demo.sh
```

---

## 🛠️ Makefile Commands

| Command | Description |
| :--- | :--- |
| `make build` | Compiles binaries into `bin/nexo` and `bin/nexo-init` |
| `make init` | Interactive setup wizard generating `nexo.yaml` |
| `make test` | Runs all unit and integration tests |
| `make test-race` | Runs tests with Go's Data Race detector enabled |
| `make run` | Runs Nexo Hub proxy locally |
| `make docker-build` | Builds Docker container image |
| `make demo` | Executes interactive demonstration script |

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
- **NSEP Protocol Enterprise Engine:** Extended JS execution sandbox (`goja`) with advanced isolation.
- **Automated PDF Compliance Evidence Generator:** Single-click compliance evidence PDF generation for Central Bank audits.
- **SSO & RBAC:** Microsoft Entra ID, Okta, and Keycloak integration.
- **SLA & 24/7 Enterprise Support.**

👉 **[Contact us to schedule an Enterprise Edition demo](mailto:pedro@nexohub.lat?subject=Nexo%20Hub%20Enterprise%20Interest)**
