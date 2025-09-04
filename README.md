[![GitHub release](https://img.shields.io/github/release/sgaunet/http-echo.svg)](https://github.com/sgaunet/http-echo/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/sgaunet/http-echo)](https://goreportcard.com/report/github.com/sgaunet/http-echo)
![GitHub Downloads](https://img.shields.io/github/downloads/sgaunet/http-echo/total)
![Coverage Badge](https://raw.githubusercontent.com/wiki/sgaunet/http-echo/coverage-badge.svg)
[![linter](https://github.com/sgaunet/http-echo/actions/workflows/coverage.yml/badge.svg)](https://github.com/sgaunet/http-echo/actions/workflows/coverage.yml)
[![coverage](https://github.com/sgaunet/http-echo/actions/workflows/coverage.yml/badge.svg)](https://github.com/sgaunet/http-echo/actions/workflows/coverage.yml)
[![Snapshot Build](https://github.com/sgaunet/http-echo/actions/workflows/snapshot.yml/badge.svg)](https://github.com/sgaunet/http-echo/actions/workflows/snapshot.yml)
[![Release Build](https://github.com/sgaunet/http-echo/actions/workflows/release.yml/badge.svg)](https://github.com/sgaunet/http-echo/actions/workflows/release.yml)
[![GoDoc](https://godoc.org/github.com/sgaunet/http-echo?status.svg)](https://godoc.org/github.com/sgaunet/http-echo)
[![License](https://img.shields.io/github/license/sgaunet/http-echo.svg)](LICENSE)

# http-echo

A comprehensive HTTP echo server that provides detailed, structured information about incoming HTTP requests. Perfect for debugging, testing, and learning about HTTP requests in development, Docker, and Kubernetes environments.

## ✨ Features

- **📊 Structured Output**: Clean, organized request information in logical sections
- **⏱️ Performance Metrics**: Request processing time and timestamps
- **🔍 Smart Parsing**: Automatic JSON formatting and form data parsing
- **🌐 Network Analysis**: Real client IP detection (proxy-aware)
- **📱 Environment Info**: Server details, Go version, container/K8s information
- **🎯 Header Analysis**: Important headers highlighted, complete header listing
- **📝 Multiple Formats**: Support for JSON, form data, query parameters
- **🐳 Container Ready**: Optimized for Docker and Kubernetes deployments

## 🚀 Quick Start

### Using Docker Compose
```bash
$ docker-compose up -d
$ curl http://localhost:8080/hello?param=value
```

### Using Docker directly
```bash
$ docker run -p 8080:8080 ghcr.io/sgaunet/http-echo:latest
```

### Build from source
```bash
$ go build .
$ ./http-echo
```

## 📖 Example Output

### Simple GET Request
```bash
$ curl "http://localhost:8080/hello?param1=value1&param2=value2"
```

```
=== REQUEST SUMMARY ===
Timestamp: 2025-01-15T10:30:45Z
Method: GET | Protocol: HTTP/1.1 | Host: localhost:8080
Full URL: /hello?param1=value1&param2=value2
Remote Address: [::1]:53818
User Agent: curl/8.7.1

=== URL INFORMATION ===
Path: /hello?param1=value1&param2=value2
Query Parameters:
  param1 = value1
  param2 = value2

=== REQUEST HEADERS ===
* Accept         : */*

All Headers:
  User-Agent          : curl/8.7.1
  Accept              : */*

=== REQUEST BODY ===
Content-Length: 0 bytes
Content-Type: 
Body: (empty)

=== FORM DATA ===
Combined Form Data (GET + POST):
  param1 = value1
  param2 = value2

=== SERVER INFORMATION ===
Server Hostname: container-abc123
Go Version: go1.24
Server OS: linux/amd64
Request Start Time: 2025-01-15T10:30:45.123456789Z

=== REQUEST COMPLETED ===
Processing Time: 76.208µs
```

### JSON POST Request
```bash
$ curl -X POST -H "Content-Type: application/json" \
  -d '{"name":"test","value":123,"nested":{"key":"data"}}' \
  http://localhost:8080/api/test
```

```
=== REQUEST SUMMARY ===
Timestamp: 2025-01-15T10:31:22Z
Method: POST | Protocol: HTTP/1.1 | Host: localhost:8080
Full URL: /api/test
Remote Address: [::1]:53825
User Agent: curl/8.7.1

=== URL INFORMATION ===
Path: /api/test
Query Parameters: (none)

=== REQUEST HEADERS ===
* Content-Type   : application/json
* Content-Length : 51
* Accept         : */*

All Headers:
  User-Agent          : curl/8.7.1
  Accept              : */*
  Content-Type        : application/json
  Content-Length      : 51

=== REQUEST BODY ===
Content-Length: 51 bytes
Content-Type: application/json
Body Content:
{
  "name": "test",
  "nested": {
    "key": "data"
  },
  "value": 123
}

=== FORM DATA ===
Form Data: (none)

=== SERVER INFORMATION ===
Server Hostname: container-abc123
Go Version: go1.24
Server OS: linux/amd64
Request Start Time: 2025-01-15T10:31:22.456789Z

=== REQUEST COMPLETED ===
Processing Time: 225.959µs
```

### Form Data POST Request
```bash
$ curl -X POST -d "username=testuser&password=secret&email=test@example.com" \
  http://localhost:8080/login
```

```
=== REQUEST SUMMARY ===
Timestamp: 2025-01-15T10:32:10Z
Method: POST | Protocol: HTTP/1.1 | Host: localhost:8080
Full URL: /login
Remote Address: [::1]:53837
User Agent: curl/8.7.1

=== URL INFORMATION ===
Path: /login
Query Parameters: (none)

=== REQUEST HEADERS ===
* Content-Type   : application/x-www-form-urlencoded
* Content-Length : 56
* Accept         : */*

All Headers:
  Accept              : */*
  Content-Length      : 56
  Content-Type        : application/x-www-form-urlencoded
  User-Agent          : curl/8.7.1

=== REQUEST BODY ===
Content-Length: 56 bytes
Content-Type: application/x-www-form-urlencoded
Body: (empty)

=== FORM DATA ===
Combined Form Data (GET + POST):
  username = testuser
  password = secret
  email = test@example.com

POST Form Data Only:
  username = testuser
  password = secret
  email = test@example.com

=== SERVER INFORMATION ===
Server Hostname: container-abc123
Go Version: go1.24
Server OS: linux/amd64
Request Start Time: 2025-01-15T10:32:10.789123Z

=== REQUEST COMPLETED ===
Processing Time: 38.5µs
```

## 🔍 What Information is Provided

- **Request Summary**: Timestamp, method, protocol, host, client IP, user agent
- **URL Analysis**: Full URL, path, and parsed query parameters
- **Header Analysis**: Important headers highlighted, complete header listing
- **Request Body**: Raw content with intelligent formatting (JSON pretty-print)
- **Form Data**: Parsed form data from both GET and POST requests
- **Server Information**: Hostname, Go version, OS, container/K8s environment details
- **Performance Metrics**: Request processing time and precise timestamps
- **Network Details**: Real client IP detection (proxy-aware)

## 🎯 Use Cases

- **API Development**: Debug HTTP requests during API development
- **Load Balancer Testing**: Verify proxy headers and real client IPs
- **Container Debugging**: Test networking in Docker/Kubernetes environments
- **HTTP Learning**: Understand how HTTP requests work
- **Integration Testing**: Validate request formatting in CI/CD pipelines
- **Webhook Testing**: Inspect incoming webhook payloads
- **Reverse Proxy Testing**: Verify header forwarding and modifications

## 🔧 Configuration

The server runs on port 8080 by default and includes:

- **Security timeouts**: Read/write/idle timeouts configured
- **Graceful shutdown**: Proper cleanup on termination
- **Error handling**: Robust error handling and logging
- **Performance**: Optimized for low latency responses

## 🐳 Container Support

- **Multi-architecture builds**: AMD64, ARM64, ARMv6, ARMv7
- **Scratch base image**: Minimal attack surface
- **Published to GHCR**: `ghcr.io/sgaunet/http-echo`
- **Kubernetes ready**: Includes environment detection
- **Health monitoring**: Simple HTTP endpoint for health checks

## 📊 Technical Details

- **Language**: Go 1.24+
- **Dependencies**: Standard library only
- **Binary size**: ~8MB (statically compiled)
- **Memory usage**: <10MB typical
- **Response time**: <1ms typical

## Project Status

🟨 **Maintenance Mode**: This project is in maintenance mode.

While we are committed to keeping the project's dependencies up-to-date and secure, please note the following:

- New features are unlikely to be added
- Bug fixes will be addressed, but not necessarily promptly
- Security updates will be prioritized

## Issues and Bug Reports

We still encourage you to use our issue tracker for:

- 🐛 Reporting critical bugs
- 🔒 Reporting security vulnerabilities
- 🔍 Asking questions about the project

Please check existing issues before creating a new one to avoid duplicates.

## Contributions

🤝 Limited contributions are still welcome.

While we're not actively developing new features, we appreciate contributions that:

- Fix bugs
- Update dependencies
- Improve documentation
- Enhance performance or security

If you're interested in contributing, please read our [CONTRIBUTING.md](link-to-contributing-file) guide for more information on how to get started.

## Support

As this project is in maintenance mode, support may be limited. We appreciate your understanding and patience.

Thank you for your interest in our project!