<div align="center">
  <h1>🛡️ Aegis</h1>
  <p><strong>The Compliance Shield for Autonomous AI. (Built for BACEN 538/2025 and LGPD).</strong></p>
  
  [![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
  [![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://golang.org/)
  [![BACEN Compliant](https://img.shields.io/badge/BACEN%20538%2F2025-Compliant-success)](#)
  [![LGPD Compliant](https://img.shields.io/badge/LGPD-Compliant-success)](#)
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
  # Agente Financeiro
  - id: "ia-fintech-support"
    modes:
      - "LGPD"      # Bloqueia GET em massa em rotas de dados sensíveis
      - "BACEN_538" # Bloqueia mutações severas (DELETE/PUT) sem aprovação
  
  # Agente de Saúde (HealthTech)
  - id: "ia-health-bot"
    modes:
      - "CFM"       # Proíbe acesso a prontuários sem validação humana
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

## 👨‍💻 Por baixo dos panos (Como o Proxy funciona)
O coração da Aegis é absurdamente simples e customizável. Ele intercepta e avalia as regras em microssegundos:

```go
// internal/proxy/proxy.go
func (p *AegisProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	agentToken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	
    // 1. Identifica e valida o agente no YAML
	policy, exists := p.policyDB[agentToken]
	if !exists {
		http.Error(w, `{"error": "Aegis: Agent identity not recognized"}`, http.StatusForbidden)
		return
	}

	// 2. Aplica as regras Brasileiras de LGPD / BACEN
	allowed, reason := compliance.EvaluateRequest(policy, r.Method, r.URL.Path)
	
	// 3. Salva o Log Imutável de Auditoria (Exigência BACEN)
	p.logger.LogAction(agentToken, r.Method, payload, reason, r.RemoteAddr)

    // 4. Bloqueia ou Encaminha
	if !allowed {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(fmt.Sprintf(`{"error": "Aegis Compliance Violation: %s"}`, reason)))
		return
	}
	p.reverse.ServeHTTP(w, r)
}
```

---

## 📜 Logs de Auditoria Prontos para o Regulador
Qualquer tentativa de acesso gera um log JSON imutável no arquivo `audit_bacen.log`.

**Exemplo de log de bloqueio:**
```json
{
  "timestamp": "2026-05-05T02:35:48Z",
  "agent_id": "ia-fintech-support",
  "action": "GET /clientes",
  "result": "BLOCKED_LGPD_BULK_DATA_ACCESS_DENIED",
  "ip_address": "[::1]:52543"
}
```

---

## 🏢 Licença Enterprise
Este é o motor (Core) Open-Source do Aegis sob licença MIT. 

Para a versão Corporativa (Obrigatória para Produção em Bancos), nós oferecemos:
- **Painel Visual Web:** Gerenciamento dos agentes por pessoas não-técnicas (CISO, Auditores).
- **Relatórios Automatizados em PDF:** Geração de provas de conformidade com 1 clique para entregar para o BACEN.
- **SSO & RBAC:** Integração com Microsoft Entra ID / Okta.
- **SLA e Suporte 24/7.**

👉 **[Fale conosco e agende uma demonstração da versão Enterprise](mailto:pedro@aegisbr.com?subject=Interesse%20Aegis%20Enterprise)**
