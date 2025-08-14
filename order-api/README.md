# Affiliate Conversions API

Este projeto implementa uma **API de Notificação de Conversões de Afiliados** em **Go**, com **MySQL** e suporte a **Docker**, garantindo:

- Recebimento seguro de notificações de conversão
- Autenticação via HMAC
- Persistência em banco MySQL
- Idempotência (evita duplicidade de transações)

---

## Estrutura do Projeto

/cmd/server/main.go
/internal/db/db.go
/internal/handlers/conversion_handler.go
/internal/models/models.go
/internal/middleware/middleware.go
/internal/utils/hmac.go
/tests/conversion_handler_test.go
/db/init.sql
/Dockerfile
/docker-compose.yml
/pkg/config/config.go
/README.md
/go.mod
/go.sum



---

## Pré-requisitos

- Go >= 1.21
- Docker e Docker Compose
- PowerShell ou terminal Unix (para testar a API)
- Git (para subir no GitHub)

---

## Configuração

O projeto utiliza **variáveis de ambiente**:

```env
DB_DSN=appuser:apppass@tcp(db:3306)/affiliate_db?parseTime=true&charset=utf8mb4
PORT=8080
```
O `config.go` lê automaticamente as variáveis de ambiente ao iniciar a API.

O `docker-compose.yml` já define:

- **Serviço db (MySQL)**
- **Serviço app (API)**
- Mapeamento de portas: `8080` e `3306`
- Volume persistente para o banco de dados

---

## Rodando o Projeto

Subir o banco e a API pelo Docker:

```bash
docker-compose up --build
```
Confirmar que os containers estão rodando:

```bash
docker ps
```
`order-api-db-1 `→ MySQL

`order-api-app-1 `→ API

Ver logs da API:

```bash
docker-compose logs -f
```
## Testando a API

Criar conversão (parceiro válido)
No PowerShell:

```powershell

$body = '{"transaction_id":"tx_test_1","partner_name":"Partner A","sale_amount":99.90}'
$mac = New-Object System.Security.Cryptography.HMACSHA256
$mac.Key = [Text.Encoding]::UTF8.GetBytes("secret_for_partner_a")
$sig = [System.BitConverter]::ToString($mac.ComputeHash([Text.Encoding]::UTF8.GetBytes($body))).Replace("-", "").ToLower()

Invoke-RestMethod -Uri "http://localhost:8080/conversions" -Method Post -Body $body -Headers @{
    "X-Partner-Id" = "partner-a"
    "X-Signature" = $sig
} -ContentType "application/json"
```
Resposta esperada:

```json
{"status":"created"}
```

**Testar duplicidade (mesma transação)**
Rodar novamente o mesmo comando deve retornar:

```json
{"status":"duplicate"}
```
**Testar parceiro inválido ou assinatura incorreta**
Parceiro desconhecido → `unknown partner`

Assinatura inválida → `invalid signature`

**Testes automatizados**
```bash
go test ./...
```
**Considerações**

- Idempotência garantida pelo tratamento de duplicidade no MySQL (erro 1062).

- Segurança via HMAC por parceiro.

- Estrutura modular seguindo boas práticas Go.

- Docker garante ambiente consistente e fácil deploy.

## Autor
**Victor Hugo**