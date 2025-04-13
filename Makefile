.PHONY: build install clean test release

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS = -ldflags "-s -w -X github.com/VesperAkshay/lazynode/pkg/version.Version=$(VERSION)"
BINARY_NAME = lazynode

build:
	go build $(LDFLAGS) -o $(BINARY_NAME) ./cmd/lazynode

install: build
	mv $(BINARY_NAME) $(GOPATH)/bin/

clean:
	rm -f $(BINARY_NAME)
	rm -rf dist/

test:
	go test -v ./...

lint:
	golangci-lint run

# Create a new release (requires goreleaser and git tag)
release:
	goreleaser release --clean

# Create a snapshot release for testing
snapshot:
	goreleaser release --snapshot --clean

# Tag a new version - usage: make tag-release VERSION=1.0.0
tag-release:
	git tag -a v$(VERSION) -m "Release v$(VERSION)"
	git push origin v$(VERSION)

# Cross compile for major platforms
cross-build:
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)_linux_amd64 ./cmd/lazynode
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BINARY_NAME)_linux_arm64 ./cmd/lazynode
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)_darwin_amd64 ./cmd/lazynode
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BINARY_NAME)_darwin_arm64 ./cmd/lazynode
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)_windows_amd64.exe ./cmd/lazynode 