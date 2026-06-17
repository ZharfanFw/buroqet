package main
import (
	"fmt"
	"log"
	"net/http"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"settlement-service/internal/domain"
	"settlement-service/internal/handler"
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
	appPort := getEnv("APP_PORT", "8081")
	pricingURL := getEnv("PRICING_SERVICE_URL", "http://localhost:8082")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPass, dbName, dbPort,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal koneksi ke database: %v", err)
	}

	if err := db.AutoMigrate(&domain.CommissionLog{}); err != nil {
		log.Fatalf("Gagal migrate database: %v", err)
	}

	pricingClient := pricing.NewPricingClient(pricingURL)
	repo := repository.NewSettlementRepository(db)
	svc := service.NewSettlementService(repo, pricingClient)
	h := handler.NewSettlementHandler(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", h.HealthCheck)
	mux.HandleFunc("/ready", h.HealthCheck)
	mux.HandleFunc("/api/v1/commissions", h.ProcessCommission)
	mux.HandleFunc("/api/v1/couriers/", h.GetCourierEarnings)

	log.Printf("Settlement Service berjalan di port %s", appPort)
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