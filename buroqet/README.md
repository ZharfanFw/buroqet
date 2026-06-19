# 🚀 Buroqet — Logistics Management System

Sistem manajemen logistik berbasis microservices, dibangun dengan Go (backend) dan React (frontend).

## Arsitektur

```
buroqet/
├── frontend/               # React + Vite + TypeScript
├── auth-service/           # Go — JWT Authentication (PostgreSQL)
├── tracking-service/       # Go — Package Tracking (MongoDB + Redis + Kafka)
├── order-management-service/ # Go — Order CRUD (PostgreSQL + Kafka)
├── dispatch-fleet-service/ # Go — Fleet & Courier Management (PostgreSQL)
├── pricing-service/        # Go — Pricing Engine (PostgreSQL)
├── settlement-service/     # Go — Payment Settlement (PostgreSQL)
├── epod-service/           # Go — Electronic Proof of Delivery (PostgreSQL)
├── warehouse-service/      # Go — Warehouse Inventory (PostgreSQL)
├── infra/kubernetes/       # Shared K8s manifests
├── docker-compose.yml      # Local dev stack
└── Makefile                # Unified commands
```

## Tech Stack

| Layer | Technology |
|---|---|
| Frontend | React 18, Vite, TypeScript, React Router v6, Zustand, Axios |
| Backend | Go 1.22+, Gin, GORM |
| Databases | PostgreSQL 16, MongoDB 7, Redis 7 |
| Messaging | Apache Kafka |
| Container | Docker, Kubernetes |
| CI/CD | Jenkins |

## Quick Start

### Prasyarat
- Docker & Docker Compose
- Node.js 20+
- Go 1.22+

### 1. Clone & Setup

```bash
git clone <repo-url>
cd buroqet

# Copy env files
cp frontend/.env.example frontend/.env
```

### 2. Jalankan Lokal

```bash
# Hanya infra (postgres, mongo, redis, kafka)
make dev-infra

# Frontend dev server (di terminal lain)
make frontend

# Atau full stack dengan Docker
make dev
```

Frontend tersedia di: **http://localhost:3000**

### 3. Akses per Service

| Service | Port | Endpoint |
|---|---|---|
| Frontend | 3000 | http://localhost:3000 |
| Auth | 8080 | http://localhost:8080 |
| Tracking | 8081 | http://localhost:8081 |
| Order | 8082 | http://localhost:8082 |
| Dispatch | 8083 | http://localhost:8083 |
| Pricing | 8084 | http://localhost:8084 |
| Settlement | 8085 | http://localhost:8085 |
| ePOD | 8086 | http://localhost:8086 |
| Warehouse | 8087 | http://localhost:8087 |

## Development

```bash
# Run tests semua service
make test

# Run test per service
make test-auth
make test-tracking

# Build semua Docker images
make build

# Stop semua container
make clean
```

## Deploy ke Kubernetes

```bash
# Apply namespace & ingress
kubectl apply -f infra/kubernetes/

# Deploy per service
kubectl apply -f auth-service/deployments/kubernetes/
kubectl apply -f tracking-service/deployments/kubernetes/
# ... dst

# Deploy frontend
kubectl apply -f frontend/deployments/kubernetes/
```

## Kontribusi

1. Buat branch baru dari `main`: `git checkout -b feature/nama-fitur`
2. Kerjakan di service yang kamu handle
3. Pastikan tests lolos: `make test-<service>`
4. Push & buat Pull Request

## Tim

> Buroqet — Cloud Computing Project, Semester 4
