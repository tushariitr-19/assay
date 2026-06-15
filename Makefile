.PHONY: build test test-unit test-integration test-all lint clean tidy

# Module directories (root is the core library; the rest are nested modules)
MODULES := . ./adk ./mcp ./cli ./examples/research-agent

# Build the assay CLI binary (lives in the cli module)
build:
	cd cli && go build -o ../assay ./cmd/assay/

# Unit tests — keyless, no network. Only the core has unit tests today.
test-unit:
	go test ./... -v

# Integration tests — spawn real servers/agents (requires built binary + creds).
test-integration: build
	go test ./tests/... -v -tags integration -run Integration

# All tests
test-all: test-unit test-integration

# Alias to match portfolio convention
test: test-unit

# Vet + format check across every module
lint:
	@for m in $(MODULES); do \
		echo "==> $$m"; \
		(cd $$m && go vet ./... && gofmt -l .); \
	done

# Tidy every module's go.mod
tidy:
	@for m in $(MODULES); do \
		echo "==> $$m"; \
		(cd $$m && go mod tidy); \
	done

# Remove build artifacts
clean:
	rm -f assay

# Build and test in one shot
ci: test-unit lint