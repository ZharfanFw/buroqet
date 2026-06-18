package repository

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Struktur tabel di database Supabase Anda
type Tariff struct {
    OriginPostalCode      string
    DestinationPostalCode string
    ServiceType           string
    BaseRate              float64
}

type PostgresPricingRepository struct {
	db *gorm.DB
}

// Tambahkan fungsi ini agar GORM tahu tabel ini ada di skema khusus
func (Tariff) TableName() string {
    return "pricing_schema.tariffs" // Sesuaikan dengan nama_schema.nama_tabel Anda
}

// Inisialisasi koneksi ke Supabase
func NewPostgresPricingRepository(dsn string) *PostgresPricingRepository {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal terhubung ke Supabase: %v", err)
	}

	return &PostgresPricingRepository{db: db}
}

// Implementasi interface dari layer service
func (r *PostgresPricingRepository) GetBaseRate(origin, destination, serviceType string) float64 {
	var tariff Tariff

	// Query ke database Supabase
	result := r.db.Where(
		"origin_postal_code = ? AND destination_postal_code = ? AND service_type = ?", 
		origin, destination, serviceType,
	).First(&tariff)

	if result.Error != nil {
		// Jika rute tidak ditemukan di DB, kita kembalikan tarif default sementara
		log.Printf("Tarif tidak ditemukan: %v", result.Error)
		return 10000.0 
	}

	return tariff.BaseRate
}