.PHONY: lint

lint:
	@echo "Running golangci-lint..."
	@golangci-lint run --config .golangci.yml

.PHONY: test

test:
	@echo "Running tests..."
	@go test ./... -v

.PHONY: fmt

fmt:
	@echo "Running go fmt..."
	@go fmt ./...