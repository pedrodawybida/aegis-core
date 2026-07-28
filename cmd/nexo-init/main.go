package main

import (
	"fmt"
	"os"
	"path/filepath"
)

const defaultYamlConfig = `# ==========================================
# Nexo Hub & NSEP Protocol Configuration
# ==========================================
target_api: "http://localhost:9000"

agents:
  # Fintech Support Agent (LGPD + BACEN 538)
  - id: "ia-fintech-support"
    modes:
      - "LGPD"      # Restricts bulk GET customer data extraction
      - "BACEN_538" # Restricts unauthorized DELETE/PUT mutations

  # HealthTech Bot (CFM Guidelines)
  - id: "ia-health-bot"
    modes:
      - "CFM"       # Restricts medical record access without human oversight
`

const defaultEnvConfig = `# Nexo Hub Environment Configuration
NEXO_PORT=8080
NEXO_CONFIG=nexo.yaml
NEXO_LOG_FILE=audit_bacen.log
NEXO_TARGET_API=http://localhost:9000
NEXO_DRY_RUN=false
`

func main() {
	fmt.Println("🛡️  Nexo Hub CLI Init Generator")
	fmt.Println("==========================================")

	yamlPath := "nexo.yaml"
	if _, err := os.Stat(yamlPath); err == nil {
		fmt.Printf("⚠️  '%s' já existe no diretório atual. Ignorando sobrescrita.\n", yamlPath)
	} else {
		if err := os.WriteFile(yamlPath, []byte(defaultYamlConfig), 0644); err != nil {
			fmt.Printf("❌ Erro ao criar '%s': %v\n", yamlPath, err)
			os.Exit(1)
		}
		fmt.Printf("✅ Criado '%s' com políticas padrão (LGPD, BACEN 538, CFM).\n", yamlPath)
	}

	envPath := ".env"
	if _, err := os.Stat(envPath); err == nil {
		fmt.Printf("⚠️  '%s' já existe no diretório atual. Ignorando sobrescrita.\n", envPath)
	} else {
		if err := os.WriteFile(envPath, []byte(defaultEnvConfig), 0644); err != nil {
			fmt.Printf("❌ Erro ao criar '%s': %v\n", envPath, err)
			os.Exit(1)
		}
		fmt.Printf("✅ Criado '%s' com portas e variáveis de ambiente.\n", envPath)
	}

	absYaml, _ := filepath.Abs(yamlPath)
	fmt.Println("\n🚀 Inicialização do Nexo Hub concluída com sucesso!")
	fmt.Printf("📍 Arquivo de configuração: %s\n", absYaml)
	fmt.Println("👉 Para executar o proxy: go run cmd/nexo/main.go")
}
