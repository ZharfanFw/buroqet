package kafka

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

// EPODEvent adalah struktur data yang akan dikirim ke Kafka topic
type EPODEvent struct {
	EventName string    `json:"event_name"` // misal: "EPOD_UPLOADED" atau "EPOD_VERIFIED"
	EPODID    string    `json:"epod_id"`
	AWB       string    `json:"awb"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

// ProducerInterface mendefinisikan fungsi yang harus dimiliki oleh producer
type ProducerInterface interface {
	PublishEvent(ctx context.Context, topic string, event EPODEvent) error
	Close() error
}

// Producer adalah implementasi konkret dari ProducerInterface
type Producer struct {
	writer *kafka.Writer
}

// NewProducer membuat instance baru untuk Kafka Producer
// Fungsi ini otomatis membaca KAFKA_BROKERS dari .env
func NewProducer() ProducerInterface {
	brokersEnv := os.Getenv("KAFKA_BROKERS")
	if brokersEnv == "" {
		brokersEnv = "localhost:9092" // default jika tidak diatur
	}

	brokers := strings.Split(brokersEnv, ",")

	// Inisialisasi Kafka Writer
	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Balancer:     &kafka.LeastBytes{},
		MaxAttempts:  3,
		WriteTimeout: 10 * time.Second,
		RequiredAcks: kafka.RequireAll, // Menjamin pesan benar-benar tersimpan di cluster
	}

	log.Printf("Kafka Producer berhasil diinisialisasi ke brokers: %v", brokers)

	return &Producer{
		writer: writer,
	}
}

// PublishEvent mengirimkan data event ke Kafka topic tertentu
func (p *Producer) PublishEvent(ctx context.Context, topic string, event EPODEvent) error {
	// Ubah struct event menjadi JSON byte
	payload, err := json.Marshal(event)
	if err != nil {
		log.Printf("❌ Kafka: Gagal marshal event ke JSON: %v", err)
		return err
	}

	// Kirim pesan ke Kafka
	err = p.writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Key:   []byte(event.AWB), // Menggunakan AWB sebagai key agar urutan per resi tetap konsisten
		Value: payload,
	})

	if err != nil {
		log.Printf("❌ Kafka: Gagal mengirim pesan ke topic %s: %v", topic, err)
		return err
	}

	log.Printf("🚀 Kafka: Berhasil mengirim event [%s] untuk AWB %s ke topic %s", event.EventName, event.AWB, topic)
	return nil
}

// Close menutup koneksi writer ke Kafka cluster (dipanggil via defer di main.go)
func (p *Producer) Close() error {
	log.Println("Mematikan Kafka Producer secara aman...")
	return p.writer.Close()
}