curl -X POST http://localhost:8080/sessions/add -H "Content-Type: application/json" -d '{"name": "minha-sessao"}'

curl -X POST http://localhost:8080/sessions/minha-sessao/connect

curl -s -X POST http://localhost:8080/message/14820c8d-50e2-42ed-a0bb-645d1b083bf7/send/text \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "559981769536",
    "body": "🚀 Teste da API WazMeow! Esta é uma mensagem de teste enviada via API REST."
  }' | jq .