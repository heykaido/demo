# Use the official Golang image as the build environment
FROM golang:1.23.2 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Download and cache Go modules
RUN go mod download

# Copy the entire project directory to the working directory
COPY . .

# Set the working directory to where main.go is located
WORKDIR /app/cmd/go

# Build the Go application
RUN go build -o /app/ev-api

# Use a minimal base image for the final build
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /root/

# Copy the built Go binary from the builder stage
COPY --from=builder /app/ev-api .

# Expose the port the application will run on
EXPOSE 8080

# Command to run the application
CMD ["./ev-api"]
