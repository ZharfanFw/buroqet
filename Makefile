# ================================================================
# Buroqet — Root Makefile
# ================================================================

.PHONY: help dev dev-infra build test clean frontend

# Default: show help
help:
	@echo ""
	@echo "  🚀 Buroqet — Makefile Commands"
	@echo ""
	@echo "  Development:"
	@echo "    make dev-infra    Start infra only (postgres, mongo, redis, kafka)"
	@echo "    make frontend     Start React frontend dev server"
	@echo "    make dev          Start full stack (docker compose)"
	@echo ""
	@echo "  Build:"
	@echo "    make build        Build semua Docker images"
	@echo "    make build-fe     Build frontend production"
	@echo ""
	@echo "  Testing:"
	@echo "    make test         Run tests semua service"
	@echo "    make test-auth    Run tests auth-service"
	@echo "    make test-tracking  Run tests tracking-service"
	@echo ""
	@echo "  Cleanup:"
	@echo "    make clean        Stop dan hapus semua container"
	@echo ""

# ─── Infrastructure Only ─────────────────────────────────────
dev-infra:
	docker compose up -d postgres mongodb redis zookeeper kafka
	@echo "✅ Infrastructure running:"
	@echo "   Postgres  → localhost:5432"
	@echo "   MongoDB   → localhost:27017"
	@echo "   Redis     → localhost:6379"
	@echo "   Kafka     → localhost:9092"

# ─── Frontend Dev Server ─────────────────────────────────────
frontend:
	cd frontend && npm run dev

# ─── Full Stack ──────────────────────────────────────────────
dev:
	docker compose up -d
	@echo "✅ Full stack running:"
	@echo "   Frontend  → http://localhost:3000"
	@echo "   Auth      → http://localhost:8080"
	@echo "   Tracking  → http://localhost:8081"
	@echo "   Order     → http://localhost:8082"
	@echo "   Dispatch  → http://localhost:8083"
	@echo "   Pricing   → http://localhost:8084"
	@echo "   Settlement → http://localhost:8085"
	@echo "   ePOD      → http://localhost:8086"
	@echo "   Warehouse → http://localhost:8087"

# ─── Build ───────────────────────────────────────────────────
build:
	docker compose build

build-fe:
	cd frontend && npm run build

# ─── Testing ─────────────────────────────────────────────────
test:
	$(MAKE) test-auth
	$(MAKE) test-tracking
	$(MAKE) test-order
	$(MAKE) test-dispatch

test-auth:
	cd auth-service && go test ./...

test-tracking:
	cd tracking-service && go test ./...

test-order:
	cd order-management-service && go test ./...

test-dispatch:
	cd dispatch-fleet-service && go test ./...

test-pricing:
	cd pricing-service && go test ./...

test-settlement:
	cd settlement-service && go test ./...

test-epod:
	cd epod-service && go test ./...

test-warehouse:
	cd warehouse-service && go test ./...

# ─── Cleanup ─────────────────────────────────────────────────
clean:
	docker compose down -v
	@echo "✅ All containers and volumes removed"

logs:
	docker compose logs -f --tail=50
