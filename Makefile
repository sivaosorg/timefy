.PHONY: dev run build build-cli test tidy deps-upgrade deps-clean-cache lint race bench cover render-examples icons-validate assets-check tree fmt fmt-check clean

LOG_DIR  := logs
BIN_NAME := timefy
BIN_DIR  := bin

# Detect OS for binary extension
ifeq ($(OS),Windows_NT)
	BIN_EXT := .exe
else
	BIN_EXT :=
endif

BIN_OUT := $(BIN_DIR)/$(BIN_NAME)$(BIN_EXT)

# ==============================================================================
# Development
# Prints CLI help output for quick reference during development.
# ==============================================================================
dev:
	go run ./main/main.go

# ==============================================================================
# Running the main application
# Executes cmd/archflow/main.go, useful for development and quick testing.
run:
	go run ./cmd/timefy

# ==============================================================================
# Building the application
# Compiles cmd/archflow into a binary at bin/archflow (or bin/archflow.exe on Windows).
build:
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_OUT) ./cmd/timefy

# Same as build, explicitly for the CLI binary.
build-cli:
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_OUT) ./cmd/timefy

# ==============================================================================
# Module support and testing
# Runs tests across all packages in the project, showing code coverage.
test:
	go test -cover ./...

# Cleans up the module by removing unused dependencies, then re-vendors.
tidy:
	go mod tidy
	go mod vendor

# ==============================================================================
# Upgrading dependencies
# Updates all direct dependencies to their latest minor/patch versions,
# then re-tidies and re-vendors.
deps-upgrade:
	go get -u -t -d -v ./...
	go mod tidy
	go mod vendor

# Removes all items from the Go module cache.
deps-clean-cache:
	go clean -modcache

# ==============================================================================
# Quality
# Runs go vet across all packages.
lint:
	go vet ./...

# Runs all tests with the race detector enabled.
race:
	go test -race -count=1 ./...

# Runs all benchmarks with memory allocation reporting.
bench:
	go test -bench=. -benchmem ./...

# Generates an HTML coverage report at coverage.html.
cover:
	@mkdir -p $(LOG_DIR)
	go test -coverprofile=./$(LOG_DIR)/coverage.out ./...
	go tool cover -html=./$(LOG_DIR)/coverage.out -o ./$(LOG_DIR)/coverage.html

# ==============================================================================
# Examples and icon management
# Runs the basic example to verify rendering works.
render-examples:
	go run ./examples/basic

# Validates icon files in the assets directory.
icons-validate:
	./$(BIN_OUT) icons validate --dir assets

# Runs the asset validation tool to audit all icon files.
assets-check:
	go run ./cmd/assetcheck --dir assets

# ==============================================================================
# Formatting
# Formats all Go source files in-place using gofmt.
fmt:
	go fmt ./...

# Checks whether all Go source files are properly formatted.
# Exits with a non-zero status if any files need formatting.
fmt-check:
	@unformatted=$$(gofmt -l .); \
	if [ -n "$$unformatted" ]; then \
		echo "The following files are not formatted:"; \
		echo "$$unformatted"; \
		exit 1; \
	fi

# ==============================================================================
# Utilities
# Creates a text file representing the project's directory structure.
tree:
	@mkdir -p $(LOG_DIR)
	tree -I ".gradle|.idea|build|logs|.vscode|.git|.github|vendor" > ./$(LOG_DIR)/tree_source_oss.txt
	cat ./$(LOG_DIR)/tree_source_oss.txt

# ==============================================================================
# Clean
# Removes the build directory and log directory.
clean:
	rm -rf $(BIN_DIR)
	rm -rf $(LOG_DIR)