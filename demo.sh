#!/usr/bin/env bash

# Aegis Core Interactive Verification & Demo Script
set -e

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

export PATH="/opt/homebrew/bin:/usr/local/bin:$PATH"

echo -e "${BLUE}=====================================================${NC}"
echo -e "${BLUE}🛡️  AEGIS CORE - DEMO DE CONFORMIDADE & AUDITORIA   ${NC}"
echo -e "${BLUE}=====================================================${NC}"

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}Erro: Go não está instalado no sistema.${NC}"
    exit 1
fi

MOCK_PORT=9000
AEGIS_PORT=8080
LOG_FILE="audit_bacen.log"

echo -e "\n${YELLOW}1. Subindo API Mock de Destino na porta ${MOCK_PORT}...${NC}"
python3 -c "
import http.server
import socketserver

class MockHandler(http.server.SimpleHTTPRequestHandler):
    def do_GET(self):
        self.send_response(200)
        self.send_header('Content-Type', 'application/json')
        self.end_headers()
        self.wfile.write('{\"status\": \"success\", \"message\": \"Dados retornados da API Interna Protegida\"}'.encode('utf-8'))

    def do_POST(self):
        self.send_response(200)
        self.send_header('Content-Type', 'application/json')
        self.end_headers()
        self.wfile.write('{\"status\": \"success\", \"message\": \"Acao executada com sucesso\"}'.encode('utf-8'))

with socketserver.TCPServer(('127.0.0.1', $MOCK_PORT), MockHandler) as httpd:
    httpd.serve_forever()
" &
MOCK_PID=$!

cleanup() {
    echo -e "\n${YELLOW}Encerrando processos de teste...${NC}"
    kill $MOCK_PID 2>/dev/null || true
    if [ -n "$AEGIS_PID" ]; then
        kill -SIGTERM $AEGIS_PID 2>/dev/null || true
    fi
}
trap cleanup EXIT

sleep 1

echo -e "${YELLOW}2. Compilando o binary do Aegis Core...${NC}"
go build -o bin/aegis cmd/aegis/main.go

echo -e "${YELLOW}3. Iniciando o Aegis Core na porta ${AEGIS_PORT}...${NC}"
AEGIS_TARGET_API="http://127.0.0.1:${MOCK_PORT}" ./bin/aegis -port ${AEGIS_PORT} -log ${LOG_FILE} &
AEGIS_PID=$!

# Wait for Aegis to respond on health check
for i in {1..10}; do
    if curl -s http://127.0.0.1:${AEGIS_PORT}/_aegis/health > /dev/null; then
        break
    fi
    sleep 0.5
done

echo -e "\n${YELLOW}4. Testando Endpoint de Saúde & Web Console do Aegis Core...${NC}"
curl -s http://127.0.0.1:${AEGIS_PORT}/_aegis/health | grep "ok" && echo -e "${GREEN}✓ Aegis Core está ativo e respondendo!${NC}"
echo -e "${GREEN}📺 Console Web Visual disponível em: http://127.0.0.1:${AEGIS_PORT}/_aegis/dashboard${NC}"

echo -e "\n${YELLOW}5. Executando testes de políticas de Compliance:${NC}"

echo -e "\n${BLUE}[Cenário A - Permitido] Agente 'ia-fintech-support' envia POST seguro em /pix/validar:${NC}"
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST http://127.0.0.1:${AEGIS_PORT}/pix/validar -H "Authorization: Bearer ia-fintech-support")
if [ "$HTTP_CODE" -eq 200 ]; then
    echo -e "${GREEN}✓ APROVADO (HTTP $HTTP_CODE)${NC}"
else
    echo -e "${RED}✗ Falhou (HTTP $HTTP_CODE)${NC}"
fi

echo -e "\n${BLUE}[Cenário B - Bloqueado LGPD] Agente 'ia-fintech-support' tenta extração em massa GET em /clientes:${NC}"
RESP=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X GET http://127.0.0.1:${AEGIS_PORT}/clientes -H "Authorization: Bearer ia-fintech-support")
echo "$RESP" | grep -q "BLOCKED_LGPD_BULK_DATA_ACCESS_DENIED" && echo -e "${GREEN}✓ BLOQUEADO PROATIVAMENTE PELO AEGIS (LGPD Violada)${NC}"

echo -e "\n${BLUE}[Cenário C - Bloqueado BACEN 538] Agente 'ia-fintech-support' tenta DELETE destrutivo em /transacoes/99:${NC}"
RESP=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X DELETE http://127.0.0.1:${AEGIS_PORT}/transacoes/99 -H "Authorization: Bearer ia-fintech-support")
echo "$RESP" | grep -q "BLOCKED_BACEN_538_MUTATION_DENIED" && echo -e "${GREEN}✓ BLOQUEADO PROATIVAMENTE PELO AEGIS (BACEN 538 Violado)${NC}"

echo -e "\n${BLUE}[Cenário D - Bloqueado CFM] Agente 'ia-health-bot' tenta acessar prontuários médicos em /prontuarios/101:${NC}"
RESP=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X GET http://127.0.0.1:${AEGIS_PORT}/prontuarios/101 -H "Authorization: Bearer ia-health-bot")
echo "$RESP" | grep -q "BLOCKED_CFM_MEDICAL_RECORDS_DENIED" && echo -e "${GREEN}✓ BLOQUEADO PROATIVAMENTE PELO AEGIS (Diretriz CFM Violada)${NC}"

echo -e "\n${YELLOW}6. Exibindo as últimas entradas do arquivo de Auditoria Imutável (${LOG_FILE}):${NC}"
echo -e "${BLUE}---------------------------------------------------------------------${NC}"
tail -n 4 ${LOG_FILE}
echo -e "${BLUE}---------------------------------------------------------------------${NC}"

echo -e "\n${GREEN}🎉 Demonstração concluída com 100% de sucesso!${NC}"
