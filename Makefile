
MIN_GO_VERSION := 1.24


GO := go
GOBUILD := $(GO) build
GOTEST := $(GO) test
GOMOD := $(GO) mod


# Phony targets are rules that are not actual files.
.PHONY: all build clean deps test tidy check-go-version

# === Main Targets ===

# The default target, which builds the extension.
all: build-all 


build-all: container_processes

# Build the Go application. It depends on the version check and dependencies being present.
container_processes: extensions/container_processes/*.go 
	$(GOBUILD) -o container_processes.ext ./extensions/container_processes/*.go

# Run all tests verbosely.
test:
	$(GOTEST) -v ./...

# Remove the built binary and clean the test cache.
clean:
	$(GO) clean -testcache

# Download and verify module dependencies. Depends on `tidy` to ensure go.mod is clean first.
deps: tidy
	$(GOMOD) download
	$(GOMOD) verify

# Tidy up the go.mod and go.sum files.
tidy:
	$(GOMOD) tidy

# === Version Check ===
check-go-version:
	@if ! $(GO) version | awk '{split($$3, v, "go"); if (system("printf \"%s\\n%s\" \""$(MIN_GO_VERSION)"\" \"" v[2] "\" | sort -VC;")) exit 1}'; then \
		echo "Error: Go version is too old. Found: `$(GO) version`. Required: >= $(MIN_GO_VERSION)"; \
		exit 1; \
	fi
	@echo "Go version check passed."