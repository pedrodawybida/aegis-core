<div align="center">
  <h1>🛡️ Nexo Hub</h1>
  <p><strong>The Open-Core Compliance & Security Shield for Autonomous AI Agents.</strong></p>
  <p><em>Casa do protocolo NSEP (Nexo Secure Execution Protocol) — Motor de Conformidade para BCB 538/2025, LGPD e CFM.</em></p>
  
  [![English](https://img.shields.io/badge/Language-English-blue)](#)
  [![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
  [![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go)](https://golang.org/)
  [![CI Pipeline](https://img.shields.io/badge/CI-Passing-brightgreen?logo=github)](https://github.com/pedrodawybida/nexo-hub/actions)
</div>

<div align="center">
  <strong><a href="README.md">English</a> | Português</strong>
</div>

<br />

## 🚨 O Problema
Agentes de IA autônomos (LangChain, CrewAI, AutoGen, LlamaIndex) utilizam *Tool-Calling* para executar ações em bancos de dados e APIs internas. Dar acesso não-supervisionado a uma IA viola diretamente os controles mínimos de cibersegurança exigidos por reguladores brasileiros:
- **Resoluções CMN nº 5.274/2025 e BCB 538/2025:** Exigem rastreabilidade absoluta, autenticação de identidades não-humanas e log imutável de ações cibernéticas em instituições financeiras.
- **LGPD (Lei Geral de Proteção de Dados):** Exige o princípio da minimização de dados e proteção contra extração não autorizada em massa.

## 💡 A Solução: Nexo Hub & NSEP
**Nexo Hub** atua como um *Reverse Proxy & Agent Gateway* de ultrabaixa latência (escrito em Go, latência < 0.1ms) posicionado entre os seus Agentes de IA e as suas APIs de backend.

O Nexo Hub verifica se a ação solicitada viola os *Templates de Compliance BR* e **bloqueia proativamente ações destrutivas ou vazamentos**, gerando logs estruturados e imutáveis prontos para a auditoria do Banco Central.

```mermaid
flowchart LR
    A[🤖 AI Agent / LLM Tool-Call] -->|Bearer Token + HTTP Req| B[🛡️ Nexo Hub Proxy]
    B -->|Check Policy O1| C{Compliance Valid?}
    C -->|YES - ALLOWED| D[🟢 Internal Backend API]
    C -->|NO - BLOCKED| E[🔴 403 Forbidden Response]
    B -->|Immutable JSON Log| F[📜 audit_bacen.log]
```

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

## ⚙️ Variáveis de Ambiente & Flags CLI

Você pode configurar o Nexo Hub via argumentos de linha de comando ou variáveis de ambiente:

| Variável de Ambiente | Flag CLI | Padrão | Descrição |
| :--- | :--- | :--- | :--- |
| `NEXO_PORT` | `-port` | `8080` | Porta onde o servidor proxy escutará |
| `NEXO_CONFIG` | `-config` | `nexo.yaml` | Caminho do arquivo de configuração YAML |
| `NEXO_LOG_FILE` | `-log` | `audit_bacen.log` | Caminho do arquivo de log imutável |
| `NEXO_TARGET_API` | - | `http://localhost:9000` | URL da API interna de destino protegida |
| `NEXO_DRY_RUN` | `-dry-run` | `false` | Ativa o Modo Sombra (Dry-Run / Audit-Only) sem bloquear requisições |

---

## 🏥 Endpoints de Saúde e Diagnósticos

O Nexo Hub expõe endpoints de sistema nativos sem necessidade de token de agente:

- `GET /_nexo/health` (ou `/healthz`): Retorna status operacional do serviço e quantidade de políticas ativas.
- `GET /_nexo/dashboard`: Abre o Console Web Visual em tempo real para testes e auditoria.
- Respostas HTTP contêm o cabeçalho `X-Nexo-Compliance-Status` indicando `ALLOWED` ou o motivo exato do bloqueio.

---

## 📺 Console Web Visual Embutido

Acessem `http://localhost:8080/_nexo/dashboard` diretamente no navegador para testar requisições dos agentes, visualizar os cartões de políticas ativas e validar simulações em tempo real sem precisar instalar nada além do Nexo Hub.

---

## 🧪 Testes Automatizados

Para executar os testes automatizados com detector de race condition:
```bash
go test -v -race ./...
```

---

## 📜 Logs de Auditoria Prontos para o Regulador
Qualquer tentativa de acesso gera um log JSON imutável e thread-safe no arquivo `audit_bacen.log`.

```json
{
  "timestamp": "2026-07-27T14:22:35Z",
  "agent_id": "ia-fintech-support",
  "action": "DELETE /transacoes/99",
  "tool_payload": "",
  "result": "BLOCKED_BACEN_538_MUTATION_DENIED",
  "ip_address": "127.0.0.1:54569"
}
```

---

## 🏛️ Arquitetura & Protocolo NSEP

O **Nexo Hub** foi projetado como um motor de políticas (*Policy Engine*) leve, conciso e extensível para segurança de identidades não-humanas.

- **Identidade & Autenticação:** No motor Core, o token do agente é validado em memória via arquivo `nexo.yaml`. Em ambientes de produção corporativos, o Nexo Hub opera integrado a API Gateways / Service Mesh (Kong, Envoy, Ambassador) delegando autenticação para mTLS, JWT/OIDC ou OAuth2.
- **Motor de Regras:** O pacote `compliance` implementa regras base para BACEN 538/2025, LGPD e CFM.

---

## 🏢 Licença Enterprise
Este repositório contém o motor (*Core*) Open-Source sob licença MIT.

Para a versão Corporativa (*Enterprise Edition*), oferecemos:
- **Protocolo NSEP (Nexo Secure Execution Protocol):** Motor de execução em sandbox JS (`goja`) com redução de tokens de até 90% e auditoria correlacionada por `execution_id`.
- **Painel Visual Web:** Gerenciamento centralizado de políticas por CISO e auditores.
- **Relatórios Automatizados em PDF:** Geração de relatórios de auditoria com 1 clique para entrega ao Banco Central.
- **SSO & RBAC:** Integração com Microsoft Entra ID, Okta e Keycloak.
- **SLA e Suporte 24/7.**

👉 **[Fale conosco e agende uma demonstração da versão Enterprise](mailto:pedro@nexohub.lat?subject=Interesse%20Nexo%20Hub%20Enterprise)**
