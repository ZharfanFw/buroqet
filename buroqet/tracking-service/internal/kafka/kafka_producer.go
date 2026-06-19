// internal/kafka/kafka_producer.go
// Implementasi KafkaProducer menggunakan segmentio/kafka-go.
// Publish tracking event ke downstream services.

package kafka

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	kafkago "github.com/segmentio/kafka-go"
)

// TrackingKafkaProducer adalah implementasi KafkaProducer untuk tracking service
type TrackingKafkaProducer struct {
	brokers []string
}

// NewTrackingKafkaProducer membuat instance producer baru
func NewTrackingKafkaProducer(broker string) *TrackingKafkaProducer {
	brokers := strings.Split(broker, ",")
	return &TrackingKafkaProducer{brokers: brokers}
}

// PublishEvent mempublish event ke Kafka topic menggunakan segmentio/kafka-go.
// Setiap publish membuat writer baru dan langsung di-close setelahnya (stateless).
// Ini lebih sederhana untuk service dengan volume event yang tidak terlalu tinggi.
func (p *TrackingKafkaProducer) PublishEvent(ctx context.Context, topic string, key string, value []byte) error {
	writer := &kafkago.Writer{
		Addr:         kafkago.TCP(p.brokers...),
		Topic:        topic,
		Balancer:     &kafkago.LeastBytes{},
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}
	defer writer.Close()

	err := writer.WriteMessages(ctx, kafkago.Message{
		Key:   []byte(key),
		Value: value,
	})
	if err != nil {
		return fmt.Errorf("kafka publish ke topic '%s' gagal: %w", topic, err)
	}

	log.Printf("📤 Kafka: published event ke topic '%s', key='%s'", topic, key)
	return nil
}

// EnsureTopics membuat topic Kafka jika belum ada.
// Dipanggil sekali saat startup untuk memastikan semua topic tersedia.
func EnsureTopics(broker string, topics []string) error {
	conn, err := kafkago.Dial("tcp", broker)
	if err != nil {
		return fmt.Errorf("gagal koneksi ke Kafka broker %s: %w", broker, err)
	}
	defer conn.Close()

	// Cari controller broker
	controller, err := conn.Controller()
	if err != nil {
		return fmt.Errorf("gagal dapat controller Kafka: %w", err)
	}

	controllerConn, err := kafkago.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return fmt.Errorf("gagal koneksi ke controller Kafka: %w", err)
	}
	defer controllerConn.Close()

	topicConfigs := make([]kafkago.TopicConfig, 0, len(topics))
	for _, t := range topics {
		topicConfigs = append(topicConfigs, kafkago.TopicConfig{
			Topic:             t,
			NumPartitions:     1,
			ReplicationFactor: 1,
		})
	}

	if err := controllerConn.CreateTopics(topicConfigs...); err != nil {
		// Topic sudah ada bukan error fatal
		log.Printf("WARNING: EnsureTopics: %v", err)
	}

	log.Printf("✅ Kafka topics ensured: %v", topics)
	return nil
}