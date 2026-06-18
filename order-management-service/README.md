# 📦 Order Management Service (OMS)

Microservice untuk pembuatan dan manajemen order pengiriman dalam sistem logistik **Buroqet**.

---

## 🔧 Tech Stack

| Layer | Teknologi |
|---|---|
| Language | Go 1.23 |
| HTTP Framework | Gin |
| ORM | GORM |
| Database | PostgreSQL (Supabase) — schema `oms`, tabel `orders` |
| Message Broker | Apache Kafka (`order.created`) |
| External Service | Pricing Service (REST HTTP) |
| Testing | Testify + gomock |

---

## 📁 Struktur Direktori

```
order-management-service/
├── cmd/
│   └── main.go                  # Entry point aplikasi
├── internal/
│   ├── domain/
│   │   └── order.go             # Entity Order, enum, DTO Request/Response
│   ├── handler/
│   │   ├── order_handler.go     # HTTP handler (Gin)
│   │   └── order_handler_test.go
│   ├── service/
│   │   ├── order_service.go     # Business logic
│   │   ├── order_service_test.go
│   │   └── pricing_client.go    # HTTP client ke Pricing Service
│   ├── repository/
│   │   └── order_repository.go  # GORM queries ke PostgreSQL
│   └── kafka/
│       └── producer.go          # Kafka producer (order.created)
├── mocks/                       # Generated mocks (gomock)
├── tests/
│   └── functional/
│       └── order_functional_test.go  # Functional/integration tests
├── deployments/
│   ├── Dockerfile
│   └── kubernetes/
│       ├── configmap.yaml
│       └── deployment.yaml
├── go.mod
└── go.sum
```

---

## 🌐 API Endpoints

| Method | Path | Deskripsi |
|---|---|---|
| `POST` | `/api/v1/orders` | Buat order baru |
| `GET` | `/api/v1/orders/:awb` | Ambil detail order by AWB |
| `GET` | `/health` | Health check |

### `POST /api/v1/orders` — Buat Order Baru

**Request Body:**
```json
{
  "sender_name":      "Rangga",
  "sender_phone":     "0812345678",
  "sender_address":   "Jl. Sukarno Hatta No. 10",
  "origin_postal":    "40286",
  "origin_city":      "Bandung",
  "receiver_name":    "Zharfan",
  "receiver_phone":   "0898765432",
  "receiver_address": "Jl. Merdeka No. 42",
  "dest_postal":      "10110",
  "dest_city":        "Jakarta",
  "weight_actual":    1.5,
  "length":           10.0,
  "width":            10.0,
  "height":           10.0,
  "service_type":     "REG",
  "payment_type":     "TRANSFER"
}
```

**Nilai enum yang valid:**

| Field | Nilai yang diterima |
|---|---|
| `service_type` | `REG`, `EXP`, `CARGO` |
| `payment_type` | `COD`, `TRANSFER`, `EWALLET`, `VA` |

**Response `201 Created`:**
```json
{
  "success": true,
  "data": {
    "order_id":    "uuid",
    "awb_number":  "BQ-2026-XKT-042",
    "payment_ref": "uuid",
    "status":      "ORDER_CREATED",
    "total_cost":  17000,
    "payment_url": "https://pay.example.com/invoice/..."
  }
}
```

> **Catatan:** `payment_url` hanya ada untuk payment_type selain `COD`.

---

## ⚙️ Konfigurasi Environment Variables

| Variable | Deskripsi | Default |
|---|---|---|
| `DATABASE_URL` | PostgreSQL DSN (Supabase) | `host=localhost user=postgres ...` |
| `PORT` | Port HTTP server | `8082` |
| `KAFKA_BROKER` | Kafka broker address | `localhost:9092` |
| `PRICING_SERVICE_URL` | Base URL Pricing Service | *(stub digunakan jika kosong)* |

---

## 🚀 Menjalankan Lokal

### 1. Prasyarat
- Go 1.23+
- Docker (untuk Kafka)
- Akses ke Supabase project Buroqet

### 2. Jalankan Infrastruktur Kafka
```bash
# Dari root folder buroqet/
docker compose up -d zookeeper kafka
```

### 3. Set Environment Variables & Jalankan Server

**CMD (Windows):**
```cmd
set "DATABASE_URL=postgresql://postgres.PROJECTREF:PASSWORD@aws-1-ap-southeast-1.pooler.supabase.com:5432/postgres?sslmode=require&search_path=oms"
set PORT=8082
set KAFKA_BROKER=localhost:9092
go run ./cmd/main.go
```

**PowerShell:**
```powershell
$env:DATABASE_URL="postgresql://postgres.PROJECTREF:PASSWORD@aws-1-ap-southeast-1.pooler.supabase.com:5432/postgres?sslmode=require&search_path=oms"
$env:PORT="8082"
$env:KAFKA_BROKER="localhost:9092"
go run ./cmd/main.go
```

> **Penting:** Gunakan **Session Pooler** (bukan Direct Connection) karena koneksi langsung hanya mendukung IPv6.

### 4. Uji Coba via curl
```bash
curl -X POST http://localhost:8082/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
    "sender_name": "Rangga",
    "sender_phone": "0812345678",
    "sender_address": "Jl. Sukarno Hatta No. 10",
    "origin_postal": "40286",
    "origin_city": "Bandung",
    "receiver_name": "Zharfan",
    "receiver_phone": "0898765432",
    "receiver_address": "Jl. Merdeka No. 42",
    "dest_postal": "10110",
    "dest_city": "Jakarta",
    "weight_actual": 1.5,
    "length": 10.0,
    "width": 10.0,
    "height": 10.0,
    "service_type": "REG",
    "payment_type": "TRANSFER"
  }'
```

---

## 🧪 Testing

### Unit Tests
```bash
go test ./internal/... -v
```

### Functional Tests (membutuhkan koneksi DB nyata)
```bash
set "TEST_DATABASE_URL=postgresql://..."
go test -v -tags=functional ./tests/functional/...
```

---

## 📨 Kafka Event

Setiap order berhasil dibuat akan mempublish event ke topic `order.created`:

```json
{
  "awb_number":     "BQ-2026-XKT-042",
  "transaction_id": "uuid",
  "sender_name":    "Rangga",
  "sender_address": "Jl. Sukarno Hatta No. 10",
  "origin_city":    "Bandung",
  "receiver_name":  "Zharfan",
  "dest_city":      "Jakarta",
  "service_type":   "REG",
  "total_price":    17000,
  "created_at":     "2026-06-18T20:49:00Z"
}
```

> Jika Kafka tidak tersedia, OMS tetap menyimpan order ke database dan mencatat warning log. Event bersifat best-effort (tidak blocking).

---

## 🐳 Docker

```bash
# Build image
docker build -f deployments/Dockerfile -t buroqet-oms .

# Jalankan via docker-compose (dari root buroqet/)
docker compose up -d order-management-service
```

---

## ☸️ Kubernetes

```bash
kubectl apply -f deployments/kubernetes/configmap.yaml
kubectl apply -f deployments/kubernetes/deployment.yaml
```

---

## 🗄️ Database Schema (Supabase)

- **Schema:** `oms`
- **Tabel:** `orders`
- **Generated Columns** (dikelola otomatis oleh PostgreSQL, tidak di-insert oleh aplikasi):
  - `volumetric_weight_kg` = `(length_cm × width_cm × height_cm) / 5000`
  - `is_cod` = `(payment_type = 'COD')`

---

## 📋 Order Status Lifecycle

```
ORDER_CREATED → PAYMENT_PENDING → PAYMENT_CONFIRMED → PICKED_UP
                                                    ↓
                                              ON_TRANSIT → AT_DESTINATION_HUB
                                                                    ↓
                                                          OUT_FOR_DELIVERY → DELIVERED
                                                                           ↘ FAILED
                                                                           ↘ RETURNED
```

---

## 🔗 Dependensi Antar Service

| Service | Interaksi | Protokol |
|---|---|---|
| **Pricing Service** (`:8084`) | Kalkulasi ongkir saat buat order | REST HTTP (sync) |
| **Dispatch Fleet Service** (`:8083`) | Consume event `order.created` | Kafka (async) |
| **Warehouse Service** (`:8087`) | Consume event `order.created` | Kafka (async) |
