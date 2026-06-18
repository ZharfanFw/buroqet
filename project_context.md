# Buroqet — PROJECT_CONTEXT.md

> **Single Source of Truth** untuk AI Coding Agent dan seluruh tim.
> Sumber kebenaran: Source Code → Dokumen Tahap 3 → Dokumen Tahap 1 → project_context.md lama.
> Last audited: 2026-06-18. Baca dokumen ini sebelum melakukan perubahan apapun.

---

## 1. Project Overview

### Tujuan Project
**Buroqet** adalah platform logistik berbasis **microservice** yang dibangun sebagai tugas akhir mata kuliah Cloud Computing Semester 4. Sistemnya meniru operasi perusahaan ekspedisi (JNE/SiCepat) dengan alur lengkap dari pembuatan order hingga bukti terima elektronik.

### Ruang Lingkup
- Manajemen order pengiriman dengan kalkulasi harga otomatis
- Pelacakan paket real-time (tracking) via Kafka + MongoDB + Redis
- Manajemen gudang (WMS): inbound, manifest, dispatch
- Penugasan kurir berbasis lokasi GPS (PostGIS)
- Electronic Proof of Delivery (ePOD) dengan upload foto
- Rekonsiliasi keuangan dan komisi kurir (settlement)

### Domain Utama
Logistik B2C: pengiriman paket dari pengirim ke penerima melalui jaringan hub & spoke.

---

## 2. System Architecture

### Tech Stack

| Layer | Teknologi | Keterangan |
|---|---|---|
| Backend Services | **Go (Golang)** | Tiap service mandiri, module terpisah |
| Frontend | **React 18 + Vite + TypeScript** | SPA, routing via React Router v6 |
| Database Relasional | **PostgreSQL 16** | Auth, Order, Dispatch, Pricing, Settlement, Warehouse |
| Database Dokumen | **MongoDB 7** | Tracking events (append-only log) |
| Cache | **Redis 7** | Status terakhir paket (tracking snapshot) |
| Message Broker | **Apache Kafka** (Confluent 7.6 + Zookeeper) | Event-driven antar service |
| HTTP Framework | **Gin** (Auth, Order, Pricing, ePOD) atau **net/http** (Tracking, Dispatch, Warehouse, Settlement) | Bervariasi per service |
| ORM | **GORM** (Auth, Order, Warehouse, Settlement) | `database/sql` untuk Dispatch |
| JWT | `github.com/golang-jwt/jwt/v5` | Auth service |
| Mock | `github.com/golang/mock` | Testing |
| Container | **Docker + Docker Compose** | Local & production |
| Orchestration | **Kubernetes** (Minikube) | Production (partial) |
| CI/CD | **Jenkins** | Tiap service punya Jenkinsfile |

### Diagram Arsitektur

```
┌──────────────────────────────────────────────────────────┐
│                   Frontend (React+Vite)                  │
│           http://localhost:3000 (Docker/prod)            │
│           http://localhost:5173 (dev server)             │
└───────────────────────┬──────────────────────────────────┘
                        │ HTTP / REST (JSON)
          ┌─────────────▼──────────────────────────┐
          │           Nginx Ingress                 │
          │        buroqet.local (K8s)              │
          └──┬─────┬─────┬─────┬─────┬─────┬──────┘
             │     │     │     │     │     │
        ┌────▼┐ ┌──▼──┐ ┌▼──┐ ┌▼──┐ ┌▼──┐ ┌▼───┐ ┌▼───┐ ┌▼───┐
        │Auth │ │Track│ │Ord│ │Dis│ │Pri│ │Set│ │ePOD│ │WMS │
        │:8080│ │:8081│ │:82│ │:83│ │:84│ │:85│ │:86 │ │:87 │
        └──┬──┘ └──┬──┘ └┬──┘ └┬──┘ └┬──┘ └┬──┘ └─┬──┘ └─┬──┘
           │       │     │     │     │     │    │      │
        ┌──▼──┐ ┌──▼──┐  └─────┴─────┴─────┴────┴──────┴───┐
        │Postg│ │Mongo│                Kafka                  │
        │ SQL │ │ DB  │        (Event Bus, Confluent 7.6)    │
        └─────┘ └──┬──┘  └─────────────────────────────────┘
                   │
                ┌──▼──┐
                │Redis│
                └─────┘
```

> **Port di atas adalah port Docker host-side** (sesuai docker-compose.yml).
> Beberapa service secara internal (di dalam container) mungkin mendengarkan port berbeda — lihat bagian per-service.

---

## 3. Business Workflow

### Pre Journey (Persiapan)
- Customer **register** dan **login** → Auth Service (JWT)
- Customer memilih layanan dan mengisi detail pengiriman

### First Mile (Pengambilan)
1. Customer membuat order → **Order Management Service**
   - OMS memanggil **Pricing Service** (REST) untuk kalkulasi ongkir
   - OMS menyimpan order ke PostgreSQL, status: `ORDER_CREATED`
   - OMS publish event `order.created` ke Kafka
2. Dispatch Fleet Service menerima event `order.created`
   - Mencari kurir terdekat (berbasis koordinat GPS + PostGIS)
   - Assign kurir ke pickup → endpoint `/v1/dispatch/assign`
3. Kurir pickup paket dari pengirim

### Mid Mile (Transit)
4. Paket masuk gudang asal → **Warehouse Service**
   - Handler: `POST /api/v1/inbound` → `ProcessInbound`
   - WMS menyimpan Package ke PostgreSQL, status: `INBOUND`
   - WMS publish event `package.arrived` ke Kafka
   - **Tracking Service** otomatis update status via Kafka
5. Paket dikelompokkan ke dalam **Manifest** (1 truk/rute)
   - Handler: `POST /api/v1/dispatch` → `DispatchManifest`
   - WMS update status semua AWB di manifest ke `ON_TRANSIT`
   - WMS publish event `manifest.dispatched` ke Kafka
6. Tracking Service mencatat event `AT_HUB` jika ada singgah di hub transit

### Last Mile (Pengiriman Akhir)
7. Kurir mengambil paket dari gudang tujuan, status `OUT_FOR_DELIVERY`
   - Dispatch Service publish event `package.dispatched` ke Kafka
   - Tracking Service update status
8. Kurir mengantar ke penerima → **ePOD Service**
   - Kurir upload foto bukti terima: `POST /upload`
   - ePOD publish event `package.delivered` ke Kafka
   - Tracking Service update status ke `DELIVERED`
9. **Settlement Service** dipicu (event atau manual) untuk mencatat komisi kurir
   - Memanggil Pricing Service untuk mendapatkan commission rate
   - Menyimpan `CommissionLog` ke PostgreSQL

### Status Lifecycle Paket (Tracking)
```
INBOUND → ON_TRANSIT → AT_HUB → OUT_FOR_DELIVERY → DELIVERED
                                                   ↘ FAILED
                                                   ↘ RETURNED
```

### Status Lifecycle Order (OMS)
```
ORDER_CREATED → PAYMENT_PENDING → PAYMENT_PAID → [masuk ke WMS]
                                               ↘ CANCELLED
```

---

## 4. Microservice Responsibilities

### 1. Auth Service — Port Docker `8080`

| Atribut | Detail |
|---|---|
| **Tujuan** | Autentikasi user, manajemen akun, penerbitan JWT |
| **Module Go** | `auth-service` |
| **HTTP Framework** | Gin |
| **ORM** | GORM |
| **Database** | PostgreSQL — DB: `buroqet` (shared instance), tabel `users` |
| **Producer Kafka** | — (tidak ada) |
| **Consumer Kafka** | — (tidak ada) |
| **Dependency** | PostgreSQL |

**Domain Model:**
```go
type User struct {
    ID        uint      `gorm:"primaryKey"`
    Name      string    `gorm:"size:100;not null"`
    Email     string    `gorm:"size:100;uniqueIndex;not null"`
    Password  string    `gorm:"not null"` // bcrypt hashed; json:"-"
    Role      string    `gorm:"size:20;not null"` // pelanggan | kurir | admin
    CreatedAt time.Time
}
```

**Endpoint:**
```
POST /api/v1/auth/register  → Daftarkan user baru (name, email, password, role)
POST /api/v1/auth/login     → Login, return access_token + refresh_token + user
GET  /health                → Liveness probe
```

> **Catatan Audit:** Handler hanya memiliki `Register` dan `Login`. Endpoint `/api/v1/auth/refresh` ada di dokumentasi lama **tapi belum diimplementasikan di handler**. Jangan mengasumsikan endpoint ini ada sampai handler-nya dibuat.

**Interfaces:**
```go
type UserRepository interface {
    Create(user *User) error
    FindByEmail(email string) (*User, error)
    GetActiveDrivers() ([]User, error)
}

type AuthService interface {
    Register(name, email, password, role string) error
    Login(email, password string) (accessToken, refreshToken string, user *User, err error)
    ValidateToken(token string) (*User, error)
}
```

---

### 2. Tracking Service — Port Docker `8081`

| Atribut | Detail |
|---|---|
| **Tujuan** | Mencatat dan membaca riwayat perjalanan paket secara real-time |
| **Module Go** | `tracking-service` |
| **HTTP Framework** | `net/http` (bukan Gin) |
| **Database** | MongoDB — DB: `tracking_db`, koleksi: `tracking_events` |
| **Cache** | Redis — key: `tracking:status:{awb}`, TTL: 7 hari |
| **Producer Kafka** | `tracking.updated` (downstream notification) |
| **Consumer Kafka** | `package.inbound`, `package.dispatched`, `package.delivered`, `manifest.arrived` |
| **Dependency** | MongoDB, Redis, Kafka |

**Penting:** Tracking service menggunakan `APP_PORT` env var, default internal `8080`. Docker-compose memetakan ke host port `8081`.

**Domain Model (MongoDB document):**
```go
type TrackingEvent struct {
    ID          string    `bson:"_id,omitempty"  json:"id"`
    AWB         string    `bson:"awb"            json:"awb"`
    Status      string    `bson:"status"         json:"status"`
    Location    string    `bson:"location"       json:"location"`
    HubID       string    `bson:"hub_id"         json:"hub_id"`
    Description string    `bson:"description"    json:"description"`
    Timestamp   time.Time `bson:"timestamp"      json:"timestamp"`
    CreatedAt   time.Time `bson:"created_at"     json:"created_at"`
    Source      string    `bson:"source"         json:"source"` // WMS | Dispatch | e-POD
}
// CATATAN AUDIT: Tidak ada field Latitude/Longitude di struct ini.
// Latitude/Longitude ada di KafkaTrackingPayload (payload lama) tapi bukan di TrackingEvent.

type TrackingStatus struct { // disimpan di Redis
    AWB           string    `json:"awb"`
    CurrentStatus string    `json:"current_status"`
    LastLocation  string    `json:"last_location"`
    LastUpdated   time.Time `json:"last_updated"`
}

type KafkaTrackingPayload struct { // payload yang diterima dari Kafka
    AWB       string    `json:"awb"`
    Status    string    `json:"status"`
    HubID     string    `json:"hub_id"`
    Location  string    `json:"location"`
    Timestamp time.Time `json:"timestamp"`
    Source    string    `json:"source"`
    // CATATAN AUDIT: Tidak ada field Description di payload ini.
}
```

**Valid Status:** `INBOUND` | `ON_TRANSIT` | `AT_HUB` | `OUT_FOR_DELIVERY` | `DELIVERED` | `FAILED` | `RETURNED`

**Endpoint:**
```
POST /api/v1/tracking/events           → Tambah event baru (dari service lain atau manual)
GET  /api/v1/tracking/{awb}/history    → Riwayat lengkap paket (dari MongoDB)
GET  /api/v1/tracking/{awb}/status     → Status terakhir (Redis → MongoDB fallback)
GET  /health                           → Liveness probe
GET  /ready                            → Readiness probe
```

**CORS yang dikonfigurasi:**
- `http://localhost:3000` (prod frontend)
- `http://localhost:8080`
- `http://127.0.0.1:8080`
- `null` (file:// origin)

> **Catatan:** `http://localhost:5173` (Vite dev) **tidak** ada di CORS whitelist default tracking service. Tambahkan via env var `ALLOWED_ORIGIN` jika diperlukan.

---

### 3. Order Management Service (OMS) — Port Docker `8082`

| Atribut | Detail |
|---|---|
| **Tujuan** | Pembuatan dan manajemen order pengiriman |
| **Module Go** | `order-management-service` (package domain: `model`) |
| **HTTP Framework** | Gin |
| **ORM** | GORM |
| **Database** | PostgreSQL — DB: `buroqet`, tabel: `orders` |
| **Producer Kafka** | `order.created` (payload: `OrderCreatedEvent`) |
| **Consumer Kafka** | — (tidak ada) |
| **Dependency** | PostgreSQL, Kafka, Pricing Service (REST HTTP call) |

**Domain Model:**
```go
type Order struct {
    ID              uint        `gorm:"primaryKey;autoIncrement"`
    AWBNumber       string      `gorm:"uniqueIndex;not null"` // format: BQ-YYYY-XXX-NNN
    TransactionID   string      `gorm:"uniqueIndex;not null"`
    Status          OrderStatus // ORDER_CREATED | PAYMENT_PENDING | PAYMENT_PAID | CANCELLED
    // Sender
    SenderName      string
    SenderPhone     string
    SenderAddress   string
    OriginCity      string
    OriginPostal    string
    // Receiver
    ReceiverName    string
    ReceiverPhone   string
    ReceiverAddress string
    DestCity        string
    DestPostal      string
    // Package
    WeightActual    float64
    WeightVolumetri float64 // dihitung: L×W×H/6000
    Length          float64
    Width           float64
    Height          float64
    // Service & Payment
    ServiceType     ServiceType // REGULER | EXPRESS
    PaymentType     PaymentType // COD | NON_COD
    TotalPrice      float64
    PaymentURL      string
    CreatedAt       time.Time
    UpdatedAt       time.Time
}
```

**Endpoint:**
```
POST /api/v1/orders        → Buat order baru (CreateOrderRequest), return 201 + AWB
GET  /api/v1/orders/:awb   → Ambil detail order berdasarkan AWB
GET  /health               → Liveness probe
```

**Kafka Event yang Dipublish — `order.created`:**
```json
{
  "awb_number":    "BQ-2024-JKT-001",
  "transaction_id": "TXN-...",
  "sender_name":   "...",
  "sender_address":"...",
  "origin_city":   "Jakarta",
  "receiver_name": "...",
  "dest_city":     "Bandung",
  "service_type":  "REGULER",
  "total_price":   25000,
  "created_at":    "2024-01-15T08:00:00Z"
}
```

**Interaksi Sinkron:** OMS memanggil Pricing Service via HTTP (`PRICING_SERVICE_URL`) untuk kalkulasi ongkir saat pembuatan order.

---

### 4. Dispatch & Fleet Service — Port Docker `8083`

| Atribut | Detail |
|---|---|
| **Tujuan** | Penugasan kurir terdekat ke pickup, manajemen armada kurir |
| **Module Go** | `dispatch-fleet` (bukan `dispatch-fleet-service`) |
| **HTTP Framework** | `net/http` (bukan Gin) |
| **ORM/DB Driver** | `database/sql` + `github.com/lib/pq` |
| **Database** | PostgreSQL dengan **PostGIS** (geospatial untuk pencarian kurir terdekat) |
| **Producer Kafka** | `package.dispatched` (saat kurir berangkat antar paket) |
| **Consumer Kafka** | `order.created` (untuk trigger assign kurir pickup) |
| **Dependency** | PostgreSQL + PostGIS |

**Domain Model:**
```go
type Courier struct {
    ID              string
    Name            string
    CurrentLocation Point // {Longitude, Latitude}
    Status          CourierStatus // available | assigned | on_delivery | offline
    CreatedAt       time.Time
    UpdatedAt       time.Time
}

type AssignCourierResult struct {
    Courier        *Courier
    DistanceMeters float64
}
```

**Endpoint:**
```
POST /v1/dispatch/assign   → Assign kurir terdekat ke titik pickup
                             Body: { pickup_lat, pickup_lon, radius_meters }
                             Response: AssignCourierResult
```

> **Catatan Audit:** Module name di go.mod adalah `dispatch-fleet`, port default internal adalah `8081` (bukan `8083`). Port `8083` hanya berlaku di docker-compose host mapping. Tidak ada `/api/v1` prefix — path langsung `/v1/dispatch/assign`.

---

### 5. Pricing Service — Port Docker `8084`

| Atribut | Detail |
|---|---|
| **Tujuan** | Kalkulasi ongkos kirim + commission rate untuk Settlement |
| **Module Go** | `pricing-service` |
| **HTTP Framework** | Gin |
| **Database** | PostgreSQL — rate card |
| **Producer Kafka** | — (tidak ada) |
| **Consumer Kafka** | — (tidak ada) |
| **Dependency** | PostgreSQL |

**Domain Model:**
```go
type CalculationRequest struct {
    OriginPostalCode      string
    DestinationPostalCode string
    OriginLat             float64
    OriginLon             float64
    DestinationLat        float64
    DestinationLon        float64
    WeightKG              float64
    LengthCM              float64
    WidthCM               float64
    HeightCM              float64
    ServiceType           string  // REGULER | EXPRESS
    UseInsurance          bool
    PromoCode             string
}

type CalculationResponse struct {
    BaseTariff float64 `json:"base_tariff"`
    Insurance  float64 `json:"insurance"`
    Discount   float64 `json:"discount"`
    Total      float64 `json:"total"`
    Estimated  string  `json:"estimated"` // SLA estimate
}
```

**Formula:** `max(weight_actual, weight_volumetric) × rate_per_kg` + insurance + promo discount

**Endpoint:**
```
POST /pricing/calculate    → Kalkulasi ongkir
GET  /health               → Liveness probe
```

> **Catatan Audit:** Port default internal `8081` (di cmd/main.go), docker-compose host port adalah `8084`. Saat ini cmd/main.go menggunakan mock repository.

---

### 6. Settlement Service — Port Docker `8085`

| Atribut | Detail |
|---|---|
| **Tujuan** | Rekonsiliasi keuangan, pencatatan komisi kurir, laporan penghasilan |
| **Module Go** | `settlement-service` |
| **HTTP Framework** | `net/http` (bukan Gin) |
| **ORM** | GORM |
| **Database** | PostgreSQL — tabel `commission_logs` |
| **Producer Kafka** | — (tidak ada) |
| **Consumer Kafka** | `package.delivered` (triggered by event) |
| **Dependency** | PostgreSQL, Pricing Service (REST HTTP call untuk commission rate) |

**Domain Model (diinfer dari service code — file domain Go ada bug berupa mis-named file):**
```go
// File settlement.go di domain/ berisi docker-compose YAML (bug di repo)
// Model aktual diinfer dari service dan main.go:

type CommissionLog struct {
    ID          string    // UUID
    CourierID   string
    AWB         string
    Amount      float64   // dari commission rate Pricing Service
    Status      string    // PENDING | PAID
    DeliveredAt time.Time
}

type CourierSummary struct {
    // Summary penghasilan kurir
}
```

**Endpoint:**
```
POST /api/v1/commissions              → Proses komisi delivery
                                        Body: { courier_id, awb, service_type }
GET  /api/v1/couriers/{courierID}/earnings → Rekap penghasilan kurir
GET  /health                          → Liveness probe
GET  /ready                           → Readiness probe
```

> **Bug diketahui:** `settlement-service/internal/domain/settlement.go` berisi YAML docker-compose, bukan Go code. Harus diperbaiki.

---

### 7. ePOD Service (Electronic Proof of Delivery) — Port Docker `8086`

| Atribut | Detail |
|---|---|
| **Tujuan** | Upload foto + tanda tangan sebagai bukti terima paket |
| **Module Go** | `epod-service` |
| **HTTP Framework** | Gin |
| **Database** | **Object Storage** (file foto disimpan lokal di `./uploads/` saat dev, harusnya MinIO/S3 di prod) |
| **Producer Kafka** | `package.delivered` (dipublish saat upload foto berhasil) |
| **Consumer Kafka** | — (tidak ada) |
| **Dependency** | Kafka |

**Domain Model:**
```go
type UploadRequest struct {
    AWB       string
    CourierID string
    Latitude  float64
    Longitude float64
    FileName  string
}

type UploadResponse struct {
    Status   string `json:"status"`   // "SUCCESS"
    ImageURL string `json:"image_url"` // URL file yang diupload
}
```

**Endpoint:**
```
POST /upload    → Upload foto bukti terima, publish event package.delivered ke Kafka
```

> **Catatan Audit:** Port internal di cmd/main.go adalah `:8080`, docker-compose host mapping ke `8086`. Tidak ada prefix `/api/v1` — path langsung `/upload`.

---

### 8. Warehouse Service (WMS) — Port Docker `8087`

| Atribut | Detail |
|---|---|
| **Tujuan** | Manajemen inbound paket, pengelompokan manifest, dispatch outbound |
| **Module Go** | `warehouse-service` |
| **HTTP Framework** | `net/http` (bukan Gin) |
| **ORM** | GORM |
| **Database** | PostgreSQL — tabel `packages`, `manifests` |
| **Producer Kafka** | `package.arrived` (inbound), `manifest.dispatched` (outbound) |
| **Consumer Kafka** | `order.created` (untuk pre-register paket) |
| **Dependency** | PostgreSQL, Kafka |

**Domain Model:**
```go
type Package struct {
    ID         string    `gorm:"primaryKey"`
    AWB        string    `gorm:"uniqueIndex;not null"`
    HubID      string
    ManifestID *string   // nullable, diisi saat dimasukkan ke manifest
    Status     string    // INBOUND | OUTBOUND | ON_TRANSIT
    ScannedAt  time.Time
    CreatedAt  time.Time
    UpdatedAt  time.Time
}

type Manifest struct {
    ID          string    `gorm:"primaryKey"`
    TruckID     string
    OriginHubID string
    DestHubID   string
    Status      string    // OPEN | DISPATCHED
    Packages    []Package `gorm:"foreignKey:ManifestID"`
    CreatedAt   time.Time
}
```

**Endpoint:**
```
POST /api/v1/inbound    → Proses paket masuk gudang (awb, hub_id), publish package.arrived
POST /api/v1/dispatch   → Dispatch manifest (manifest_id), publish manifest.dispatched
GET  /health            → Liveness probe
GET  /ready             → Readiness probe
```

> **Catatan Audit — Kafka Topic:** Source code WMS mempublish `package.arrived` (bukan `package.inbound`) untuk inbound, dan `manifest.dispatched` (bukan `manifest.arrived`) untuk dispatch. Tracking service subscribe ke `package.inbound` dan `manifest.arrived`. Ada **ketidaksesuaian topic** antara WMS (producer) dan Tracking (consumer) — ini perlu diselaraskan.

---

## 5. Communication Rules

### Synchronous (REST)

**Kapan digunakan:** Saat membutuhkan respons langsung untuk melanjutkan proses bisnis.

| Caller | Called | Endpoint | Tujuan |
|---|---|---|---|
| Order Management Service | Pricing Service | `POST /pricing/calculate` | Kalkulasi ongkir saat buat order |
| Settlement Service | Pricing Service | HTTP GET commission rate | Ambil commission rate per service type |
| Frontend | Auth Service | `POST /api/v1/auth/login` | Login |
| Frontend | Tracking Service | `GET /api/v1/tracking/{awb}/status` | Cek status paket |

**Format:** JSON request/response. Auth via JWT Bearer token untuk endpoint protected.

### Asynchronous (Kafka)

**Kapan digunakan:** Untuk event yang tidak butuh respons langsung, atau untuk decoupling antar service.

| Topic | Producer | Consumer | Trigger |
|---|---|---|---|
| `order.created` | Order Service | Dispatch Service, (WMS) | Order baru berhasil dibuat |
| `package.arrived` | Warehouse Service | Tracking Service | Paket scan masuk gudang |
| `manifest.dispatched` | Warehouse Service | Tracking Service | Manifest (kumpulan paket) dikirim |
| `package.dispatched` | Dispatch Service | Tracking Service | Kurir berangkat antar paket |
| `package.delivered` | ePOD Service | Tracking Service, Settlement Service | Paket berhasil diterima penerima |
| `tracking.updated` | Tracking Service | (downstream notification) | Status paket diupdate |

> **Catatan Penting — Ketidaksesuaian Topic:** WMS mempublish `package.arrived` dan `manifest.dispatched`, namun tracking service subscribe ke `package.inbound` dan `manifest.arrived`. Topic-topic ini **belum selaras** dan harus diselesaikan sebelum integrasi end-to-end.

**Key Kafka:** AWB Number digunakan sebagai message key untuk semua event per paket → menjamin ordering per AWB dalam satu partition.

**Retry Policy:** Tidak ada retry otomatis di current implementation — producer mencatat WARNING dan lanjut (eventual consistency pattern). Implementasi retry untuk production.

**Idempotency:** Belum diimplementasikan — perlu pengecekan duplikasi di consumer.

---

## 6. Database Rules

### Prinsip Dasar
- **Database per service** — setiap service memiliki schema/tabel sendiri di PostgreSQL atau database tersendiri (MongoDB)
- **Tidak boleh cross query** — service A tidak boleh langsung query database service B
- **Data sharing** via event Kafka atau REST API — bukan direct DB access

### Database per Service

| Service | Database | DB Name / Collection |
|---|---|---|
| Auth | PostgreSQL | `buroqet` → tabel `users` |
| Order | PostgreSQL | `buroqet` → tabel `orders` |
| Dispatch | PostgreSQL + PostGIS | `dispatch_db` |
| Pricing | PostgreSQL | `buroqet` |
| Settlement | PostgreSQL | `settlement_db` → tabel `commission_logs` |
| Warehouse | PostgreSQL | `wms_db` → tabel `packages`, `manifests` |
| Tracking | MongoDB | `tracking_db` → koleksi `tracking_events` |
| Tracking (cache) | Redis | key pattern: `tracking:status:{awb}` |

> **Catatan:** Auth, Order, dan Pricing saat ini menggunakan shared PostgreSQL instance (`buroqet` DB di docker-compose), namun menggunakan tabel terpisah. Settlement dan Warehouse menggunakan DB name berbeda dan diintensi untuk dipisah.

### Transaction Boundary
- Setiap service mengelola transaksi databasenya sendiri
- Tidak ada distributed transaction — gunakan Saga pattern / eventual consistency via Kafka

---

## 7. API Rules

- **Format:** JSON (Content-Type: application/json)
- **Auth:** JWT Bearer token di header `Authorization: Bearer <token>`
- **Versioning:** `/api/v1/` untuk semua endpoint (kecuali Dispatch: `/v1/`, ePOD: tidak ada prefix)
- **Validation:** Gin binding tags (`binding:"required"`, `oneof=...`) untuk request body
- **Error Response Format:**
  ```json
  { "error": "pesan error" }
  // atau (OMS):
  { "success": false, "error": "pesan error" }
  ```
- **Success Response Format:**
  ```json
  { "data": {...} }
  // atau (OMS):
  { "success": true, "data": {...} }
  ```
- **HTTP Status Codes:**
  - `200 OK` — berhasil (GET, PUT)
  - `201 Created` — resource berhasil dibuat (POST)
  - `400 Bad Request` — input tidak valid
  - `401 Unauthorized` — tidak ada/invalid JWT
  - `404 Not Found` — resource tidak ditemukan
  - `405 Method Not Allowed`
  - `500 Internal Server Error`

---

## 8. Kafka Event Catalog

### `order.created`
| Field | Value |
|---|---|
| **Topic** | `order.created` |
| **Producer** | Order Management Service |
| **Consumer** | Dispatch Fleet Service |
| **Trigger** | Order berhasil dibuat dan disimpan ke DB |
| **Key** | AWB Number |

**Payload:**
```json
{
  "awb_number":     "BQ-2024-JKT-001",
  "transaction_id": "TXN-...",
  "sender_name":    "Budi Santoso",
  "sender_address": "Jl. Sudirman No.1, Jakarta",
  "origin_city":    "Jakarta",
  "receiver_name":  "Ani Wijaya",
  "dest_city":      "Bandung",
  "service_type":   "REGULER",
  "total_price":    25000,
  "created_at":     "2024-01-15T08:00:00Z"
}
```

---

### `package.arrived`
| Field | Value |
|---|---|
| **Topic** | `package.arrived` |
| **Producer** | Warehouse Service |
| **Consumer** | Tracking Service (setelah topic diselaraskan) |
| **Trigger** | Paket di-scan masuk gudang (`POST /api/v1/inbound`) |
| **Key** | AWB |

**Payload:**
```json
{
  "awb":    "BQ-2024-JKT-001",
  "hub_id": "HUB-JKT-01",
  "status": "INBOUND"
}
```

---

### `manifest.dispatched`
| Field | Value |
|---|---|
| **Topic** | `manifest.dispatched` |
| **Producer** | Warehouse Service |
| **Consumer** | Tracking Service (setelah topic diselaraskan) |
| **Trigger** | Manifest di-dispatch (`POST /api/v1/dispatch`) |
| **Key** | Manifest ID |

**Payload:**
```json
{
  "manifest_id": "MNF-...",
  "awbs":        ["BQ-2024-JKT-001", "BQ-2024-JKT-002"],
  "status":      "ON_TRANSIT"
}
```

---

### `package.dispatched`
| Field | Value |
|---|---|
| **Topic** | `package.dispatched` |
| **Producer** | Dispatch Fleet Service |
| **Consumer** | Tracking Service |
| **Trigger** | Kurir berangkat antar paket ke penerima |

**Payload (KafkaTrackingPayload format):**
```json
{
  "awb":       "BQ-2024-JKT-001",
  "status":    "OUT_FOR_DELIVERY",
  "hub_id":    "HUB-BDG-01",
  "location":  "Gudang Bandung",
  "timestamp": "2024-01-15T14:00:00Z",
  "source":    "Dispatch"
}
```

---

### `package.delivere


---

### `tracking.updated`
| Field | Value |
|---|---|
| **Topic** | `tracking.updated` |
| **Producer** | Tracking Service |
| **Consumer** | Downstream services (notification, dll.) |
| **Trigger** | Status paket diupdate |

---

## 9. Frontend Rules

| Atribut | Detail |
|---|---|
| **Framework** | React 18 + TypeScript |
| **Build Tool** | Vite |
| **Routing** | React Router v6 (`BrowserRouter`) |
| **State Management** | Zustand (auth state: `useAuthStore`) |
| **HTTP Client** | Axios + interceptor (auto-attach JWT, handle 401/refresh) |
| **Peta** | React Leaflet + Leaflet.js (OpenStreetMap / CartoDB dark tile) |

### Environment Variables
```env
VITE_API_BASE_URL=/api              # Via Nginx proxy di prod
VITE_TRACKING_API=http://localhost:8081  # Direct ke tracking service (dev)
```

### Halaman dan Routes

| Route | Komponen | Auth | Status |
|---|---|---|---|
| `/login` | `LoginPage` | Public | ✅ Selesai |
| `/tracking` | `TrackingPage` | Protected (via PrivateRoute) | ✅ UI selesai |
| `/orders` | `OrderPage` | Protected | ⚠️ Placeholder |
| `/dispatch` | `DispatchPage` | Protected | ⚠️ Placeholder |
| `/pricing` | `PricingPage` | Protected | ⚠️ Placeholder |
| `/settlement` | `SettlementPage` | Protected | ⚠️ Placeholder |
| `/epod` | `EpodPage` | Protected | ⚠️ Placeholder |
| `/warehouse` | `WarehousePage` | Protected | ⚠️ Placeholder |
| `/` | redirect → `/tracking` | — | — |

> **Catatan Audit:** Route `/tracking` membutuhkan autentikasi (PrivateRoute). Dokumentasi lama menyebutnya "Public" — kode menunjukkan Protected. Sesuaikan dokumentasi dengan implementasi.

### Struktur Frontend (`frontend/src/`)
```
src/
├── App.tsx                   # Router + PrivateRoute + layout wrapping
├── main.tsx                  # Entry point
├── index.css                 # Global styles
├── types/index.ts            # Shared TypeScript interfaces
├── services/api-client.ts    # Axios instance + JWT interceptor
├── store/auth.store.ts       # Zustand: isAuthenticated, fetchMe, dll.
├── utils/api-config.ts       # API endpoint constants
├── components/layout/        # MainLayout, AuthLayout, Sidebar
└── pages/
    ├── auth/LoginPage.tsx
    ├── tracking/
    │   ├── TrackingPage.tsx      # Search UI + state
    │   ├── TrackingResult.tsx    # Status hero + stepper + timeline
    │   ├── TrackingMap.tsx       # Leaflet interaktif
    │   └── *.css
    ├── order/OrderPage.tsx       # placeholder
    ├── dispatch/DispatchPage.tsx # placeholder
    ├── pricing/PricingPage.tsx   # placeholder
    ├── settlement/SettlementPage.tsx # placeholder
    ├── epod/EpodPage.tsx         # placeholder
    └── warehouse/WarehousePage.tsx # placeholder
```

---

## 10. Coding Standards

### Pola Arsitektur

Setiap service mengikuti **Clean Architecture** dengan layer:

```
cmd/main.go           → Entry point + Dependency Injection + routing
internal/
├── domain/           → Model, interface, konstanta — TANPA external dependency
├── handler/          → HTTP handler — decode request, call service, encode response
├── service/          → Business logic — orchestrasi, validasi bisnis
├── repository/       → DB access (PostgreSQL/MongoDB)
├── cache/            → Redis access (tracking service)
├── kafka/            → Producer/Consumer implementation
mocks/                → Mock implementasi untuk unit test
tests/functional/     → Functional/integration test
```

### Prinsip
- **Dependency Injection:** Service menerima interface (bukan konkret), diinject dari `cmd/main.go`
- **Interface-first:** Semua dependency direpresentasikan sebagai interface agar bisa di-mock
- **SOLID:** Single Responsibility per layer
- **DRY/KISS:** Tidak duplikasi logika antar layer
- **No hardcode:** Semua config dari environment variable
- **Graceful shutdown:** HTTP server dan Kafka consumer shutdown bersih (contoh: tracking service)

### Environment Variables (umum)
```env
DB_HOST      → Host PostgreSQL
DB_PORT      → Port PostgreSQL (default: 5432)
DB_USER      → User PostgreSQL
DB_PASSWORD  → Password PostgreSQL
DB_NAME      → Nama database
APP_PORT     → Port yang didengarkan service
KAFKA_BROKER → Kafka broker address
JWT_SECRET   → Secret key JWT (jangan hardcode, jangan commit)
```

### Error Handling
- Return error dari setiap layer — jangan panic kecuali startup failure
- Kafka publish failure: log WARNING + lanjut (eventual consistency pattern)
- DB failure saat startup: `log.Fatalf` (fail fast)

### Logging
- Gunakan `log.Printf` / `log.Println` dari standard library
- Format: `[KAFKA] topic=X key=Y value=Z` untuk Kafka events
- Tidak ada structured logging library saat ini

---

## 11. Development Workflow

```
Requirement
    ↓
Domain Design (domain/ layer — interface + model)
    ↓
Repository/Cache Implementation
    ↓
Service / Business Logic
    ↓
Handler (HTTP) + Kafka Producer/Consumer
    ↓
Unit Test (mock dependencies)
    ↓
Dockerfile + docker-compose entry
    ↓
Frontend (page + API integration)
    ↓
Functional/Integration Test
    ↓
Jenkinsfile (CI/CD)
    ↓
Kubernetes manifest (jika diperlukan)
```

### Cara Menjalankan

**Mode Dev (Hybrid — hanya infra Docker, service lokal):**
```bash
# 1. Start infrastructure
make dev-infra
# atau: docker compose up -d postgres mongodb redis zookeeper kafka

# 2. Start tracking service (contoh)
cd tracking-service
MONGO_URI="mongodb://buroqet:buroqet123@localhost:27017" \
MONGO_DB="tracking_db" \
REDIS_ADDR="localhost:6379" \
REDIS_PASSWORD="buroqet123" \
KAFKA_BROKER="localhost:9092" \
APP_PORT="8081" \
go run ./cmd/main.go

# 3. Start frontend dev server
cd frontend && npm run dev
# → http://localhost:5173
```

**Mode Full Docker:**
```bash
make dev    # atau: docker compose up -d
# Frontend: http://localhost:3000
```

### Make Commands
```bash
make help              # Lihat semua command
make dev-infra         # Jalankan infra saja (DB, Kafka)
make dev               # Full stack Docker
make frontend          # Jalankan frontend dev server (npm run dev)
make build             # Build semua Docker images
make build-fe          # Build frontend production
make test              # Run semua unit test (auth, tracking, order, dispatch)
make test-auth         # Test auth-service saja
make test-tracking     # Test tracking-service saja
make test-order        # Test order-management-service saja
make test-dispatch     # Test dispatch-fleet-service saja
make test-pricing      # Test pricing-service saja
make test-settlement   # Test settlement-service saja
make test-epod         # Test epod-service saja
make test-warehouse    # Test warehouse-service saja
make clean             # Stop + hapus semua container dan volume
make logs              # Tail semua logs
```

---

## 12. File Structure

```
buroqet/
├── Makefile                       # Unified dev commands
├── docker-compose.yml             # Full stack orchestration (semua services + infra)
├── project_context.md             # Dokumen ini
├── .gitignore
│
├── auth-service/                  # Go 1.23, Gin, GORM, PostgreSQL, JWT
│   ├── cmd/main.go                # Entry + DI + routing
│   ├── internal/
│   │   ├── domain/user.go         # User model + interfaces
│   │   ├── handler/auth_handler.go
│   │   ├── repository/user_postgres.go
│   │   └── service/auth_service.go
│   ├── mocks/
│   ├── tests/functional/
│   ├── deployments/Dockerfile
│   └── Jenkinsfile
│
├── tracking-service/              # Go 1.22, net/http, MongoDB, Redis, Kafka
│   ├── cmd/main.go                # Entry + DI + routing + Kafka consumer loop + CORS
│   ├── internal/
│   │   ├── domain/tracking.go     # Models, interfaces, status constants
│   │   ├── handler/tracking_handler.go
│   │   ├── service/
│   │   ├── repository/            # MongoDB implementation
│   │   ├── cache/                 # Redis implementation
│   │   └── kafka/                 # Producer + Consumer (placeholder, TODO: kafka-go)
│   ├── mocks/
│   ├── tests/functional/
│   └── deployments/Dockerfile
│
├── order-management-service/      # Go, Gin, GORM, PostgreSQL, Kafka
│   ├── cmd/main.go
│   ├── internal/
│   │   ├── domain/order.go        # package model (bukan package domain!)
│   │   ├── handler/order_handler.go
│   │   ├── service/
│   │   ├── repository/
│   │   └── kafka/producer.go      # Real kafka-go implementation
│   ├── mocks/
│   └── deployments/Dockerfile
│
├── dispatch-fleet-service/        # Go, net/http, database/sql, PostgreSQL+PostGIS
│   ├── cmd/main.go                # module: dispatch-fleet
│   ├── internal/
│   │   ├── domain/fleet.go + interfaces.go + errors.go
│   │   ├── handler/dispatch_handler.go
│   │   ├── repository/
│   │   └── service/
│   ├── mocks/
│   └── deployments/Dockerfile
│
├── pricing-service/               # Go, Gin, GORM, PostgreSQL
│   ├── cmd/main.go                # Saat ini menggunakan mock repository
│   ├── internal/
│   │   ├── domain/pricing.go
│   │   ├── handler/pricing_handler.go
│   │   └── service/
│   ├── mocks/
│   └── deployments/Dockerfile
│
├── settlement-service/            # Go, net/http, GORM, PostgreSQL
│   ├── cmd/main.go
│   ├── internal/
│   │   ├── domain/settlement.go   # ⚠️ BUG: file ini berisi YAML, bukan Go code
│   │   ├── handler/settlement_handler.go
│   │   ├── service/settlement_service.go
│   │   ├── repository/
│   │   └── pricing/               # HTTP client ke Pricing Service
│   ├── mocks/
│   └── deployments/Dockerfile
│
├── epod-service/                  # Go, Gin, Local storage, Kafka
│   ├── cmd/main.go                # Port internal: 8080
│   ├── internal/
│   │   ├── domain/epod.go
│   │   ├── handler/epod_handler.go
│   │   ├── service/epod_service.go
│   │   ├── storage/               # Local file storage (placeholder: MinIO/S3)
│   │   └── kafka/
│   ├── mocks/
│   └── deployments/Dockerfile
│
├── warehouse-service/             # Go, net/http, GORM, PostgreSQL, Kafka
│   ├── cmd/main.go
│   ├── internal/
│   │   ├── domain/warehouse.go    # Package, Manifest, interfaces
│   │   ├── handler/warehouse_handler.go
│   │   ├── service/warehouse_service.go
│   │   ├── repository/
│   │   └── kafka/producer.go      # Placeholder (log only, TODO: real Kafka)
│   ├── mocks/
│   └── deployments/Dockerfile
│
├── frontend/                      # React 18 + Vite + TypeScript
│   ├── src/
│   │   ├── App.tsx                # Router + PrivateRoute + layout
│   │   ├── types/index.ts
│   │   ├── services/api-client.ts # Axios + JWT interceptor
│   │   ├── store/auth.store.ts    # Zustand auth state
│   │   ├── utils/api-config.ts    # API endpoint constants
│   │   ├── components/layout/
│   │   └── pages/                 # auth/, tracking/, order/, ...
│   ├── .env.example
│   ├── Dockerfile                 # Nginx multi-stage build
│   └── nginx.conf
│
└── infra/
    └── kubernetes/
        ├── ingress.yaml           # Nginx Ingress routing ke semua service
        ├── namespace.yaml
        └── minikube-jenkins.yaml  # Jenkins di Minikube
```

---

## 13. AI Working Rules

Saat membantu project ini, AI **wajib**:

1. **Baca PROJECT_CONTEXT.md ini terlebih dahulu** sebelum membuat perubahan apapun
2. **Ikuti arsitektur yang ada:** Clean Architecture, interface-based DI, layer separation
3. **Jaga backward compatibility:** Jangan ubah API contract tanpa alasan kuat
4. **Ikuti style code per service:** Perhatikan apakah service menggunakan Gin atau `net/http`, GORM atau `database/sql`
5. **Jangan buat dependency baru** tanpa alasan kuat — cek go.mod terlebih dahulu
6. **Jangan pindahkan business logic** antar layer (service ke handler, dll.)
7. **Jangan ubah kontrak Kafka** (topic name, payload structure) tanpa update consumers
8. **Jelaskan dampak** perubahan sebelum melakukan refactor besar
9. **Perhatikan inkonsistensi yang diketahui** (lihat bagian Decision Log) sebelum melaporkan sebagai bug baru

---

## 14. Constraints

AI **tidak boleh**:

- Mengubah nama service (auth-service, tracking-service, dll.)
- Mengubah database yang digunakan per service
- Mengubah business workflow (urutan Pre/First/Mid/Last mile)
- Mengubah Kafka event name tanpa meng-update producer DAN consumer
- Mengubah API endpoint path/method tanpa alasan kuat
- Menggabungkan service (anti-monolith)
- Mengakses database service lain secara langsung
- Menambahkan library baru tanpa alasan
- Membuat asumsi tentang fitur yang belum diimplementasikan

---

## 15. Decision Log

### Mengapa Kafka untuk event antar service?
**Keputusan:** Menggunakan Apache Kafka (Confluent 7.6) sebagai message broker.
**Alasan:** Decoupling antar service, durability event, mendukung multiple consumer pada event yang sama (mis: `package.delivered` dikonsumsi oleh Tracking DAN Settlement sekaligus). Kafka juga mendukung replay event untuk audit trail.

### Mengapa setiap service memiliki database sendiri?
**Keputusan:** Database per service (PostgreSQL per service, MongoDB untuk tracking).
**Alasan:** Prinsip microservice — bounded context. Mencegah coupling database, memungkinkan scaling independen, memudahkan migrasi schema per service.

### Mengapa Tracking menggunakan MongoDB?
**Keputusan:** MongoDB untuk koleksi `tracking_events`.
**Alasan:** Tracking events bersifat append-only (log pattern), schema-less (setiap event bisa punya field berbeda), dan MongoDB lebih efisien untuk query "ambil semua events untuk AWB X terurut berdasarkan waktu" dibanding PostgreSQL.

### Mengapa Tracking menggunakan Redis?
**Keputusan:** Redis sebagai cache status terakhir paket.
**Alasan:** Query status terakhir paket (`GET /tracking/{awb}/status`) adalah operasi paling sering dilakukan. Redis menyimpan snapshot status terakhir → response time O(1) tanpa harus scan semua events di MongoDB. TTL 7 hari untuk auto-expiry paket yang sudah delivered/returned.

### Mengapa Dispatch menggunakan PostGIS?
**Keputusan:** PostgreSQL + PostGIS untuk geospatial query.
**Alasan:** Pencarian "kurir terdekat dari titik pickup" membutuhkan query geospatial yang efisien. PostGIS menyediakan index geospatial dan fungsi distance calculation yang optimal.

### Mengapa Settlement dipicu melalui event Kafka?
**Keputusan:** Settlement mendengarkan event `package.delivered` dari Kafka (bukan dipanggil langsung oleh ePOD).
**Alasan:** Decoupling — ePOD tidak perlu tahu tentang settlement. Settlement bisa diproses secara asinkron setelah delivery. Jika settlement gagal, event bisa di-replay dari Kafka.

### Mengapa HTTP Framework bervariasi (Gin vs net/http)?
**Keputusan:** Tidak ada standardisasi framework HTTP — setiap service bebas memilih.
**Alasan:** Setiap service dikembangkan secara independen oleh anggota tim berbeda. Auth, Order, Pricing, ePOD menggunakan Gin. Tracking, Dispatch, Warehouse, Settlement menggunakan `net/http` standard library.
**Dampak:** Tidak ada global middleware yang bisa diterapkan sekaligus — perlu diimplementasikan per service.

---

## 16. Known Bugs & Inconsistencies

> Bagian ini mendokumentasikan inkonsistensi yang **diketahui** agar AI tidak melaporkannya sebagai temuan baru.

| # | Masalah | Lokasi | Status |
|---|---|---|---|
| 1 | `settlement-service/internal/domain/settlement.go` berisi YAML (docker-compose test), bukan Go code | settlement-service | ⚠️ Bug — belum diperbaiki |
| 2 | Kafka topic mismatch: WMS publish `package.arrived` tapi Tracking subscribe ke `package.inbound` | warehouse-service + tracking-service | ⚠️ Perlu diselaraskan |
| 3 | Kafka topic mismatch: WMS publish `manifest.dispatched` tapi Tracking subscribe ke `manifest.arrived` | warehouse-service + tracking-service | ⚠️ Perlu diselaraskan |
| 4 | Auth service handler tidak memiliki `/refresh` endpoint meskipun terdokumentasi | auth-service/handler | ⚠️ Belum diimplementasikan |
| 5 | Port default internal ePOD (`:8080`) berbeda dari docker-compose host mapping (`:8086`) | epod-service/cmd/main.go | ℹ️ Normal (mapping Docker) |
| 6 | Port default internal Pricing (`:8081`) berbeda dari docker-compose host mapping (`:8084`) | pricing-service/cmd/main.go | ℹ️ Normal (mapping Docker) |
| 7 | `order-management-service/internal/domain/` menggunakan `package model` bukan `package domain` | order-management-service | ℹ️ Inkonsistensi naming |
| 8 | TrackingEvent domain struct tidak memiliki field `Latitude`/`Longitude` meskipun ada di payload lama | tracking-service | ℹ️ Field dihapus dari domain |
| 9 | Warehouse dan tracking Kafka implementation adalah placeholder (log only, tidak kirim ke Kafka sungguhan) | warehouse-service, tracking-service | ⚠️ TODO untuk production |
| 10 | CORS tracking service tidak include `localhost:5173` (Vite dev port) | tracking-service/cmd/main.go | ⚠️ Dev issue |

---

## 17. Port Reference

### Docker Host Ports (sesuai docker-compose.yml)

| Service | Container Name | Host Port | Internal Port |
|---|---|---|---|
| PostgreSQL | `buroqet-postgres` | 5432 | 5432 |
| MongoDB | `buroqet-mongo` | 27017 | 27017 |
| Redis | `buroqet-redis` | 6379 | 6379 |
| Zookeeper | `buroqet-zookeeper` | — (internal) | 2181 |
| Kafka | `buroqet-kafka` | 9092 | 9092 (host), 29092 (internal) |
| Auth | `buroqet-auth` | **8080** | 8080 |
| Tracking | `buroqet-tracking` | **8081** | 8080 (APP_PORT) |
| Order | `buroqet-order` | **8082** | 8082 |
| Dispatch | `buroqet-dispatch` | **8083** | 8081 (default internal) |
| Pricing | `buroqet-pricing` | **8084** | 8081 (default internal) |
| Settlement | `buroqet-settlement` | **8085** | 8081 (default internal) |
| ePOD | `buroqet-epod` | **8086** | 8080 (hardcoded) |
| Warehouse | `buroqet-warehouse` | **8087** | 8080 (APP_PORT) |
| Frontend | `buroqet-frontend` | **3000** | 80 (Nginx) |

### Kafka Listener
- **Dari host (development):** `localhost:9092`
- **Dari dalam Docker network (antar container):** `kafka:29092`

---

## 18. Test Data

AWB yang tersedia di MongoDB lokal untuk testing:

| AWB | Route | Status |
|---|---|---|
| `BQ-2024-GPS-001` | Jakarta → Bandung | DELIVERED |
| `BQ-2024-SBY-042` | Surabaya → Malang | OUT_FOR_DELIVERY |
| `BQ-2024-MDN-007` | Medan → Pekanbaru | FAILED |
| `BQ-2024-YGY-099` | Jakarta → Yogyakarta | DELIVERED |

---

## 19. Implementation Status

| Komponen | Backend | Unit Test | Frontend UI | Docker | K8s |
|---|---|---|---|---|---|
| Auth Service | ✅ | ✅ | ⚠️ Login only | ✅ | ✅ |
| Tracking Service | ✅ | ✅ | ✅ Selesai | ✅ | — |
| Order Service | ✅ | ✅ | ❌ placeholder | ✅ | — |
| Dispatch Service | ✅ | ✅ | ❌ placeholder | ✅ | — |
| Pricing Service | ✅ (mock repo) | ✅ | ❌ placeholder | ✅ | — |
| Settlement Service | ✅ (bug domain) | ✅ | ❌ placeholder | ✅ | — |
| ePOD Service | ✅ | ✅ | ❌ placeholder | ✅ | — |
| Warehouse Service | ✅ | ✅ | ❌ placeholder | ✅ | — |
| Frontend Shell | — | — | ✅ Layout+Router | ✅ | — |
| Kafka Integration | ⚠️ Placeholder | — | — | ✅ | — |
| Ingress/K8s | — | — | — | — | ✅ (partial) |
