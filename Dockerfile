# Build stage
FROM golang:1.27.2-alpine AS builder

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /server ./cmd/server

# Runtime stage
FROM alpine:3.24.2

RUN apk --no-cache add ca-certificates && \
    adduser -D -g '' appuser

WORKDIR /app

COPY --from=builder /server .

USER appuser

EXPOSE 8080

ENTRYPOINT ["./server"]