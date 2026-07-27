#!/usr/bin/env python3
"""
Simulador de API Bancária Protegida para Gravação de Vídeo & Demonstrações do Aegis Core
"""

import http.server
import socketserver
import json
import sys

PORT = 9000

class BankingBackendHandler(http.server.SimpleHTTPRequestHandler):
    def _send_json(self, status, payload):
        self.send_response(status)
        self.send_header('Content-Type', 'application/json; charset=utf-8')
        self.end_headers()
        self.wfile.write(json.dumps(payload, ensure_ascii=False).encode('utf-8'))

    def do_GET(self):
        if "/clientes" in self.path:
            self._send_json(200, {
                "system": "Core Banking",
                "total_customers": 15420,
                "data": [
                    {"id": 1, "name": "João Silva", "cpf": "123.456.789-00", "balance": 15000.00},
                    {"id": 2, "name": "Maria Santos", "cpf": "987.654.321-11", "balance": 42000.50}
                ]
            })
        elif "/prontuarios" in self.path:
            self._send_json(200, {
                "system": "Hospital Health System",
                "patient": "Carlos Oliveira",
                "diagnosis": "CID-10 J45.0 Asthma",
                "records_unlocked": True
            })
        else:
            self._send_json(200, {"status": "ok", "message": "API Interna Protegida operando normalmente"})

    def do_POST(self):
        self._send_json(200, {
            "status": "APPROVED",
            "transaction_id": "TX-987654321",
            "message": "Operação de PIX realizada com sucesso no Core Banking"
        })

    def do_DELETE(self):
        self._send_json(200, {
            "status": "DELETED",
            "message": "Conta bancária encerrada no sistema central"
        })

if __name__ == "__main__":
    print(f"🏦 [MOCK BANKING API] Subindo Core Banking na porta {PORT}...")
    with socketserver.TCPServer(('127.0.0.1', PORT), BankingBackendHandler) as httpd:
        try:
            httpd.serve_forever()
        except KeyboardInterrupt:
            print("\nEncerrando Mock Banking API.")
            sys.exit(0)
