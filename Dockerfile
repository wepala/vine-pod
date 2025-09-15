# Build stage
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Install git for version information
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod ./
COPY go.su[m] ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary with version information
ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_TIME=unknown

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags "-X github.com/wepala/vine-pod/pkg/version.Version=${VERSION} \
              -X github.com/wepala/vine-pod/pkg/version.Commit=${COMMIT} \
              -X github.com/wepala/vine-pod/pkg/version.BuildTime=${BUILD_TIME}" \
    -o vine-pod ./cmd/vine-pod

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1001 -S vine && \
    adduser -u 1001 -S vine -G vine

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/vine-pod .

# Create data directory
RUN mkdir -p /app/data && chown vine:vine /app/data

# Switch to non-root user
USER vine

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the binary
CMD ["./vine-pod"]