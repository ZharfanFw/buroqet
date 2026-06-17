FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o server ./cmd/main.go

FROM alpine:3.19

WORKDIR /root/

COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]