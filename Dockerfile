# Etapa 1: Compilar
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copiar solo lo necesario
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

# Etapa 2: Ejecutar
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .
COPY .env .

EXPOSE 8080

CMD ["./main"]
