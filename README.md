# KKT 54-FZ Monitoring

![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)
![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go)

A comprehensive monitoring system for Cash Register Equipment (KKT) compliant with Russian Federal Law 54-FZ "On the use of cash register equipment".

**Author's brand / website:** https://run-as-daemon.ru

## Description

`kkt-54fz-monitoring` is a complete solution for monitoring and analyzing cash register equipment operating under 54-FZ. The system collects data from various sources (logs, OFD HTTP API), aggregates metrics, and provides them to Prometheus. It includes a built-in AI subsystem for error clustering and alert recommendations.

## Features

- 📊 **Metric Collection** from file logs and OFD HTTP API
- 📈 **Prometheus Exporter** with ready-to-use metrics
- 🚨 **Pre-configured Alert Rules** for common issues
- 📉 **Grafana Dashboards** for visualization
- 🤖 **AI Subsystem** for error analysis and recommendations
- ⚙️ **Flexible Configuration** via YAML
- 🔒 **Security** - built-in checks and validation

## Quick Start

### Installation

```bash
# Clone the repository
git clone https://github.com/ranas-mukminov/kkt-54fz-monitoring.git
cd kkt-54fz-monitoring

# Build
make build

# Or with Go
go build -o kkt-monitor ./cmd/kkt-monitor
```

### Running

```bash
# Run with default configuration
./kkt-monitor --config configs/config.yaml

# Run with Docker
docker-compose up -d
```

### Check Metrics

```bash
curl http://localhost:9090/metrics
```

## Architecture

```
┌─────────────┐     ┌─────────────┐
│  File Logs  │     │  HTTP OFD   │
└──────┬──────┘     └──────┬──────┘
       │                   │
       └───────┬───────────┘
               │
        ┌──────▼──────┐
        │ Collectors  │
        └──────┬──────┘
               │
        ┌──────▼──────┐
        │   Domain    │
        │    Model    │
        └──────┬──────┘
               │
        ┌──────▼──────┐
        │  Prometheus │
        │   Exporter  │
        └──────┬──────┘
               │
               ├────────┐
               │        │
        ┌──────▼──────┐ │
        │ Prometheus  │ │
        └─────────────┘ │
                        │
                 ┌──────▼──────┐
                 │   Grafana   │
                 └─────────────┘
```

## Configuration

Example configuration file `configs/config.yaml`:

```yaml
server:
  port: 9090
  metrics_path: /metrics

collectors:
  file_log:
    enabled: true
    path: /var/log/kkt/*.log
    format: json
    poll_interval: 10s
  
  http_ofd:
    enabled: true
    url: https://ofd.example.ru/api/v1
    api_key: ${OFD_API_KEY}
    poll_interval: 30s

ai:
  provider: mock  # mock, openai, anthropic
  error_clustering:
    enabled: true
    min_cluster_size: 5
  alert_advisor:
    enabled: true

logging:
  level: info
  format: json
```

## Metrics

The system exports the following metrics:

- `kkt_status` - KKT status (0=unavailable, 1=running, 2=error)
- `kkt_documents_total` - total number of fiscal documents
- `kkt_errors_total` - number of errors by type
- `kkt_ofd_sync_status` - OFD synchronization status
- `kkt_shift_status` - shift status (open/closed)
- `kkt_last_document_timestamp` - timestamp of last document

## Alerts

Pre-configured alert rules are in `configs/alerts/kkt-alerts.yaml`:

- KKT unavailable for more than 5 minutes
- Critical fiscal drive error
- OFD synchronization issues
- Fiscal drive memory overflow
- Fiscal drive expiration

## Development

### Requirements

- Go 1.23+
- Make
- Docker and Docker Compose (for local development)

### Build and Test

```bash
# Install dependencies
make deps

# Run linter
make lint

# Run tests
make test

# Run integration tests
make test-integration

# Security check
make security-check

# Performance check
make perf-check

# Full check (lint + test + security)
make check
```

### Project Structure

```
.
├── cmd/
│   └── kkt-monitor/        # Application entry point
├── internal/
│   ├── domain/             # Domain models
│   ├── config/             # Configuration loading and validation
│   ├── collector/          # Data collectors
│   ├── exporter/           # Prometheus exporter
│   └── ai/                 # AI subsystem
├── pkg/
│   ├── utils/              # Utilities
│   └── logger/             # Logging
├── configs/
│   ├── config.yaml         # Main configuration
│   ├── alerts/             # Alert rules
│   └── dashboards/         # Grafana dashboards
├── deployments/
│   ├── docker/             # Docker files
│   └── kubernetes/         # K8s manifests
├── test/
│   ├── testdata/           # Test data
│   └── integration/        # Integration tests
└── docs/                   # Documentation
```

## AI Subsystem

### Error Clustering

The AI module automatically groups similar errors for easier analysis:

```bash
curl http://localhost:9090/api/v1/ai/error-clusters
```

### Alert Advisor

Get alert configuration recommendations based on historical data:

```bash
curl http://localhost:9090/api/v1/ai/alert-recommendations
```

### Providers

Supported AI providers:
- **mock** - stub for development and testing
- **openai** - OpenAI GPT API
- **anthropic** - Anthropic Claude API

## 54-FZ Compliance

The system complies with:
- Federal Law No. 54-FZ dated 22.05.2003
- Order of the Federal Tax Service of Russia dated 21.03.2017 No. MMV-7-20/229@
- Technical requirements for fiscal document formats

See [LEGAL](LEGAL) file for detailed compliance information.

## License

Apache License 2.0. See [LICENSE](LICENSE) file for details.

## Support

- 📧 Email: support@run-as-daemon.ru
- 🐛 Issues: https://github.com/ranas-mukminov/kkt-54fz-monitoring/issues
- 📖 Documentation: https://github.com/ranas-mukminov/kkt-54fz-monitoring/wiki

## Author

© 2024 [run-as-daemon.ru](https://run-as-daemon.ru)

---

[Русская версия](README.ru.md)
