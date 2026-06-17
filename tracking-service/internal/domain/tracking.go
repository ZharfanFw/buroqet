// internal/domain/tracking.go
// Domain layer: berisi model data, konstanta status, dan kontrak interface.
// Layer ini TIDAK memiliki ketergantungan pada library eksternal —
// hanya menggunakan tipe bawaan Go dan standard library.

package domain

import (
	"context"
	"time"
)

// =========================================================
// STATUS CONSTANTS
// Mendefinisikan semua status paket yang valid di sistem
// =========================================================

const (
	StatusInbound    = "INBOUND"     // Paket diterima di gudang asal
	StatusOnTransit  = "ON_TRANSIT"  // Paket sedang dalam perjalanan antar hub
	StatusAtHub      = "AT_HUB"      // Paket tiba di hub transit
	StatusOutForDelivery = "OUT_FOR_DELIVERY" // Paket sedang diantar ke penerima
	StatusDelivered  = "DELIVERED"   // Paket sudah diterima oleh penerima (e-POD)
	StatusFailed     = "FAILED"      // Pengiriman gagal (tidak ada di rumah, dsb)
	StatusReturned   = "RETURNED"    // Paket dikembalikan ke pengirim
)

// =========================================================
// TRACKING EVENT — Dokumen MongoDB
// Setiap event adalah satu baris riwayat perjalanan paket.
// Koleksi ini akan terus bertambah (append-only / log pattern).
// =========================================================

// TrackingEvent merepresentasikan satu kejadian/pergerakan paket
type TrackingEvent struct {
	ID          string    `bson:"_id,omitempty"    json:"id"`
	AWB         string    `bson:"awb"              json:"awb"`               // Nomor Airway Bill
	Status      string    `bson:"status"           json:"status"`            // Status baru setelah event
	Location    string    `bson:"location"         json:"location"`          // Nama hub/lokasi
	HubID       string    `bson:"hub_id"           json:"hub_id"`            // ID hub
	Description string    `bson:"description"      json:"description"`       // Keterangan tambahan
	Timestamp   time.Time `bson:"timestamp"        json:"timestamp"`         // Waktu event terjadi
	CreatedAt   time.Time `bson:"created_at"       json:"created_at"`        // Waktu dicatat di sistem
	Source      string    `bson:"source"           json:"source"`            // Service pengirim event (WMS, Dispatch, e-POD)
}

// TrackingHistory adalah respons API: kumpulan events terurut berdasarkan waktu
type TrackingHistory struct {
	AWB    string          `json:"awb"`
	Events []TrackingEvent `json:"events"`
	Total  int             `json:"total"`
}

// =========================================================
// TRACKING STATUS — Data di Redis
// Menyimpan status TERAKHIR paket agar query status cepat
// tanpa harus scan semua event di MongoDB
// =========================================================

// TrackingStatus adalah snapshot status terakhir paket (disimpan di Redis)
type TrackingStatus struct {
	AWB          string    `json:"awb"`
	CurrentStatus string   `json:"current_status"`
	LastLocation  string   `json:"last_location"`
	LastUpdated   time.Time `json:"last_updated"`
}

// =========================================================
// KAFKA EVENT — Payload yang diterima dari service lain
// Tracking service mendengarkan beberapa topic Kafka
// =========================================================

// KafkaTrackingPayload adalah format pesan Kafka yang diterima
type KafkaTrackingPayload struct {
	AWB       string    `json:"awb"`
	Status    string    `json:"status"`
	HubID     string    `json:"hub_id"`
	Location  string    `json:"location"`
	Timestamp time.Time `json:"timestamp"`
	Source    string    `json:"source"`
}

// =========================================================
// REQUEST / RESPONSE DTO (Data Transfer Object)
// =========================================================

// AddTrackingEventRequest adalah payload untuk menambahkan event baru
// (digunakan baik dari HTTP API maupun dari konsumsi Kafka)
type AddTrackingEventRequest struct {
	AWB       string    `json:"awb"`
	Status    string    `json:"status"`
	HubID     string    `json:"hub_id"`
	Location  string    `json:"location"`
	Timestamp time.Time `json:"timestamp"`
	Source    string    `json:"source"`
}

// =========================================================
// INTERFACES — Kontrak yang harus dipenuhi oleh implementasi
// Interface inilah yang akan di-mock saat unit test
// =========================================================

// TrackingRepository adalah kontrak akses ke MongoDB (primary store)
type TrackingRepository interface {
	// InsertEvent menyimpan satu event tracking ke MongoDB
	InsertEvent(ctx context.Context, event *TrackingEvent) error

	// GetEventsByAWB mengambil semua riwayat perjalanan satu AWB,
	// terurut berdasarkan timestamp ascending
	GetEventsByAWB(ctx context.Context, awb string) ([]TrackingEvent, error)

	// GetLatestEventByAWB mengambil event terbaru satu AWB
	GetLatestEventByAWB(ctx context.Context, awb string) (*TrackingEvent, error)
}

// TrackingCache adalah kontrak akses ke Redis (status terakhir paket)
type TrackingCache interface {
	// SetStatus menyimpan status terakhir paket ke Redis dengan TTL
	SetStatus(ctx context.Context, status *TrackingStatus) error

	// GetStatus mengambil status terakhir paket dari Redis
	// Mengembalikan nil, nil jika tidak ada (cache miss)
	GetStatus(ctx context.Context, awb string) (*TrackingStatus, error)

	// DeleteStatus menghapus status dari cache (saat paket delivered/returned)
	DeleteStatus(ctx context.Context, awb string) error
}

// KafkaConsumer adalah kontrak untuk Kafka consumer
// Interface ini memungkinkan mock saat testing
type KafkaConsumer interface {
	// Subscribe mendaftarkan consumer ke satu atau lebih topic
	Subscribe(topics []string) error

	// ReadMessage membaca satu pesan dari Kafka (blocking)
	ReadMessage(ctx context.Context) (topic string, key string, value []byte, err error)

	// Close menutup koneksi consumer
	Close() error
}

// KafkaProducer adalah kontrak untuk publish event
// (jika tracking service perlu mempublish ulang event)
type KafkaProducer interface {
	PublishEvent(ctx context.Context, topic string, key string, value []byte) error
}