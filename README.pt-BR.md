<div align="center">
  <h1>🛡️ Nexo Hub</h1>
  <p><strong>The Open-Core Compliance & Security Shield for Autonomous AI Agents.</strong></p>
  <p><em>Casa do protocolo NSEP (Nexo Secure Execution Protocol) — Motor de Conformidade para BCB 538/2025, LGPD e CFM.</em></p>
  
  [![English](https://img.shields.io/badge/Language-English-blue)](#)
  [![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
  [![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go)](https://golang.org/)
  [![CI Pipeline](https://img.shields.io/badge/CI-Passing-brightgreen?logo=github)](https://github.com/pedrodawybida/nexo-hub/actions)
  [![Protocolo NSEP](https://img.shields.io/badge/Protocolo%20NSEP-Economia%2090%25%20Tokens-emerald)](#)
</div>

<div align="center">
  <strong><a href="README.md">English</a> | Português</strong>
</div>

<br />

## 🚨 O Problema
Agentes de IA autônomos (LangChain, CrewAI, AutoGen, LlamaIndex) utilizam *Tool-Calling* para executar ações em bancos de dados e APIs internas. Dar acesso não-supervisionado a uma IA viola diretamente os controles mínimos de cibersegurança exigidos por reguladores brasileiros:
- **Resoluções CMN nº 5.274/2025 e BCB 538/2025:** Exigem rastreabilidade absoluta, autenticação de identidades não-humanas e log imutável de ações cibernéticas em instituições financeiras.
- **LGPD (Lei Geral de Proteção de Dados):** Exige o princípio da minimização de dados e proteção contra extração não autorizada em massa.
- **Explosão do Custo de Tokens:** O MCP tradicional obriga o LLM a ler schemas de 50.000+ tokens por chamada de ferramenta, tornando a latência e o custo inviáveis.

## 💡 A Solução: Nexo Hub & Protocolo NSEP
**Nexo Hub** atua como um *Reverse Proxy & Agent Gateway* de ultrabaixa latência (escrito em Go, latência < 0.1ms) posicionado entre os seus Agentes de IA e as suas APIs de backend.

Ele embute o **NSEP (Nexo Secure Execution Protocol)** — um ambiente de execução em JavaScript isolado (`goja`) que permite ao agente escrever scripts orquestrados tipados em vez de realizar 10 chamadas de ferramentas individuais, reduzindo o custo de tokens em até 90%.

```mermaid
flowchart LR
    A[🤖 AI Agent / LLM Tool-Call] -->|Bearer Token + MCP Req| B[🛡️ Nexo Hub Proxy]
    B -->|Check Policy O1| C{Compliance Valid?}
    C -->|YES - ALLOWED| D[🟢 Internal Backend API]
    C -->|NO - BLOCKED| E[🔴 403 Forbidden Response]
    B -->|Immutable JSON Log| F[📜 audit_bacen.log]
```

---

## ⚡ Benchmark do Protocolo NSEP (Economia de Tokens)

| Protocolo | Estratégia | Pegada no Prompt | Latência | Rastreabilidade de Auditoria |
| :--- | :--- | :--- | :--- | :--- |
| **MCP Tradicional** | 1 Tool-Call por Operação | ~45.000 tokens | 12.4s | Logs Descorrelacionados |
| **NSEP (Nexo Hub)** | Script JS Tipado + Sandbox | **~2.800 tokens** | **0.8s** | **Correlacionado por `execution_id`** |

---

## 🔌 Integração Universal MCP (Claude, ChatGPT, Cursor)

O Nexo Hub expõe compatibilidade nativa com o padrão **MCP (Model Context Protocol)** no endpoint `POST /_nexo/mcp`:

- **`nsep.search`**: Busca ultra-rápida por intenção de rotas com retorno enxuto (< 300 tokens).
- **`nsep.execute`**: Executa o código orquestrado dentro da VM segura Goja.

---

## ⚡ Demonstração Interativa em 5 Segundos

Execute o script interativo para testar o Nexo Hub bloqueando requisições não conformes em tempo real:

```bash
make demo
# ou
./demo.sh
```

---

## 🚀 Como Iniciar (Quickstart)

### 1. Crie a Configuração (`nexo.yaml`)
Defina as permissões de cada agente:
```yaml
target_api: "http://localhost:9000" 

agents:
  # Agente de Suporte Fintech
  - id: "ia-fintech-support"
    modes:
      - "LGPD"      # Bloqueia GET em massa em rotas de clientes
      - "BACEN_538" # Bloqueia mutações destrutivas (DELETE/PUT) sem homologação
  
  # Agente de Saúde (HealthTech)
  - id: "ia-health-bot"
    modes:
      - "CFM"       # Proíbe acesso a prontuários médicos sem validação humana
```

### 2. Rodando via Docker (Conteinerizado)
```bash
docker build -t nexo-hub .
docker run -p 8080:8080 -d nexo-hub
```

### 3. Rodando Localmente com Makefile
```bash
make build
make run
```

---

## 🛠️ Comandos do Makefile

| Comando | Descrição |
| :--- | :--- |
| `make build` | Compila o executável `bin/nexo` |
| `make test` | Executa a suíte de testes unitários e de integração |
| `make test-race` | Executa testes com detector de *Data Race* ativo |
| `make run` | Executa o proxy localmente |
| `make docker-build` | Constrói a imagem Docker do Nexo Hub |
| `make demo` | Executa a demonstração prática interativa |

---

## 📜 Logs de Auditoria Prontos para o Regulador
Qualquer tentativa de acesso gera um log JSON imutável e thread-safe no arquivo `audit_bacen.log`.

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

## 🏢 Licença Enterprise
Este repositório contém o motor (*Core*) Open-Source sob licença MIT.

Para a versão Corporativa (*Enterprise Edition*), oferecemos:
- **Protocolo NSEP Enterprise Engine:** Motor de execução estendido em sandbox JS (`goja`) com isolamento avançado.
- **Relatórios Automatizados em PDF:** Geração de relatórios de auditoria com 1 clique para entrega ao Banco Central.
- **SSO & RBAC:** Integração com Microsoft Entra ID, Okta e Keycloak.
- **SLA e Suporte 24/7.**

👉 **[Fale conosco e agende uma demonstração da versão Enterprise](mailto:pedro@nexohub.lat?subject=Interesse%20Nexo%20Hub%20Enterprise)**
