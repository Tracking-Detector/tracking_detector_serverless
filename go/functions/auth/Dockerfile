# Use the official Golang image as the base image
FROM golang:1.18-alpine

WORKDIR /app
# # Copy the Go module files
# COPY go.mod /app/go.mod

# Download and cache Go modules
COPY ./go/go.mod /app/go.mod
COPY ./go/go.sum /app/go.sum

RUN go mod download
# Copy the function source code and shared library
COPY ./go/functions/auth/main.go /app/main.go
COPY ./go/shared /app/shared

# Build the function executable
RUN go build -o main .

# Set the command to run the function executable
CMD ["./main"]