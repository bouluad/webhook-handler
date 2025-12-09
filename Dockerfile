# --- Stage 1: Build the Go application ---
FROM golang:1.21-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the go mod and sum files
COPY go.mod go.sum ./

# Download dependencies. Dependencies will be cached if the go.mod and go.sum files are unchanged
RUN go mod download

# Copy the source code
COPY . .

# Build the application. CGO_ENABLED=0 creates a statically linked binary.
# The binary is named 'webhook-handler'
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o /webhook-handler ./cmd/webhook-handler

# --- Stage 2: Create the final lean image ---
FROM alpine:3.18

# Set the environment variable for security
ENV PORT=8080

# Expose the port the application runs on
EXPOSE 8080

# Copy the static binary from the builder stage
COPY --from=builder /webhook-handler /webhook-handler

# Set the entry point to run the application
ENTRYPOINT ["/webhook-handler"]
