// internal/kafka/kafka_consumer.go
// Implementasi KafkaConsumer untuk Tracking Service menggunakan segmentio/kafka-go.
//
// Tracking service MENDENGARKAN event dari service lain:
//   - Topic "package.inbound"    → dari WMS (paket masuk gudang)
//   - Topic "package.dispatched" → dari Dispatch Service (paket dikirim)
//   - Topic "package.delivered"  → dari e-POD Service (paket diterima penerima)
//   - Topic "manifest.arrived"   → dari Warehouse (manifest tiba di hub)

package kafka

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	kafkago "github.com/segmentio/kafka-go"
)

// Topics yang dikirim oleh masing-masing service (nama sesuai source code producer):
//
// WMS (warehouse-service) → package.arrived, manifest.dispatched
// Dispatch Fleet Service  → package.dispatched
// ePOD Service            → package.delivered
//
// CATATAN INKONSISTENSI (Bug #2 & #3 dari PROJECT_CONTEXT.md):
// Dokumentasi lama menyebut "package.inbound" dan "manifest.arrived",
// namun WMS source code publish ke "package.arrived" dan "manifest.dispatched".
// Solusi: subscribe ke KEDUA nama (lama + baru) agar tidak ada event yang terlewat.

const (
	// Dari WMS — nama yang benar sesuai warehouse-service source code
	TopicPackageArrived     = "package.arrived"     // paket scan masuk gudang
	TopicManifestDispatched = "manifest.dispatched" // manifest kumpulan paket dikirim

	// Dari Dispatch Fleet Service
	TopicPackageDispatched = "package.dispatched" // kurir berangkat antar paket

	// Dari ePOD Service
	TopicPackageDelivered = "package.delivered" // paket diterima penerima

	// Alias lama (nama dari dokumentasi lama) — dipertahankan untuk backward compat
	TopicPackageInbound  = "package.inbound"  // alias lama → sama dgn package.arrived
	TopicManifestArrived = "manifest.arrived" // alias lama → sama dgn manifest.dispatched
)

// TrackingKafkaConsumer adalah implementasi KafkaConsumer menggunakan segmentio/kafka-go
type TrackingKafkaConsumer struct {
	brokers []string
	groupID string
	reader  *kafkago.Reader
}

// NewTrackingKafkaConsumer membuat instance consumer baru
func NewTrackingKafkaConsumer(broker string, groupID string) *TrackingKafkaConsumer {
	return &TrackingKafkaConsumer{
		brokers: strings.Split(broker, ","),
		groupID: groupID,
	}
}

// Subscribe mendaftarkan consumer ke satu atau lebih topic Kafka.
// segmentio/kafka-go menggunakan GroupTopics untuk multiple topic subscription.
func (c *TrackingKafkaConsumer) Subscribe(topics []string) error {
	if len(topics) == 0 {
		return fmt.Errorf("minimal satu topic harus di-subscribe")
	}

	c.reader = kafkago.NewReader(kafkago.ReaderConfig{
		Brokers:     c.brokers,
		GroupID:     c.groupID,
		GroupTopics: topics,

		// Batch settings untuk efisiensi
		MinBytes:    10e3, // 10KB
		MaxBytes:    10e6, // 10MB
		MaxWait:     1 * time.Second,

		// Retry settings
		MaxAttempts: 3,

		// Start dari earliest jika group baru
		StartOffset: kafkago.FirstOffset,

		// Logger (opsional, bisa dimatikan di production)
		Logger: kafkago.LoggerFunc(func(msg string, args ...interface{}) {
			log.Printf("[kafka-consumer] "+msg, args...)
		}),
		ErrorLogger: kafkago.LoggerFunc(func(msg string, args ...interface{}) {
			log.Printf("[kafka-consumer ERROR] "+msg, args...)
		}),
	})

	log.Printf("📡 Kafka consumer subscribe ke topics: %v (broker: %v, group: %s)",
		topics, c.brokers, c.groupID)
	return nil
}

// ReadMessage membaca satu pesan dari Kafka (blocking).
// Mengembalikan topic, key, value, dan error.
// Akan block sampai ada message atau context di-cancel.
func (c *TrackingKafkaConsumer) ReadMessage(ctx context.Context) (topic string, key string, value []byte, err error) {
	if c.reader == nil {
		return "", "", nil, fmt.Errorf("consumer belum di-subscribe, panggil Subscribe() dulu")
	}

	msg, err := c.reader.ReadMessage(ctx)
	if err != nil {
		return "", "", nil, err
	}

	return msg.Topic, string(msg.Key), msg.Value, nil
}

// Close menutup koneksi consumer ke Kafka dengan bersih
func (c *TrackingKafkaConsumer) Close() error {
	if c.reader != nil {
		log.Println("🔌 Kafka consumer ditutup")
		return c.reader.Close()
	}
	return nil
}

// DefaultTopics mengembalikan daftar semua topic yang di-listen tracking service.
// Include KEDUA nama (lama + baru) untuk menangani ketidaksesuaian
// antara dokumentasi dan implementasi WMS producer.
func DefaultTopics() []string {
	return []string{
		// Dari WMS (nama benar sesuai source code)
		TopicPackageArrived,
		TopicManifestDispatched,
		// Dari Dispatch Fleet Service
		TopicPackageDispatched,
		// Dari ePOD Service
		TopicPackageDelivered,
		// Alias lama (backward compat)
		TopicPackageInbound,
		TopicManifestArrived,
	}
}