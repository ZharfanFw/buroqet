# ================================================
# STAGE 1: BUILD STAGE
# Menggunakan image Go yang besar untuk compile
# ================================================
FROM golang:1.25.6-alpine AS builder

# Install dependency yang dibutuhkan untuk build
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Selama dua file ini tidak berubah, Docker tidak akan re-download dependencies
COPY go.mod go.sum ./

RUN go mod download

# Baru copy source code (sering berubah, jadi ditaruh belakangan)
COPY . .

# Build binary — CGO_ENABLED=0 agar binary bisa jalan di alpine tanpa C library
# -ldflags="-w -s" untuk memperkecil ukuran binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/server ./cmd/main.go

# ================================================
# STAGE 2: RUNTIME STAGE  
# Image final yang SANGAT kecil — hanya berisi binary
# Image builder (golang:1.22-alpine ~300MB) tidak ikut ke sini
# ================================================
FROM alpine:3.19

# Tambahkan certificate untuk HTTPS calls
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy HANYA binary dari stage builder
COPY --from=builder /app/server .

# Expose port service
EXPOSE 8080

# Health check — Kubernetes akan pakai ini
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --quiet --tries=1 --spider http://localhost:8080/health || exit 1

# Jalankan binary
CMD ["./server"]