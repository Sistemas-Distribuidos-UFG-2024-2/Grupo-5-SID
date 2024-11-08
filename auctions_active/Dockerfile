# Etapa de build
FROM golang:1.21-alpine AS builder

RUN apk update && apk add --no-cache gcc musl-dev

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o main .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .
EXPOSE 6003

CMD ["./main"]
