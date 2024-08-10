# Use an official Golang image as the base
FROM golang:1.22.4-alpine as builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod tidy

# Copy the source code
COPY . .

# Build the Go application
RUN go build -o main .

# Use a minimal image for the final container
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

# Copy .env file
COPY .env .env

# Expose the port the app runs on
EXPOSE 3000

# Command to run the application
CMD ["./main"]
