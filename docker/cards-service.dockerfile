FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata
WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /build/cards-service \
    ./cmd/cards_service

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /build/cards-service .
COPY static ./static

EXPOSE 8082
ENTRYPOINT ["./cards-service"]
