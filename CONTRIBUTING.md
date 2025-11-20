# Contributing to KKT 54-FZ Monitoring

Thank you for your interest in contributing to KKT 54-FZ Monitoring! This document provides guidelines for contributing to the project.

## Code of Conduct

Be respectful, professional, and inclusive. We welcome contributions from everyone.

## How to Contribute

### Reporting Bugs

If you find a bug, please create an issue with:
- Clear description of the problem
- Steps to reproduce
- Expected vs actual behavior
- Environment details (OS, Go version, etc.)
- Relevant logs or error messages

### Suggesting Features

For feature requests:
- Describe the feature and its use case
- Explain why it would be valuable
- Provide examples if possible

### Pull Requests

1. **Fork the repository**
2. **Create a feature branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **Make your changes**
   - Follow the coding standards
   - Add tests for new functionality
   - Update documentation as needed

4. **Run tests and linters**
   ```bash
   make test
   make lint
   make security-check
   ```

5. **Commit your changes**
   - Use clear, descriptive commit messages
   - Reference issues if applicable

6. **Push and create a PR**
   ```bash
   git push origin feature/your-feature-name
   ```

## Development Setup

### Prerequisites

- Go 1.21 or higher
- Make
- Docker and Docker Compose (optional)

### Getting Started

```bash
# Clone the repository
git clone https://github.com/ranas-mukminov/kkt-54fz-monitoring.git
cd kkt-54fz-monitoring

# Install dependencies
make deps

# Build the project
make build

# Run tests
make test

# Run the application
./build/kkt-monitor --config configs/config.yaml
```

## Coding Standards

### Go Style

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` and `goimports` for formatting
- Keep functions small and focused
- Write clear, self-documenting code

### Testing

- Write unit tests for all new code
- Maintain or improve code coverage
- Use table-driven tests where appropriate
- Mock external dependencies

### Documentation

- Document all exported functions and types
- Update README.md for user-facing changes
- Add architecture documentation for significant changes

### Commit Messages

Format:
```
<type>: <subject>

<body>

<footer>
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `test`: Test changes
- `refactor`: Code refactoring
- `perf`: Performance improvements
- `chore`: Build/tooling changes

Example:
```
feat: add OpenAI provider for AI subsystem

Implement OpenAI GPT integration for error clustering
and alert recommendations. Includes configuration
options and error handling.

Closes #123
```

## Project Structure

```
.
├── cmd/                 # Application entry points
├── internal/            # Private application code
│   ├── domain/         # Domain models
│   ├── config/         # Configuration
│   ├── collector/      # Data collectors
│   ├── exporter/       # Prometheus exporter
│   └── ai/             # AI subsystem
├── pkg/                # Public libraries
├── configs/            # Configuration files
├── test/               # Test files and data
└── docs/               # Documentation
```

## Testing Guidelines

### Unit Tests

- Test individual functions and methods
- Use mocks for external dependencies
- Aim for >80% code coverage

### Integration Tests

- Test component interactions
- Use real dependencies where possible
- Tag with `// +build integration`

### Running Tests

```bash
# All tests
make test

# With coverage
make test-coverage

# Integration tests only
make test-integration
```

## Security

- Never commit secrets or credentials
- Use environment variables for sensitive data
- Run security scans before submitting PRs
- Report security issues privately

## License

By contributing, you agree that your contributions will be licensed under the Apache License 2.0.

## Questions?

Feel free to:
- Open an issue for discussion
- Contact the maintainers
- Check existing issues and PRs

Thank you for contributing!
