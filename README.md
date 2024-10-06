# Rate Limiter Service

This project implements a rate-limiting service using Go, Memcached, and `testcontainers-go` for integration testing. It demonstrates how to build a scalable rate limiter with a clean architecture, leveraging Memcached for storing rate limit counters and using Docker containers for testing the service.

## Features

- **Rate Limiting**: Implements a rate limiter with customizable limits for each user and application.
- **Memcached Integration**: Uses Memcached to store request counters with time-based expiration.
- **Clean Architecture**: The project is structured following the principles of clean architecture for better maintainability and scalability.
- **Integration Tests**: Uses `testcontainers-go` to run Memcached in a Docker container during tests.

## Project Structure

```bash
.
├── cmd                  # Main entry point for the application
│   └── rate-limiter
│       └── main.go      # Main Go file for starting the service
├── internal             # Internal packages for business logic
│   ├── app              # Core business logic (rate limiter)
│   ├── cache            # Memcached repository implementation
│   └── http             # HTTP handlers and request routing
├── domain               # Domain models and structures
├── Dockerfile           # Dockerfile for building the application
├── Makefile             # Makefile for running tasks like tests, lint, fmt
├── go.mod               # Go module file
├── go.sum               # Go dependencies
└── .golangci.yml        # Configuration for golangci-lint
```

## Prerequisites

Ensure you have the following installed:

- Go 1.19 or higher
- Docker (for running tests)
- Memcached (for running the application locally)

## Running the Application

1. Clone the repository:

    ```bash
    git clone https://github.com/yourusername/rate-limiter.git
    cd rate-limiter
    ```

2. Build and run the service:

    ```bash
    go build -o rate-limiter ./cmd/rate-limiter
    ./rate-limiter
    ```

3. The service will start and listen for HTTP requests. You can send requests to the `/verify` endpoint to test the rate limiter:

    ```bash
    curl -X POST http://localhost:8080/verify \
         -H "Content-Type: application/json" \
         -H "User-id: 12345" \
         -H "App-id: my_app" \
         -d '{"order_id": "order_001"}'
    ```

## Running the Tests

The project includes integration tests that use Docker to run a Memcached instance in a container. To run the tests:

1. Ensure Docker is running on your machine.
2. Use the following command to run the tests:

    ```bash
    make test
    ```

This will start a Memcached container using `testcontainers-go` and run the tests against the running container.

## Code Formatting and Linting

- To format the code using `go fmt`:

    ```bash
    make fmt
    ```

- To run the linter (`golangci-lint`):

    ```bash
    make lint
    ```

## Docker

You can build and run the service inside a Docker container:

1. Build the Docker image:

    ```bash
    docker build -t go-rate-limiter .
    ```

2. Run the Docker container:

    ```bash
    docker run -p 8080:8080 go-rate-limiter
    ```

## Configuration

Configuration settings (such as the Memcached address, rate limits, and TTL) can be customized in the `cache` package.
