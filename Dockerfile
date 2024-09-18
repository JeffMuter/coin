# syntax=docker/dockerfile:1

# Use the official Golang image
FROM golang:1.23-alpine

# Set the working directory inside the container
WORKDIR /workspace

# Copy go.mod first (if it exists) to avoid redownloading dependencies on every build
COPY go.mod ./

# Download dependencies (this won't do anything if you have no dependencies)
RUN go mod download || true

# Copy the rest of the application code
COPY . .

# Build the Go application from main.go
RUN go build -o main ./main.go

# Expose the port on which your Go server will run
EXPOSE 8080

# Command to run the Go server
CMD ["./main"]
