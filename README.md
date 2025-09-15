# Vine Pod

A Solid Server implementation in Go following standard microservice architecture patterns.

## Overview

Vine Pod is a modern, cloud-native implementation of a Solid server built with Go. It follows standard Go project layout conventions and microservice best practices.

## Features

- 🚀 **Standard Go Project Layout** - Clean, maintainable code structure
- 🔧 **Configuration Management** - Environment-based configuration
- 📊 **Structured Logging** - JSON logging with multiple levels
- 🛡️ **Middleware Support** - CORS, logging, recovery middleware
- 🐳 **Docker Ready** - Multi-stage Docker builds
- 🔄 **Graceful Shutdown** - Proper signal handling
- ✅ **Health Checks** - Built-in health and version endpoints
- 🧪 **Test Coverage** - Comprehensive test suite

## Quick Start

### Prerequisites

- Go 1.24+ 
- Docker (optional)
- Make (optional, for convenience commands)

### Build and Run

```bash
# Clone the repository
git clone https://github.com/wepala/vine-pod.git
cd vine-pod

# Build the application
make build

# Run the application
make run
```

### Using Docker

```bash
# Build and run with Docker
make docker-run
```

### Development

```bash
# Setup development environment
make dev-setup

# Run tests
make test

# Run with coverage
make test-cover

# Format code
make fmt

# Lint code
make lint
```

## Configuration

Configure the application using environment variables:

```bash
export SERVER_HOST=0.0.0.0
export SERVER_PORT=8080
export LOG_LEVEL=info
export SOLID_DATA_PATH=./data
```

See [API Documentation](docs/API.md) for complete configuration options.

## Project Structure

```
├── cmd/vine-pod/          # Application entry point
├── internal/              # Private application code
│   ├── app/              # Application logic
│   ├── config/           # Configuration management
│   ├── handler/          # HTTP handlers
│   ├── middleware/       # HTTP middleware
│   └── server/           # HTTP server setup
├── pkg/                  # Public library code
│   ├── logger/           # Logging utilities
│   └── version/          # Version information
├── api/v1/               # API definitions (future)
├── deployments/          # Deployment configurations
├── docs/                 # Documentation
├── scripts/              # Build and deployment scripts
└── test/                 # Additional test files
```

## API Endpoints

- `GET /health` - Health check
- `GET /version` - Version information  
- `GET /` - Welcome message
- `ALL /solid/*` - Solid protocol endpoints (placeholder)

See [API Documentation](docs/API.md) for detailed endpoint information.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run the test suite
6. Submit a pull request

## License

This project is licensed under the GNU Affero General Public License v3.0 - see the [LICENSE](LICENSE) file for details.
