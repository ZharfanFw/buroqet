# Makefile — Dispatch & Fleet Microservice
# ─────────────────────────────────────────────────────────────────────────────
# Usage:
#   make help           → tampilkan semua target
#   make ci             → jalankan full pipeline (unit → lint → build → functional)
#   make test-unit      → unit test saja
#   make test-func      → functional test saja (membutuhkan Docker)
# ─────────────────────────────────────────────────────────────────────────────

.PHONY: help ci generate mock test-unit lint build test-func db-up db-down clean

# ─── Variabel ────────────────────────────────────────────────────────────────
BINARY_NAME     := dispatch-fleet
BUILD_DIR       := ./bin
GO_PACKAGES     := ./internal/...
FUNCTIONAL_DIR  := ./tests/functional/...
TEST_DB_URL     := "host=localhost port=5433 user=test_user password=test_pass dbname=dispatch_test sslmode=disable"

# ─── Default target ───────────────────────────────────────────────────────────
help: ## Tampilkan daftar semua target
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

# ─── Code Generation ─────────────────────────────────────────────────────────
generate: ## Generate semua mock dengan go generate
	@echo "🔧 Generating mocks..."
	go generate ./internal/domain/...
	@echo "✅ Mocks generated."

# ─── STEP 1: Unit Test ────────────────────────────────────────────────────────
test-unit: ## Jalankan unit test (tanpa DB)
	@echo "🧪 Running unit tests..."
	go test \
		-v \
		-race \
		-count=1 \
		-timeout=60s \
		-coverprofile=coverage.unit.out \
		$(GO_PACKAGES)
	@echo "✅ Unit tests done."
	@go tool cover -func=coverage.unit.out | tail -1

# ─── STEP 2: Lint ─────────────────────────────────────────────────────────────
lint: ## Jalankan golangci-lint
	@echo "🔍 Running linter..."
	@which golangci-lint > /dev/null 2>&1 || \
		(echo "golangci-lint not found, installing..." && \
		 go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run ./...
	@echo "✅ Lint passed."

# ─── STEP 3: Build ────────────────────────────────────────────────────────────
build: ## Build binary
	@echo "🏗️  Building binary..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/...
	@echo "✅ Build done → $(BUILD_DIR)/$(BINARY_NAME)"

# ─── STEP 4: Functional Test ──────────────────────────────────────────────────
db-up: ## Jalankan test database container
	@echo "🐳 Starting test database..."
	docker compose -f docker-compose.test.yml up -d --wait
	@echo "✅ Test DB ready on port 5433."

db-down: ## Hentikan dan hapus test database container
	@echo "🛑 Stopping test database..."
	docker compose -f docker-compose.test.yml down -v
	@echo "✅ Test DB stopped."

test-func: db-up ## Jalankan functional test (otomatis start/stop DB)
	@echo "🧪 Running functional tests..."
	TEST_DATABASE_URL=$(TEST_DB_URL) \
	go test \
		-v \
		-count=1 \
		-timeout=120s \
		-tags=functional \
		$(FUNCTIONAL_DIR) ; \
	EXIT_CODE=$$? ; \
	$(MAKE) db-down ; \
	exit $$EXIT_CODE

# ─── Full CI Pipeline ─────────────────────────────────────────────────────────
ci: test-unit lint build test-func ## 🚀 Jalankan full pipeline: unit → lint → build → functional
	@echo ""
	@echo "🎉 All pipeline stages passed!"

# ─── Utilities ────────────────────────────────────────────────────────────────
clean: ## Hapus build artifacts dan coverage files
	@rm -rf $(BUILD_DIR) coverage.*.out
	@echo "🧹 Cleaned."

coverage-html: ## Buka laporan coverage unit test di browser
	go tool cover -html=coverage.unit.out