# Build variables
BINARY_NAME := sekret
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -X github.com/yquansah/sekret/cmd.version=$(VERSION) -X github.com/yquansah/sekret/cmd.commit=$(COMMIT) -X github.com/yquansah/sekret/cmd.date=$(DATE)

# Default target
.PHONY: all
all: build

# Build the binary
.PHONY: build
build:
	go build -ldflags "$(LDFLAGS)" -o $(BINARY_NAME) .

# Install the binary
.PHONY: install
install:
	go install -ldflags "$(LDFLAGS)" .

# Clean build artifacts
.PHONY: clean
clean:
	rm -f $(BINARY_NAME)

# Run tests
.PHONY: test
test:
	go test ./...

# Show version info that would be embedded
.PHONY: version-info
version-info:
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT)"
	@echo "Date: $(DATE)"

# Build for multiple platforms
.PHONY: build-all
build-all:
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(BINARY_NAME)-linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(BINARY_NAME)-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(BINARY_NAME)-windows-amd64.exe .

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build        Build the binary with version information"
	@echo "  install      Install the binary with version information"
	@echo "  clean        Remove build artifacts"
	@echo "  test         Run tests"
	@echo "  version-info Show version information that would be embedded"
	@echo "  build-all    Build for multiple platforms"
	@echo "  help         Show this help message"