package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"settlement-service/internal/domain"
	"settlement-service/internal/handler"
	"settlement-service/internal/kafka"
	"settlement-service/internal/pricing"
	"settlement-service/internal/repository"
	"settlement-service/internal/service"
)

func main() {
	dbHost := getEnv("DB_HOST", "localhost")
	dbUser := getEnv("DB_USER", "settlementuser")
	dbPass := getEnv("DB_PASSWORD", "settlementpassword")
	dbName := getEnv("DB_NAME", "settlement_db")
	dbPort := getEnv("DB_PORT", "5432")
	appPort := getEnv("APP_PORT", "8085")
	pricingURL := getEnv("PRICING_SERVICE_URL", "http://pricing-service:8084")
	kafkaBroker := getEnv("KAFKA_BROKER", "localhost:9092")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPass, dbName, dbPort,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: "settlement.", // Set schema ke settlement
		},
	})
	if err != nil {
		log.Fatalf("Gagal koneksi ke database: %v", err)
	}

	// Buat schema jika belum ada
	db.Exec("CREATE SCHEMA IF NOT EXISTS settlement")

	if err := db.AutoMigrate(&domain.CommissionLog{}); err != nil {
		log.Fatalf("Gagal migrate database: %v", err)
	}

	// Inisialisasi semua layer (Dependency Injection)
	pricingClient := pricing.NewPricingClient(pricingURL)
	repo := repository.NewSettlementRepository(db)
	svc := service.NewSettlementService(repo, pricingClient)
	h := handler.NewSettlementHandler(svc)

	// Setup context yang bisa dibatalkan untuk graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())

	// Inisialisasi dan jalankan Kafka consumer sebagai goroutine
	// Consumer subscribe ke topic "package.delivered" (sesuai Kafka Event Catalog)
	consumer := kafka.NewConsumer(kafkaBroker, "settlement-service-group", func(ctx context.Context, awb string) error {
		// Saat event package.delivered diterima, payload hanya berisi AWB.
		// courier_id dan service_type tidak tersedia dari event ini,
		// sehingga kita gunakan AWB sebagai key identifier, courier_id dari message key (AWB),
		// dan service_type default "REGULER" sebagai fallback.
		// Operator dapat override via manual endpoint POST /api/v1/commissions jika diperlukan.
		log.Printf("[KAFKA] Event package.delivered diterima untuk AWB: %s", awb)
		// Gunakan AWB sebagai courier_id placeholder — ini akan dioverride jika
		// ada manual trigger dengan data lengkap. Di production, butuh lookup ke OMS.
		return svc.ProcessDeliveryCommission(ctx, awb, awb, "REGULER")
	})

	go consumer.Start(ctx)
	log.Printf("[KAFKA] Consumer package.delivered dimulai, broker=%s", kafkaBroker)

	// Setup HTTP routes
	mux := http.NewServeMux()
	mux.HandleFunc("/health", h.HealthCheck)
	mux.HandleFunc("/ready", h.HealthCheck)
	mux.HandleFunc("/api/v1/commissions", h.CommissionsHandler)
	mux.HandleFunc("/api/v1/couriers/", h.GetCourierEarnings)

	// Graceful shutdown: dengarkan sinyal OS (SIGINT, SIGTERM)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Jalankan HTTP server di goroutine terpisah
	serverErr := make(chan error, 1)
	go func() {
		log.Printf("Settlement Service berjalan di port %s", appPort)
		serverErr <- http.ListenAndServe(":"+appPort, mux)
	}()

	// Blok sampai ada sinyal shutdown atau server error
	select {
	case sig := <-sigChan:
		log.Printf("Menerima sinyal %v, memulai graceful shutdown...", sig)
	case err := <-serverErr:
		log.Printf("HTTP server error: %v", err)
	}

	// Batalkan context → hentikan Kafka consumer
	cancel()
	if err := consumer.Close(); err != nil {
		log.Printf("WARNING: gagal menutup Kafka consumer: %v", err)
	}
	log.Println("Settlement Service berhenti.")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}