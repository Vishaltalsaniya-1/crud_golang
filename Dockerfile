# Use the official Golang image
FROM golang:1.23-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files for dependency management
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project into the container
COPY . .

# Build the Go app
RUN go build -o main .

# Expose port 8081 for the app
EXPOSE 8081

# Command to run the app
CMD ["./main"]
