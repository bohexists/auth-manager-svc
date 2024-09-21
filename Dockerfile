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

RUN ls -la /

RUN chmod +x /auth-manager-svc

# Run the app
ENTRYPOINT ["/auth-manager-svc"]