FROM golang:1.25-alpine AS builder

WORKDIR /app

# Docker Cache untuk dependency
COPY go.mod go.sum ./
RUN go mod download

# Baru copy source code
COPY . .

# Build dengan mematikan CGO & optimasi ukuran binary (-ldflags)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/server ./cmd/main.go

FROM alpine:3.19

WORKDIR /root/
COPY --from=builder /app/server .

EXPOSE 8080
CMD ["./server"]