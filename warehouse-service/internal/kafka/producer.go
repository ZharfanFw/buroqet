// warehouse-service/internal/kafka/producer.go
// Implementasi nyata Kafka producer menggunakan kafka-go.
// Menggantikan placeholder log-only yang ada sebelumnya (Bug #9 di project_context.md).

package kafka

import (
	"context"
	"log"

	kafkago "github.com/segmentio/kafka-go"
)

// kafkaProducer adalah implementasi nyata yang mengirim event ke Kafka broker.
// Mengimplementasikan interface domain.KafkaProducer.
type kafkaProducer struct {
	broker string
	writer *kafkago.Writer
}

// NewKafkaProducer membuat producer baru yang terhubung ke broker.
// Writer dikonfigurasi dengan AllowAutoTopicCreation=true agar tidak perlu
// membuat topic secara manual di development environment.
func NewKafkaProducer(broker string) *kafkaProducer {
	writer := &kafkago.Writer{
		Addr:                   kafkago.TCP(broker),
		Balancer:               &kafkago.LeastBytes{},
		RequiredAcks:           kafkago.RequireOne,
		AllowAutoTopicCreation: true, // Otomatis buat topic jika belum ada (dev-friendly)
	}
	log.Printf("[KAFKA] Producer diinisialisasi, broker=%s", broker)
	return &kafkaProducer{broker: broker, writer: writer}
}

// PublishEvent mengirim event ke topic Kafka yang ditentukan.
// key digunakan sebagai Kafka message key (biasanya AWB) untuk menjamin ordering per AWB.
// Jika pengiriman gagal, error dikembalikan ke caller untuk di-handle (log WARNING + lanjut).
func (k *kafkaProducer) PublishEvent(ctx context.Context, topic string, key string, value []byte) error {
	err := k.writer.WriteMessages(ctx, kafkago.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: value,
	})
	if err != nil {
		log.Printf("[KAFKA] ERROR: gagal publish ke topic=%s key=%s: %v", topic, key, err)
		return err
	}
	log.Printf("[KAFKA] topic=%s key=%s value=%s", topic, key, string(value))
	return nil
}

// Close menutup writer Kafka dengan bersih — panggil saat shutdown service.
func (k *kafkaProducer) Close() error {
	return k.writer.Close()
}