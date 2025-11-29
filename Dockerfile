# syntax=docker/dockerfile:1
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/app/main.go

FROM golang:1.25-alpine AS test
WORKDIR /app
COPY --from=builder /app /app
CMD ["go", "test", "./...", "-v", "-cover"]

FROM alpine:latest AS runtime
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
COPY internal/db/migrations ./migrations
EXPOSE 8080
CMD ["./main"]
