// settlement-service/internal/kafka/consumer.go
// Kafka consumer untuk menerima event package.delivered dari ePOD Service.

package kafka

import (
	"context"
	"log"

	kafkago "github.com/segmentio/kafka-go"
)

// DeliveryEventHandler adalah fungsi callback yang akan dipanggil consumer
// ketika event package.delivered diterima.
// Parameter: awb = Air Waybill number dari paket yang delivered.
type DeliveryEventHandler func(ctx context.Context, awb string) error

// Consumer membungkus kafka-go Reader dan menyediakan loop konsumsi event.
type Consumer struct {
	reader  *kafkago.Reader
	handler DeliveryEventHandler
}

// NewConsumer membuat consumer baru yang subscribe ke topic package.delivered.
func NewConsumer(broker string, groupID string, handler DeliveryEventHandler) *Consumer {
	reader := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers: []string{broker},
		// topic package.delivered sesuai Kafka Event Catalog di project_context.md section 8
		Topic:   "package.delivered",
		GroupID: groupID,
		// MinBytes dan MaxBytes mengontrol batching; default cocok untuk tugas ini
		MinBytes: 1,    // 1B
		MaxBytes: 10e6, // 10MB
	})
	return &Consumer{reader: reader, handler: handler}
}

// Start memulai loop konsumsi event. Harus dijalankan sebagai goroutine.
// Loop berhenti ketika ctx dibatalkan (graceful shutdown).
func (c *Consumer) Start(ctx context.Context) {
	log.Println("[KAFKA] Settlement consumer dimulai, subscribe ke topic: package.delivered")
	for {
		// ReadMessage akan block sampai ada pesan baru atau ctx dibatalkan
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			// ctx.Done() = shutdown sinyal — berhenti dengan bersih
			if ctx.Err() != nil {
				log.Println("[KAFKA] Settlement consumer berhenti (context cancelled)")
				return
			}
			// Error lain: log dan lanjutkan (eventual consistency pattern)
			log.Printf("[KAFKA] WARNING: gagal membaca pesan dari package.delivered: %v", err)
			continue
		}

		// Payload package.delivered adalah AWB sebagai plain string (sesuai spec section 8)
		// "Payload: AWB string (current implementation: hanya kirim AWB sebagai string message)"
		awb := string(msg.Value)
		if awb == "" {
			log.Printf("[KAFKA] WARNING: menerima pesan dengan payload kosong, key=%s", string(msg.Key))
			continue
		}

		log.Printf("[KAFKA] topic=package.delivered key=%s awb=%s — memproses komisi", string(msg.Key), awb)

		// Panggil handler untuk proses komisi.
		// Handler akan mencatat commission log berdasarkan AWB.
		if err := c.handler(ctx, awb); err != nil {
			// Log WARNING tapi TIDAK gagal — eventual consistency pattern
			// Jika gagal, operator bisa manual trigger via POST /api/v1/commissions
			log.Printf("[KAFKA] WARNING: gagal memproses komisi untuk AWB %s: %v", awb, err)
		}
	}
}

// Close menutup reader Kafka dengan bersih.
func (c *Consumer) Close() error {
	return c.reader.Close()
}
