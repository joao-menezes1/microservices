FROM golang:1.25-alpine AS builder
WORKDIR /app

# Copia TUDO o que está na raiz (Projeto-Ruan) para dentro do container
# Precisamos disso por causa do 'replace' no go.mod
COPY . .

# Entra na pasta do serviço específico para buildar
WORKDIR /app/microservices/order
RUN go mod download
RUN go build -o main cmd/main.go

FROM alpine:latest
WORKDIR /root/
# Copia o binário gerado
COPY --from=builder /app/microservices/order/main .
CMD ["./main"]