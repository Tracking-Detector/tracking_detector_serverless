FROM golang:1.18-alpine AS builder

WORKDIR /app

COPY ./go/go.mod /app/go.mod
COPY ./go/go.sum /app/go.sum
RUN go mod download

COPY ./go/functions/auth/main.go /app/main.go
COPY ./go/shared /app/shared

# Build the function executable with optimizations enabled
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .

CMD ["./main"]