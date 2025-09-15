# API Documentation

## Overview

Vine Pod is a Solid Server implementation in Go that provides a microservice architecture for handling Solid protocol requests.

## Endpoints

### Health Check

**GET** `/health`

Returns the health status of the service.

**Response:**
```json
{
  "status": "healthy",
  "service": "vine-pod"
}
```

### Version Information

**GET** `/version`

Returns version and build information.

**Response:**
```json
{
  "version": "v1.0.0",
  "commit": "abc123...",
  "build_time": "2024-01-01T00:00:00Z",
  "go_version": "go1.24.7",
  "platform": "linux/amd64"
}
```

### Root

**GET** `/`

Returns welcome message and available endpoints.

**Response:**
```json
{
  "message": "Welcome to Vine Pod - Solid Server",
  "version": "v1.0.0",
  "endpoints": {
    "health": "/health",
    "version": "/version",
    "solid": "/solid/"
  }
}
```

### Solid Protocol (Placeholder)

**ALL** `/solid/*`

Placeholder endpoint for Solid protocol implementation.

**Response:**
```json
{
  "message": "Solid protocol endpoint",
  "method": "GET",
  "path": "/solid/test",
  "status": "not_implemented",
  "note": "This endpoint will be implemented with full Solid protocol support"
}
```

## Configuration

The service can be configured using environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_HOST` | `0.0.0.0` | Server bind address |
| `SERVER_PORT` | `8080` | Server port |
| `SERVER_READ_TIMEOUT` | `30s` | HTTP read timeout |
| `SERVER_WRITE_TIMEOUT` | `30s` | HTTP write timeout |
| `SERVER_IDLE_TIMEOUT` | `60s` | HTTP idle timeout |
| `LOG_LEVEL` | `info` | Log level (debug, info, warn, error) |
| `SOLID_DATA_PATH` | `./data` | Path to Solid data storage |
| `SOLID_ALLOW_ORIGIN` | `*` | CORS allow origin header |
| `SOLID_ENABLE_CORS` | `true` | Enable CORS middleware |

## Examples

### Start the server

```bash
# Using default configuration
./vine-pod

# Using environment variables
export SERVER_PORT=9000
export LOG_LEVEL=debug
./vine-pod
```

### Docker

```bash
# Build image
docker build -t vine-pod .

# Run container
docker run -p 8080:8080 vine-pod

# Run with custom configuration
docker run -p 9000:8080 -e SERVER_PORT=8080 -e LOG_LEVEL=debug vine-pod
```