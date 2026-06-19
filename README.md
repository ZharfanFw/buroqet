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
├── settlement-service/     # Go — Payment Settlement (PostgreSQL + Kafka)
├── epod-service/           # Go — Electronic Proof of Delivery (Kafka)
├── warehouse-service/      # Go — Warehouse Management (PostgreSQL + Kafka)
├── infra/kubernetes/       # Shared K8s manifests
├── docker-compose.yml      # Local dev stack
└── Makefile                # Unified commands
```

## Tech Stack

| Layer | Technology |
|---|---|
| Frontend | React 18, Vite, TypeScript, React Router v6, Zustand, Axios |
| Backend | Go 1.22+, Gin / net/http, GORM |
| Databases | PostgreSQL 16, MongoDB 7, Redis 7 |
| Messaging | Apache Kafka (Confluent 7.6) |
| Container | Docker, Docker Compose |
| Orchestration | Kubernetes (Minikube) |
| CI/CD | Jenkins |

---

## ⚡ Quick Start (Full Docker)

> Cara paling mudah — semua service berjalan lewat Docker Compose.

### Prasyarat

- [Docker Desktop](https://www.docker.com/products/docker-desktop/) (sudah include Docker Compose)
- Git

### Langkah

```bash
# 1. Clone repository
git clone <repo-url>
cd buroqet

# 2. Buat database yang dibutuhkan (wms_db & settlement_db)
#    Jalankan infra dulu
docker compose up -d postgres
#    Tunggu ~10 detik, lalu buat database
docker exec -it buroqet-postgres psql -U buroqet -c "CREATE DATABASE wms_db;"
docker exec -it buroqet-postgres psql -U buroqet -c "CREATE DATABASE settlement_db;"

# 3. Jalankan full stack
docker compose up -d

# 4. Cek semua container running
docker compose ps
```

Akses frontend di: **http://localhost:3000**

---

## 🛠️ Menjalankan Lokal (Hybrid Mode)

Mode ini menjalankan **infrastruktur via Docker** (PostgreSQL, MongoDB, Redis, Kafka) dan **service Go secara langsung** di terminal. Cocok untuk development dan debugging.

### Prasyarat

- Docker Desktop
- Go 1.22+
- Node.js 20+ (untuk frontend)

### Step 1 — Jalankan Infrastruktur

```bash
# Jalankan hanya infra (DB, Redis, Kafka)
docker compose up -d postgres mongodb redis zookeeper kafka

# Tunggu ~15 detik sampai Kafka ready, lalu buat database
docker exec -it buroqet-postgres psql -U buroqet -c "CREATE DATABASE wms_db;"
docker exec -it buroqet-postgres psql -U buroqet -c "CREATE DATABASE settlement_db;"
```

### Step 2 — Jalankan Service (PILIH SALAH SATU CARA)

> [!WARNING]
> **PENTING:** Ada DUA cara untuk menjalankan sistem ini.
> **CARA 1 (Disarankan):** Jalankan semuanya dengan Docker Compose. Kamu tidak perlu menjalankan `go run` apa pun. Cukup ketik `docker compose up -d` lalu langsung lompat ke Step 3 (Frontend).
> **CARA 2 (Untuk Debugging):** Jika kamu ingin me-run servicenya manual di terminal, pastikan kamu **hanya menjalankan database saja** di Docker (`docker compose up -d postgres mongodb redis zookeeper kafka`), lalu kamu bisa menjalankan `go run` di bawah ini satu per satu.
> 
> 🛑 **SANGAT PENTING JIKA MENDAPAT ERROR "Only one usage of each socket address" (Port Conflict):**
> Jika kamu sebelumnya sudah terlanjur menjalankan `docker compose up -d` (yang menyalakan semua service backend), port 8080, 8081, dsb **sudah dipakai oleh Docker**. Kamu WAJIB mematikan containernya terlebih dahulu sebelum menjalankan `go run` di terminal:
> ```bash
> docker compose stop auth-service tracking-service order-management-service dispatch-fleet-service pricing-service settlement-service epod-service warehouse-service
> ```

---
**Instruksi di bawah ini HANYA untuk CARA 2 (Manual Run). Jika kamu menggunakan CARA 1, LEWATI BAGIAN INI!**

### Pilihan A — Koneksi ke Local Infra (Docker)
Gunakan environment variables berikut jika kamu ingin connect ke `postgres`, `mongodb`, dan `redis` yang berjalan di Docker lokalmu.


#### Auth Service (port 8080)
```bash
cd auth-service
$env:DB_HOST="localhost"
$env:DB_USER="buroqet"
$env:DB_PASSWORD="buroqet123"
$env:DB_NAME="buroqet"
$env:JWT_SECRET="dev-secret-key"
go run ./cmd/main.go
```

#### Tracking Service (port 8081)
```bash
$env:MONGO_URI="mongodb+srv://zharfanfaz21_db_user:bsUIL2Yyewp5hfUj@buroqet-tracking.qdsmduw.mongodb.net/?appName=buroqet-tracking"
$env:MONGO_DB="tracking_db"
$env:REDIS_ADDR="rediss://default:gQAAAAAAAfchAAIgcDEwOGNiOWUwNTA4NmU0MDJkOWI3YjAxYzk4NTE3ZjA0MA@distinct-dinosaur-128801.upstash.io:6379"
$env:KAFKA_BROKER="localhost:9092"
$env:APP_PORT="8081"
go run ./cmd/main.go
```

#### Order Management Service (port 8082)
```bash
cd order-management-service
$env:DB_HOST="localhost"
$env:DB_USER="buroqet"
$env:DB_PASSWORD="buroqet123"
$env:DB_NAME="buroqet"
$env:KAFKA_BROKERS="localhost:9092"
$env:PRICING_SERVICE_URL="http://localhost:8084"
go run ./cmd/main.go
```

#### Dispatch Fleet Service (port 8083)
```bash
cd dispatch-fleet-service
$env:DB_HOST="localhost"
$env:DB_USER="buroqet"
$env:DB_PASSWORD="buroqet123"
$env:DB_NAME="buroqet"
go run ./cmd/main.go
```

#### Pricing Service (port 8084)
```bash
cd pricing-service
$env:DB_HOST="localhost"
$env:DB_USER="buroqet"
$env:DB_PASSWORD="buroqet123"
$env:DB_NAME="buroqet"
go run ./cmd/main.go
```

#### Settlement Service (port 8085)
```bash
cd settlement-service
$env:DB_HOST="localhost"
$env:DB_USER="buroqet"
$env:DB_PASSWORD="buroqet123"
$env:DB_NAME="settlement_db"
$env:APP_PORT="8085"
$env:KAFKA_BROKER="localhost:9092"
$env:PRICING_SERVICE_URL="http://localhost:8084"
go run ./cmd/main.go
```

#### ePOD Service (port 8086)
```bash
cd epod-service
$env:KAFKA_BROKER="localhost:9092"
go run ./cmd/main.go
```

#### Warehouse Service (port 8087)
```bash
cd warehouse-service
$env:DB_HOST="localhost"
$env:DB_USER="buroqet"
$env:DB_PASSWORD="buroqet123"
$env:DB_NAME="wms_db"
$env:APP_PORT="8087"
$env:KAFKA_BROKER="localhost:9092"
go run ./cmd/main.go
```

### Pilihan B — Koneksi ke Cloud Database (Supabase & MongoDB Atlas)
Jika kamu ingin menjalankan `go run` di terminal tapi connect ke **Database Cloud (Supabase untuk Postgres & Atlas untuk MongoDB)**, gunakan setup `env` berikut.
*(Catatan: Redis sebaiknya tetap menggunakan `localhost:6379` karena Upstash yang tersedia adalah versi REST API, sedangkan kode Go menggunakan koneksi Redis TCP murni).*

#### Auth Service (port 8080)
```powershell
cd auth-service
$env:DB_HOST="aws-1-ap-southeast-1.pooler.supabase.com"
$env:DB_PORT="5432"
$env:DB_USER="postgres.gvmknmaybejslbghzsov"
$env:DB_PASSWORD="HfuyrF1xd1iUVagY"
$env:DB_NAME="postgres"
$env:JWT_SECRET="dev-secret-key"
go run ./cmd/main.go
```

#### Tracking Service (port 8081)
```powershell
cd tracking-service
$env:MONGO_URI="mongodb+srv://zharfanfaz21_db_user:bsUIL2Yyewp5hfUj@buroqet-tracking.qdsmduw.mongodb.net/?appName=buroqet-tracking"
$env:MONGO_DB="tracking_db"
$env:REDIS_ADDR="localhost:6379"
$env:REDIS_PASSWORD="buroqet123"
$env:KAFKA_BROKER="localhost:9092"
$env:APP_PORT="8081"
go run ./cmd/main.go
```

#### Settlement Service (port 8085)
```powershell
cd settlement-service
$env:DB_HOST="aws-1-ap-southeast-1.pooler.supabase.com"
$env:DB_PORT="5432"
$env:DB_USER="postgres.gvmknmaybejslbghzsov"
$env:DB_PASSWORD="HfuyrF1xd1iUVagY"
$env:DB_NAME="postgres"
$env:APP_PORT="8085"
$env:KAFKA_BROKER="localhost:9092"
$env:PRICING_SERVICE_URL="http://localhost:8084"
go run ./cmd/main.go
```

#### Warehouse Service (port 8087)
```powershell
cd warehouse-service
$env:DB_HOST="aws-1-ap-southeast-1.pooler.supabase.com"
$env:DB_PORT="5432"
$env:DB_USER="postgres.gvmknmaybejslbghzsov"
$env:DB_PASSWORD="HfuyrF1xd1iUVagY"
$env:DB_NAME="postgres"
$env:APP_PORT="8087"
$env:KAFKA_BROKER="localhost:9092"
go run ./cmd/main.go
```

### Step 3 — Jalankan Frontend

```bash
cd frontend
npm install
npm run dev
# Akses di http://localhost:5173
```

---

## 🌐 Port Reference

| Service | Host Port | URL |
|---|---|---|
| Frontend (Docker) | 3000 | http://localhost:3000 |
| Frontend (Dev) | 5173 | http://localhost:5173 |
| Auth | 8080 | http://localhost:8080/health |
| Tracking | 8081 | http://localhost:8081/health |
| Order | 8082 | http://localhost:8082/health |
| Dispatch | 8083 | http://localhost:8083 |
| Pricing | 8084 | http://localhost:8084/health |
| Settlement | 8085 | http://localhost:8085/health |
| ePOD | 8086 | http://localhost:8086 |
| Warehouse | 8087 | http://localhost:8087/health |
| PostgreSQL | 5432 | — |
| MongoDB | 27017 | — |
| Redis | 6379 | — |
| Kafka | 9092 | — |

---

## 🧪 Testing Endpoint Penting

### Cek Status Semua Service
```powershell
# Health check (PowerShell)
Invoke-RestMethod -Uri "http://localhost:8080/health"   # Auth
Invoke-RestMethod -Uri "http://localhost:8081/health"   # Tracking
Invoke-RestMethod -Uri "http://localhost:8085/health"   # Settlement
Invoke-RestMethod -Uri "http://localhost:8087/health"   # Warehouse
```

### Warehouse — Scan Paket Masuk Gudang
```powershell
Invoke-RestMethod -Uri "http://localhost:8087/api/v1/inbound" -Method Post -ContentType "application/json" -Body '{"awb": "BQ-2024-JKT-001", "hub_id": "HUB-JKT-01"}'
# Expected: 201 Created
```

### Settlement — Catat Komisi Kurir
```powershell
Invoke-RestMethod -Uri "http://localhost:8085/api/v1/commissions" -Method Post -ContentType "application/json" -Body '{"courier_id": "COURIER-001", "awb": "BQ-2024-JKT-001", "service_type": "REGULER"}'
# Expected: 201 Created
```

### Settlement — Lihat Penghasilan Kurir
```powershell
Invoke-RestMethod -Uri "http://localhost:8085/api/v1/couriers/COURIER-001/earnings"
# Expected: 200 dengan ringkasan total, pending, paid
```

### Auth — Register & Login
```powershell
# Register
Invoke-RestMethod -Uri "http://localhost:8080/auth/register" -Method Post -ContentType "application/json" -Body '{"name": "Test User", "email": "test@buroqet.id", "password": "secret123", "role": "pelanggan"}'

# Login
Invoke-RestMethod -Uri "http://localhost:8080/auth/login" -Method Post -ContentType "application/json" -Body '{"email": "test@buroqet.id", "password": "secret123"}'
```

---

## 🪵 Melihat Log

```bash
# Log semua service (Docker)
docker compose logs -f

# Log service tertentu
docker compose logs -f warehouse-service
docker compose logs -f settlement-service
docker compose logs -f tracking-service

# Log singkat (50 baris terakhir)
docker compose logs --tail=50 settlement-service
```

---

## 🔧 Development

```bash
# Run tests semua service
make test

# Run test per service
make test-auth
make test-tracking
make test-warehouse
make test-settlement

# Build semua Docker images
make build

# Stop semua container
make clean

# Rebuild dan restart service tertentu
docker compose up -d --build warehouse-service
docker compose up -d --build settlement-service
```

---

## ☸️ Deploy ke Kubernetes

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

---

## Kontribusi

1. Buat branch baru dari `main`: `git checkout -b feature/nama-fitur`
2. Kerjakan di service yang kamu handle
3. Pastikan tests lolos: `make test-<service>`
4. Push & buat Pull Request

## Tim

> Buroqet — Cloud Computing Project, Semester 4
