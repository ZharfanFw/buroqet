.PHONY: ci test-unit lint build test-func db-up db-down

BINARY_NAME     := auth-service
TEST_DB_URL     := "host=localhost port=5433 user=postgres password=password dbname=auth_db_test sslmode=disable"

# 1. Unit Test
test-unit:
	@echo "🧪 Running unit tests..."
	go test -v -short ./internal/service/...

# 2. Lint/Vet
lint:
	@echo "🔍 Running go vet..."
	go vet ./...

# 3. Build Image Local
build:
	@echo "🏗️  Building docker image..."
	docker build -t imamaziz/auth-service:v1 .

# 4. Angkat & Turunkan DB Test untuk Functional Test
db-up:
	@echo "🐳 Starting test database..."
	docker compose -f docker-compose.test.yml up -d

db-down:
	@echo "🛑 Stopping test database..."
	docker compose -f docker-compose.test.yml down -v

# 5. Functional Test
test-func: db-up
	@echo "🧪 Running functional tests..."
	@sleep 3 # Tunggu DB benar-benar ready
	TEST_DATABASE_URL=$(TEST_DB_URL) go test -v ./tests/functional/...; \
	EXIT_CODE=$$?; \
	$(MAKE) db-down; \
	exit $$EXIT_CODE

# Target otomatis untuk Jenkins Pipeline (Langkah 2 sampai 5)
ci: test-unit lint build test-func