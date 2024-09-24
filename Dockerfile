# Use the official Golang image as a base
FROM golang:1.22-alpine as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules files
COPY go.mod go.sum ./

# Download and cache Go modules dependencies
RUN go mod download

# Copy the entire source code to the working directory
COPY . .

# Build the Go app (binary)
RUN go build -o services-manager-svc cmd/main.go

# Use a smaller image for the final stage
FROM alpine:3.18

# Set environment variables (optional)
ENV GIN_MODE=release

# Expose the port that the app will run on
EXPOSE 8080

# Copy the compiled Go binary from the builder stage
COPY --from=builder /app/auth-manager-svc /auth-manager-svc

# Copy migration files to the final container
COPY ./database/migrations /app/database/migrations

# Run the Go application
ENTRYPOINT ["/auth-manager-svc"]