# Use the official Golang image as a base
FROM golang:1.22-alpine as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules files
COPY go.mod go.sum ./

# Copy the .env file to the container
COPY .env /app/.env

# Download and cache Go modules dependencies
RUN go mod download

# Copy the entire source code to the working directory
COPY . .

# Copy the .env file to the container after copying the source code
COPY .env /app/.env

# Build the Go app (binary) and target the ARM architecture
RUN GOARCH=arm64 go build -o auth-manager-svc cmd/main.go

# List the files to ensure binary was created
RUN ls -la /app

# Create a smaller final image using the scratch base
FROM alpine:3.18

# Set environment variables (optional)
ENV GIN_MODE=release

# Expose the port that the app will run on
EXPOSE 8080

# Copy the compiled binary from the builder stage to the final image
COPY --from=builder /app/auth-manager-svc /auth-manager-svc

# Copy migration files to the final container
COPY ./database/migrations /app/database/migrations

# Run the app
ENTRYPOINT ["/auth-manager-svc"]