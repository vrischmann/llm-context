# Justfile for llm-context project

# Build the project
build:
	go build -o llm-context

# Run the project
run:
	./llm-context

# Test goreleaser snapshot build
test-release:
	goreleaser release --snapshot --clean

# Run CI checks
ci:
	go mod download
	go mod verify
	if [ "$(gofmt -s -l . | wc -l)" -gt 0 ]; then echo "Code is not properly formatted:"; gofmt -s -d .; exit 1; fi
	go vet ./...
	go build -v ./...
	go test -v -timeout=10s ./...

# Clean build artifacts
clean:
	rm -f llm-context
	rm -rf dist

# Show version
version:
	./llm-context --version
