# Start from the official Golang image
FROM golang:1.23-alpine AS builder

# Set working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Start a new stage from scratch
FROM alpine:latest  

# Install necessary certificates and timezone data
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/main .
COPY --from=builder /app/config ./config
COPY --from=builder /app/config/app.conf ./config/app.conf

# Optional: Create a default .env file
RUN touch .env

# Expose port 8087 (your application port)
EXPOSE 8087

# Command to run the executable
CMD ["./main"]