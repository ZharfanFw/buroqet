FROM golang:alpine AS builder

WORKDIR /app

COPY . .

RUN go mod tidy

# KITA TEGASKAN OS & ARSITEKTURNYA, DAN UBAH NAMA OUTPUT JADI "buroqet-app"
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o buroqet-app ./cmd/main.go

FROM alpine:3.19

WORKDIR /root/

# COPY FILE YANG BARU
COPY --from=builder /app/buroqet-app .

EXPOSE 8080

# JALANKAN FILE YANG BARU
CMD ["./buroqet-app"]