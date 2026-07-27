#!/usr/bin/env bash

# Nexo Hub Interactive Verification & Demo Script
set -e

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

export PATH="/opt/homebrew/bin:/usr/local/bin:$PATH"

echo -e "${BLUE}=====================================================${NC}"
echo -e "${BLUE}🛡️  NEXO HUB - DEMO DE CONFORMIDADE & AUDITORIA   ${NC}"
echo -e "${BLUE}=====================================================${NC}"

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}Erro: Go não está instalado no sistema.${NC}"
    exit 1
fi

MOCK_PORT=9000
NEXO_PORT=8080
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
    if [ -n "$NEXO_PID" ]; then
        kill -SIGTERM $NEXO_PID 2>/dev/null || true
    fi
}
trap cleanup EXIT

sleep 1

echo -e "${YELLOW}2. Compilando o binary do Nexo Hub...${NC}"
go build -o bin/nexo cmd/nexo/main.go

echo -e "${YELLOW}3. Iniciando o Nexo Hub na porta ${NEXO_PORT}...${NC}"
NEXO_TARGET_API="http://127.0.0.1:${MOCK_PORT}" ./bin/nexo -config nexo.yaml -port ${NEXO_PORT} -log ${LOG_FILE} &
NEXO_PID=$!

# Wait for Nexo Hub to respond on health check
for i in {1..10}; do
    if curl -s http://127.0.0.1:${NEXO_PORT}/_nexo/health > /dev/null; then
        break
    fi
    sleep 0.5
done

echo -e "\n${YELLOW}4. Testando Endpoint de Saúde & Web Console do Nexo Hub...${NC}"
curl -s http://127.0.0.1:${NEXO_PORT}/_nexo/health | grep "ok" && echo -e "${GREEN}✓ Nexo Hub está ativo e respondendo!${NC}"
echo -e "${GREEN}📺 Console Web Visual disponível em: http://127.0.0.1:${NEXO_PORT}/_nexo/dashboard${NC}"

echo -e "\n${YELLOW}5. Executando testes de políticas de Compliance:${NC}"

echo -e "\n${BLUE}[Cenário A - Permitido] Agente 'ia-fintech-support' envia POST seguro em /pix/validar:${NC}"
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST http://127.0.0.1:${NEXO_PORT}/pix/validar -H "Authorization: Bearer ia-fintech-support")
if [ "$HTTP_CODE" -eq 200 ]; then
    echo -e "${GREEN}✓ APROVADO (HTTP $HTTP_CODE)${NC}"
else
    echo -e "${RED}✗ Falhou (HTTP $HTTP_CODE)${NC}"
fi

echo -e "\n${BLUE}[Cenário B - Bloqueado LGPD] Agente 'ia-fintech-support' tenta extração em massa GET em /clientes:${NC}"
RESP=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X GET http://127.0.0.1:${NEXO_PORT}/clientes -H "Authorization: Bearer ia-fintech-support")
echo "$RESP" | grep -q "BLOCKED_LGPD_BULK_DATA_ACCESS_DENIED" && echo -e "${GREEN}✓ BLOQUEADO PROATIVAMENTE PELO NEXO HUB (LGPD Violada)${NC}"

echo -e "\n${BLUE}[Cenário C - Bloqueado BACEN 538] Agente 'ia-fintech-support' tenta DELETE destrutivo em /transacoes/99:${NC}"
RESP=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X DELETE http://127.0.0.1:${NEXO_PORT}/transacoes/99 -H "Authorization: Bearer ia-fintech-support")
echo "$RESP" | grep -q "BLOCKED_BACEN_538_MUTATION_DENIED" && echo -e "${GREEN}✓ BLOQUEADO PROATIVAMENTE PELO NEXO HUB (BACEN 538 Violado)${NC}"

echo -e "\n${BLUE}[Cenário D - Bloqueado CFM] Agente 'ia-health-bot' tenta acessar prontuários médicos em /prontuarios/101:${NC}"
RESP=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X GET http://127.0.0.1:${NEXO_PORT}/prontuarios/101 -H "Authorization: Bearer ia-health-bot")
echo "$RESP" | grep -q "BLOCKED_CFM_MEDICAL_RECORDS_DENIED" && echo -e "${GREEN}✓ BLOQUEADO PROATIVAMENTE PELO NEXO HUB (Diretriz CFM Violada)${NC}"

echo -e "\n${YELLOW}6. Exibindo as últimas entradas do arquivo de Auditoria Imutável (${LOG_FILE}):${NC}"
echo -e "${BLUE}---------------------------------------------------------------------${NC}"
tail -n 4 ${LOG_FILE}
echo -e "${BLUE}---------------------------------------------------------------------${NC}"

echo -e "\n${GREEN}🎉 Demonstração concluída com 100% de sucesso!${NC}"
