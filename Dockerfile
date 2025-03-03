# Use the official Golang image
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o your-binary-name

# Use a lightweight Alpine image for the final stage
FROM alpine:latest

# Copy the binary from the builder stage
COPY --from=builder /app/your-binary-name /your-binary-name

# Expose the port your app runs on
EXPOSE 8080

# Run the application
CMD ["/simple"]