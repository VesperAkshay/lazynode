.PHONY: build clean run test release install

# Default target
all: build

# Build for current platform
build:
	go build -o bin/lazynode cmd/lazynode/main.go

# Clean build artifacts
clean:
	rm -rf bin dist

# Run the application
run:
	go run cmd/lazynode/main.go

# Run tests
test:
	go test ./...

# Create a new release
release:
	./scripts/build.sh

# Install locally
install: build
	cp bin/lazynode /usr/local/bin/lazynode

# Cross compile for all platforms
crossbuild:
	./scripts/build.sh

# Generate docs
docs:
	mkdir -p docs
	# Add doc generation commands here

# Help command
help:
	@echo "LazyNode Make Targets:"
	@echo "  build      - Build for current platform"
	@echo "  clean      - Remove build artifacts"
	@echo "  run        - Run the application"
	@echo "  test       - Run tests"
	@echo "  release    - Create release builds for all platforms"
	@echo "  install    - Install locally (requires sudo)"
	@echo "  crossbuild - Cross compile for all platforms"
	@echo "  docs       - Generate documentation"
	@echo "  help       - Show this help" 