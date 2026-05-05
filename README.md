<div align="center">
  <h1>🛡️ Aegis</h1>
  <p><strong>The Compliance Shield for Autonomous AI. (Built for BACEN 538/2025 and LGPD).</strong></p>
</div>

<br />

## 🚨 O Problema
Agentes de IA autônomos (LangChain, CrewAI, etc.) utilizam *Tool-Calling* para executar ações em bancos de dados e APIs internas. O problema? Dar acesso não-supervisionado a uma IA viola diretamente os controles mínimos de cibersegurança exigidos por reguladores brasileiros:
- **Resoluções CMN nº 5.274/2025 e BCB 538/2025:** Exigem rastreabilidade absoluta e log imutável de ações cibernéticas em instituições financeiras.
- **LGPD (Lei Geral de Proteção de Dados):** Exige o princípio da minimização de dados e proteção contra extração de dados em massa.

## 💡 A Solução
**Aegis** atua como um *Reverse Proxy* ultrarrápido (escrito em Go, latência < 5ms) que se posiciona entre o seu Agente de IA e as suas APIs internas.

Em vez do Agente acessar o seu banco de dados diretamente, ele envia a requisição para o Aegis. O Aegis verifica se a ação viola os *Templates de Compliance BR* e **bloqueia proativamente ações destrutivas ou vazamentos**, gerando logs estruturados prontos para a auditoria do Banco Central.

---

## 🚀 Como Iniciar (Quickstart)

### 1. Crie a Configuração
Defina as permissões no arquivo `aegis.yaml`:
```yaml
target_api: "http://localhost:9000" 

agents:
  - id: "ia-fintech-support"
    modes:
      - "LGPD"      # Bloqueia GET em massa em rotas de dados sensíveis
      - "BACEN_538" # Bloqueia mutações severas (DELETE/PUT) sem aprovação
```

### 2. Rodando via Docker (Recomendado On-Premise)
Bancos e fintechs exigem isolamento de rede. O Aegis é 100% conteinerizado:
```bash
docker build -t aegis-core .
docker run -p 8080:8080 -d aegis-core
```

### 3. Rodando localmente para testes (Dev)
```bash
go mod tidy
go run cmd/aegis/main.go
```

---

## 📜 Logs de Auditoria Prontos para o Regulador
Qualquer tentativa de acesso gera um log JSON imutável (exigência do BACEN) no arquivo `audit_bacen.log`.

**Exemplo de log de bloqueio:**
```json
{
  "timestamp": "2026-05-05T02:35:48Z",
  "agent_id": "ia-fintech-support",
  "action": "GET /clientes",
  "tool_payload": "",
  "result": "BLOCKED_LGPD_BULK_DATA_ACCESS_DENIED",
  "ip_address": "[::1]:52543"
}
```

---

## 🏢 Licença Enterprise
Este é o motor (Core) Open-Source do Aegis.
Para adquirir a versão Enterprise (que inclui Deploy Automatizado On-Premise, Painel Visual de Gerenciamento, exportação de Relatórios PDF de Auditoria e Integração SSO para corporações), entre em contato conosco.
