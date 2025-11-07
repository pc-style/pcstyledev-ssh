# Build stage
FROM golang:latest-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git openssh-keygen

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o ssh-server ./cmd/server

# Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates openssh-keygen

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /build/ssh-server .

# Generate SSH host keys if they don't exist
RUN mkdir -p .ssh && \
    if [ ! -f .ssh/id_ed25519 ]; then \
        ssh-keygen -t ed25519 -f .ssh/id_ed25519 -N "" -C "ssh-server-host-key"; \
    fi

# Expose SSH port
EXPOSE 2222

# Run the server
CMD ["./ssh-server"]
