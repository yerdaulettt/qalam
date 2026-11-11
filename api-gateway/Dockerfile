FROM golang:1.27.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o api ./cmd/main.go

FROM alpine:3.23.3
COPY --from=builder app/api /api
COPY --from=builder app/config.yaml config.yaml

ENTRYPOINT ["/api"]