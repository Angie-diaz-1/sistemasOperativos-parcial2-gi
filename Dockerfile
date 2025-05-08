# Etapa 1: Compilar
FROM golang:1.24-alpine AS builder

#LABEL about the custom image
LABEL maintainer="AngieDiaz" \
      version="0.1" \
      description="This is a custom Docker image for the Apache services"


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
