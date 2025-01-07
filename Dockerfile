# Stage 1: Build the Go application
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Install Git to clone the repository
RUN apk add --no-cache git

# Clone the repository (replace with your repository URL)
RUN git clone https://github.com/Vishaltalsaniya-1/crud_golang.git .

# Download dependencies
RUN go mod download

# Build the Go app
RUN go build -o main .

# Stage 2: Create a minimal runtime image
FROM alpine:latest

# Install certificates for secure connections (e.g., to databases)
RUN apk --no-cache add ca-certificates

# Copy the Go application binary from the builder stage
COPY --from=builder /app/main /app/main

# Set the working directory in the runtime container
WORKDIR /app

# Expose the application port (update as needed)
EXPOSE 8081

# Command to run the app
CMD ["./main"]
