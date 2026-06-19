# Stage 1: Build the application
FROM golang:1.23-alpine AS builder

# Install git (dibutuhkan jika ada dependency dari private repo)
RUN apk add --no-cache git

WORKDIR /app

# Copy go mod dan sum dulu agar bisa caching layer
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh source code
COPY . .

# Build binary-nya. Kita arahkan ke main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /dispatch-fleet ./cmd/main.go

# Stage 2: Final image yang ringan
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary dari stage builder
COPY --from=builder /dispatch-fleet .

# Ekspos port sesuai yang ada di main.go kamu
EXPOSE 8081

# Jalankan aplikasinya
CMD ["./dispatch-fleet"]