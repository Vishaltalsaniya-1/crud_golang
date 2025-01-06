# Use the official Golang image
FROM golang:1.23-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files for dependency management
COPY go.mod go.sum ./

# Download dependencies
RUN go mod tidy

# Copy the entire project into the container
COPY . .

# Build the Go app
RUN go build -o main .

# Use a minimal base image (Alpine) for the final image
FROM alpine:latest

#  Install dependencies for PostgreSQL and MongoDB client if needed
RUN apk --no-cache add ca-certificates

#  Copy the Go application from the builder stage
COPY --from=builder /app/main .

# Expose port 8081 for the app
EXPOSE 8081

# Command to run the app
CMD ["./main"]
