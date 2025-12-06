.PHONY: all build run test proto clean harness execute

# Build settings
BINARY_NAME=server
CMD_PATH=./cmd/server

# Build the application
build:
	go build -o $(BINARY_NAME) $(CMD_PATH)

# Run the application
run: build
	./$(BINARY_NAME)

# Run tests
test:
	go test ./...

# Run tests with verbose output
test-v:
	go test -v ./...

# Generate protobuf code
proto:
	buf generate

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME)
	rm -f *.exe

# Download dependencies
deps:
	go mod tidy

# Run go vet
vet:
	go vet ./...

# Run harness script
harness:
	py scripts/execute_harness.py

# Show status only
harness-status:
	py scripts/execute_harness.py --status

# Dry run
harness-dry:
	py scripts/execute_harness.py --dry-run

# Build and test
check: vet test build

# Full rebuild
rebuild: clean deps build
