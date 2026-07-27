#!/usr/bin/env bash
# Script para iniciar o Ambiente de Simulação Real e Gravação do Aegis Core

set -e

export PATH="/opt/homebrew/bin:/usr/local/bin:$PATH"

GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
PURPLE='\033[0;35m'
NC='\033[0m'

echo -e "${BLUE}=====================================================================${NC}"
echo -e "${BLUE}🎬 AEGIS CORE - AMBIENTE DE SIMULAÇÃO REAL & GRAVAÇÃO DE VÍDEO      ${NC}"
echo -e "${BLUE}=====================================================================${NC}"

echo -e "\n${YELLOW}1. Compilando o Aegis Core...${NC}"
go build -o bin/aegis cmd/aegis/main.go

echo -e "${YELLOW}2. Iniciando a API Bancária Mock (Porta 9000)...${NC}"
python3 sim/server.py &
MOCK_PID=$!

sleep 1

echo -e "${YELLOW}3. Iniciando o Proxy de Segurança Aegis Core (Porta 8080)...${NC}"
AEGIS_TARGET_API="http://127.0.0.1:9000" ./bin/aegis -port 8080 -log audit_bacen.log &
AEGIS_PID=$!

cleanup() {
    echo -e "\n${PURPLE}Encerrando ambiente de simulação...${NC}"
    kill $MOCK_PID 2>/dev/null || true
    kill $AEGIS_PID 2>/dev/null || true
}
trap cleanup EXIT

sleep 2

echo -e "\n${GREEN}=====================================================================${NC}"
echo -e "${GREEN}✅ AMBIENTE DE SIMULAÇÃO ATIVO E PRONTO PARA GRAVAÇÃO!               ${NC}"
echo -e "${GREEN}=====================================================================${NC}"
echo -e "📺 Console Web Visual: ${YELLOW}http://localhost:8080/_aegis/dashboard${NC}"
echo -e "📝 Logs de Auditoria:   ${YELLOW}audit_bacen.log${NC}"
echo -e "🎯 API Bancária Mock:   ${YELLOW}http://localhost:9000${NC}"
echo -e "${GREEN}=====================================================================${NC}"

echo -e "\n${PURPLE}Pressione ENTER para testar uma simulação de chamada de Agente...${NC}"
read -r

echo -e "\n${BLUE}[Simulação Agente IA] Enviando requisição bloqueada por LGPD (GET /clientes):${NC}"
curl -i -X GET http://localhost:8080/clientes -H "Authorization: Bearer ia-fintech-support" || true

echo -e "\n${PURPLE}Ambiente ativo! Pressione CTRL+C quando terminar a gravação do seu vídeo.${NC}"
wait
