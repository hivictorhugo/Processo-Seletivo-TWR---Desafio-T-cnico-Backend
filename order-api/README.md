# Affiliate Conversions - Final

## O que está incluso
- API em Go para receber notificações de conversão.
- HMAC por parceiro (`X-Partner-Id` + `X-Signature`).
- Persistência em MySQL (tabelas `partners` e `conversions`).
- Idempotência via `transaction_id` único.
- Docker Compose para rodar tudo localmente.
- Teste básico para criar e detectar duplicidade.

## Como rodar local (Docker)
1. Copie `.env.example` se houver e ajuste variáveis.
2. `docker-compose up --build` (o MySQL inicializa e roda `db/init.sql`).
3. API em `http://localhost:8080/conversions`.

## Exemplo de requisição
Veja `db/init.sql` — existe um parceiro `partner-a` com `secret_key` = `secret_for_partner_a`.

```bash
BODY='{"transaction_id":"tx_123","partner_name":"Partner A","sale_amount":199.90}'
SIGNATURE=$(echo -n "$BODY" | openssl dgst -sha256 -hmac "secret_for_partner_a" -hex | sed 's/^.* //')

curl -X POST http://localhost:8080/conversions \
  -H "Content-Type: application/json" \
  -H "X-Partner-Id: partner-a" \
  -H "X-Signature: $SIGNATURE" \
  -d "$BODY"
```

## Próximos passos recomendados (extras para produção)
- Rotacionar chaves e armazená-las cifradas.
- TLS/mTLS e rate limiting.
- Colocar worker assíncrono para processamento e fila (RabbitMQ/Redis).
- Monitoramento e alertas.
- Testes de integração automatizados no pipeline CI.