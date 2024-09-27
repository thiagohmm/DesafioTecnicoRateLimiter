Desafio Técnico Rate Limiter
Este projeto implementa um sistema de limitação de taxa (rate limiting) para controlar o acesso a uma API.

Como executar
Para rodar o projeto, utilize o seguinte comando:

Bash
docker compose up --build


Testando o Rate Limiter
Você pode testar o funcionamento do rate limiter utilizando os seguintes comandos bash:

Sem API Key:

Bash
for i in {1..100}; do
  curl http://localhost:8080/
done


Com API Key:

Bash
for i in {1..100}; do
  curl -v -H "API_KEY: your_api_key" http://localhost:8080
done


