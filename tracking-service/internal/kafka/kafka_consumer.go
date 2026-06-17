// internal/kafka/kafka_consumer.go
// Implementasi KafkaConsumer untuk Tracking Service.
//
// Tracking service MENDENGARKAN event dari service lain:
//   - Topic "package.inbound"   → dari WMS (paket masuk gudang)
//   - Topic "package.dispatched" → dari Dispatch Service (paket dikirim)
//   - Topic "package.delivered" → dari e-POD Service (paket diterima penerima)
//
// TODO: Implementasi nyata menggunakan github.com/segmentio/kafka-go

package kafka

import (
	"context"
	"fmt"
	"log"
)

// =========================================================
// Topics yang di-subscribe oleh Tracking Service
// =========================================================

const (
	TopicPackageInbound    = "package.inbound"
	TopicPackageDispatched = "package.dispatched"
	TopicPackageDelivered  = "package.delivered"
	TopicManifestArrived   = "manifest.arrived"
)

// TrackingKafkaConsumer adalah implementasi KafkaConsumer untuk tracking service
type TrackingKafkaConsumer struct {
	broker  string
	groupID string
	topics  []string
}

// NewTrackingKafkaConsumer membuat instance consumer baru
func NewTrackingKafkaConsumer(broker string, groupID string) *TrackingKafkaConsumer {
	return &TrackingKafkaConsumer{
		broker:  broker,
		groupID: groupID,
	}
}

// Subscribe mendaftarkan consumer ke satu atau lebih topic Kafka
// TODO: Implementasi nyata menggunakan segmentio/kafka-go:
//
//	reader := kafka.NewReader(kafka.ReaderConfig{
//	    Brokers:  []string{c.broker},
//	    GroupID:  c.groupID,
//	    Topic:    topics[0], // atau gunakan GroupTopics untuk multiple topics
//	    MinBytes: 10e3,
//	    MaxBytes: 10e6,
//	})
func (c *TrackingKafkaConsumer) Subscribe(topics []string) error {
	c.topics = topics
	log.Printf("📡 Kafka consumer subscribe ke topics: %v (broker: %s, group: %s)",
		topics, c.broker, c.groupID)
	// TODO: Implementasi nyata
	return nil // Tidak error saat startup (placeholder)
}

// ReadMessage membaca satu pesan dari Kafka (blocking).
// Mengembalikan topic, key, value, dan error.
// TODO: Implementasi nyata:
//
//	msg, err := reader.ReadMessage(ctx)
//	return msg.Topic, string(msg.Key), msg.Value, err
func (c *TrackingKafkaConsumer) ReadMessage(ctx context.Context) (topic string, key string, value []byte, err error) {
	// Placeholder — blok sampai context di-cancel
	<-ctx.Done()
	return "", "", nil, fmt.Errorf("consumer stopped: %w", ctx.Err())
}

// Close menutup koneksi consumer ke Kafka
// TODO: Implementasi nyata: return reader.Close()
func (c *TrackingKafkaConsumer) Close() error {
	log.Println("🔌 Kafka consumer ditutup")
	return nil
}

// DefaultTopics mengembalikan daftar topic default yang di-listen tracking service
func DefaultTopics() []string {
	return []string{
		TopicPackageInbound,
		TopicPackageDispatched,
		TopicPackageDelivered,
		TopicManifestArrived,
	}
}