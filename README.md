# 🚨 Real-Time Alerting System

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Enterprise-grade distributed real-time alerting system built with Go, Kubernetes, and Terraform.

## 📋 Overview

This system provides real-time alerting capabilities with high availability, horizontal scalability, and comprehensive observability. Designed to handle thousands of concurrent users with minimal latency.

### Key Features

- **Real-time Communication**: WebSocket-based bidirectional communication
- **Event-Driven Architecture**: Redis Streams/NATS for reliable message processing
- **Multi-Channel Notifications**: Slack, Email, and SMS integrations
- **High Availability**: Kubernetes-native with auto-scaling
- **Full Observability**: Prometheus metrics, Grafana dashboards, Jaeger tracing
- **Infrastructure as Code**: Terraform modules for AWS deployment

## 🛠️ Tech Stack

| Category | Technologies |
|----------|-------------|
| **Backend** | Go 1.22+, Fiber Framework |
| **Message Broker** | Redis Streams / NATS |
| **Database** | PostgreSQL, Redis |
| **Infrastructure** | Terraform, AWS EKS |
| **Containers** | Docker, Kubernetes |
| **CI/CD** | GitHub Actions |
| **Observability** | Prometheus, Grafana, Jaeger |

## 📁 Project Structure
```
realtime-alerting-system/
├── cmd/                        # Application entrypoints
│   └── api/                    # Main API server
├── internal/                   # Private application code
│   ├── domain/                 # Business logic & entities
│   │   ├── entity/             # Domain models
│   │   ├── repository/         # Repository interfaces (ports)
│   │   └── service/            # Domain service interfaces
│   ├── application/            # Use cases & app services
│   │   ├── dto/                # Data Transfer Objects
│   │   ├── usecase/            # Use case implementations
│   │   └── service/            # Application services
│   ├── infrastructure/         # External implementations
│   │   ├── config/             # Configuration (Viper)
│   │   ├── database/           # PostgreSQL, Redis
│   │   ├── messaging/          # Message broker
│   │   ├── notification/       # Slack, Email, SMS
│   │   └── logger/             # Structured logging
│   └── presentation/           # HTTP & WebSocket layer
│       ├── http/               # REST API handlers
│       └── websocket/          # WebSocket server
├── pkg/                        # Public reusable packages
├── deployments/                # Deployment configurations
│   ├── docker/                 # Additional Dockerfiles
│   └── kubernetes/             # K8s manifests, Helm charts
├── terraform/                  # Infrastructure as Code
│   ├── modules/                # Reusable Terraform modules
│   └── environments/           # Environment configs
├── migrations/                 # Database migrations
├── scripts/                    # Automation scripts
├── docs/                       # Documentation
└── test/                       # Integration/E2E tests
```

## 🚀 Quick Start

### Prerequisites

- Go 1.22 or higher
- Docker & Docker Compose
- Make
- golangci-lint

### Installation

1. **Clone the repository**
```bash
   git clone https://github.com/TU_USUARIO/realtime-alerting-system.git
   cd realtime-alerting-system
```

2. **Install dependencies**
```bash
   make deps
```

3. **Setup environment**
```bash
   make env
   # Edit .env with your configuration
```

4. **Run locally**
```bash
   make run
```

### Using Docker Compose
```bash
# Start all services
make docker-up

# View logs
make docker-logs

# Stop services
make docker-down
```

## 📝 Available Commands

Run `make help` to see all available commands:
```
Usage:
  make <target>

Targets:
  help            Show this help message
  run             Run the application
  build           Build the application binary
  clean           Clean build artifacts
  test            Run tests
  test-coverage   Run tests with coverage
  lint            Run linter
  lint-fix        Run linter and fix issues
  fmt             Format code
  check           Run all code quality checks
  deps            Download dependencies
  tidy            Tidy dependencies
  docker-build    Build Docker image
  docker-up       Start all services with docker-compose
  docker-down     Stop all services
  migrate-up      Run database migrations up
  migrate-down    Run database migrations down
```

## ⚙️ Configuration

The application supports multiple configuration sources (in order of priority):

1. Environment variables
2. `.env` file
3. `config.yaml` file
4. Default values

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `APP_ENV` | Environment (development/staging/production) | development |
| `SERVER_PORT` | HTTP server port | 8080 |
| `DATABASE_HOST` | PostgreSQL host | localhost |
| `DATABASE_PORT` | PostgreSQL port | 5432 |
| `DATABASE_USER` | PostgreSQL user | postgres |
| `DATABASE_PASSWORD` | PostgreSQL password | postgres |
| `DATABASE_NAME` | Database name | alerting_db |
| `REDIS_HOST` | Redis host | localhost |
| `REDIS_PORT` | Redis port | 6379 |
| `JWT_SECRET` | JWT signing secret | - |
| `LOG_LEVEL` | Log level (debug/info/warn/error) | debug |

## 🧪 Testing
```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run short tests only
make test-short
```

## 🔍 Code Quality
```bash
# Run linter
make lint

# Run linter and auto-fix
make lint-fix

# Format code
make fmt

# Run all checks (fmt, vet, lint)
make check
```

## 📊 Architecture

The system follows **Clean Architecture** principles with clear separation of concerns:
```
┌─────────────────────────────────────────────────────────────┐
│                      Presentation                            │
│                   (HTTP, WebSocket)                          │
├─────────────────────────────────────────────────────────────┤
│                      Application                             │
│                (Use Cases, Services)                         │
├─────────────────────────────────────────────────────────────┤
│                        Domain                                │
│              (Entities, Business Rules)                      │
├─────────────────────────────────────────────────────────────┤
│                    Infrastructure                            │
│          (Database, Cache, External APIs)                    │
└─────────────────────────────────────────────────────────────┘
```

## 🗺️ Roadmap

- [x] Project structure setup
- [x] Configuration management
- [x] Linting & code quality
- [ ] Docker Compose environment
- [ ] Domain entities & repositories
- [ ] REST API endpoints
- [ ] WebSocket server
- [ ] Event-driven messaging
- [ ] Observability stack
- [ ] CI/CD pipeline
- [ ] Terraform infrastructure
- [ ] Kubernetes deployment

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Commit Convention

This project follows [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` New feature
- `fix:` Bug fix
- `docs:` Documentation changes
- `chore:` Maintenance tasks
- `refactor:` Code refactoring
- `test:` Adding tests

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👤 Author

**Daniel**

- GitHub: [@daniel-caso-github](https://github.com/daniel-caso-github)

---

⭐ If you find this project useful, please consider giving it a star!
