package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"warehouse-service/internal/domain"
	"warehouse-service/internal/handler"
	"warehouse-service/internal/kafka"
	"warehouse-service/internal/repository"
	"warehouse-service/internal/service"
)

func main() {
	dbHost := getEnv("DB_HOST", "localhost")
	dbUser := getEnv("DB_USER", "wmsuser")
	dbPass := getEnv("DB_PASSWORD", "wmspassword")
	dbName := getEnv("DB_NAME", "wms_db")
	dbPort := getEnv("DB_PORT", "5432")
	appPort := getEnv("APP_PORT", "8080")
	kafkaBroker := getEnv("KAFKA_BROKER", "localhost:9092")

	// Setup koneksi database
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPass, dbName, dbPort,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal koneksi ke database: %v", err)
	}

	// Auto migrate tabel (manifest dulu karena packages punya FK ke manifests)
	if err := db.AutoMigrate(&domain.Manifest{}, &domain.Package{}); err != nil {
		log.Fatalf("Gagal migrate database: %v", err)
	}

	// Inisialisasi semua layer
	kafkaProducer := kafka.NewKafkaProducer(kafkaBroker)
	repo := repository.NewWarehouseRepository(db)
	svc := service.NewWarehouseService(repo, kafkaProducer)
	h := handler.NewWarehouseHandler(svc)

	// Setup routes
	mux := http.NewServeMux()
	mux.HandleFunc("/health", h.HealthCheck)
	mux.HandleFunc("/ready", h.HealthCheck)
	mux.HandleFunc("/api/v1/inbound", h.ProcessInbound)
	mux.HandleFunc("/api/v1/dispatch", h.DispatchManifest)

	log.Printf("Warehouse Service berjalan di port %s", appPort)
	if err := http.ListenAndServe(":"+appPort, mux); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

