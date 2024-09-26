# Etapa de build
FROM golang:1.23-alpine AS build
WORKDIR /app
COPY . .

# Instala o certificado SSL necessário
RUN apk add --no-cache ca-certificates && update-ca-certificates

# Correção para o comando de build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o web-service ./cmd/main.go

# Etapa final
FROM alpine:latest
WORKDIR /app
COPY --from=build /app/web-service .
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ 
ENTRYPOINT ["./web-service"]