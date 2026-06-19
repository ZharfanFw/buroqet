package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"order-management-service/internal/domain"
	"order-management-service/internal/handler"
	"order-management-service/internal/kafka"
	"order-management-service/internal/repository"
	"order-management-service/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// buildDSN constructs a PostgreSQL DSN from either a full DATABASE_URL
// or individual DB_HOST / DB_PORT / DB_USER / DB_PASSWORD / DB_NAME env vars
// (the format used by docker-compose).
func buildDSN() string {
	if dsn := os.Getenv("DATABASE_URL"); dsn != "" {
		return dsn
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	if host != "" && user != "" {
		if port == "" {
			port = "5432"
		}
		return fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			host, user, password, dbname, port,
		)
	}

	// Local fallback (matches docker-compose postgres service defaults)
	return "host=localhost user=buroqet password=buroqet123 dbname=buroqet port=5432 sslmode=disable"
}

// kafkaBrokerList reads KAFKA_BROKERS (plural) or KAFKA_BROKER (singular), comma-separated.
func kafkaBrokerList() []string {
	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		brokers = os.Getenv("KAFKA_BROKER")
	}
	if brokers == "" {
		return []string{"localhost:9092"}
	}
	result := []string{}
	for _, b := range strings.Split(brokers, ",") {
		if b = strings.TrimSpace(b); b != "" {
			result = append(result, b)
		}
	}
	if len(result) == 0 {
		return []string{"localhost:9092"}
	}
	return result
}

func main() {
	// --- Database setup ---
	dsn := buildDSN()
	log.Printf("Connecting to database...")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Ensure new columns exist manually to bypass any GORM AutoMigrate constraint issues
	db.Exec("ALTER TABLE orders ADD COLUMN IF NOT EXISTS payment_provider VARCHAR(255) DEFAULT ''")
	db.Exec("ALTER TABLE orders ADD COLUMN IF NOT EXISTS payment_code VARCHAR(255) DEFAULT ''")

	// AutoMigrate is enabled to ensure the new columns payment_provider and payment_code exist.
	if err := db.AutoMigrate(&domain.Order{}); err != nil {
		log.Printf("WARNING: AutoMigrate failed: %v", err)
	}

	// --- Kafka setup ---
	kafkaBrokers := kafkaBrokerList()
	log.Printf("Kafka brokers: %v", kafkaBrokers)

	producer := kafka.NewProducer(kafkaBrokers)
	defer producer.Close()

	// --- Wire dependencies ---
	orderRepo := repository.NewOrderRepository(db)

	// --- Kafka consumer setup (Order status automation via tracking.updated events) ---
	consumerGroupID := os.Getenv("KAFKA_GROUP_ID")
	if consumerGroupID == "" {
		consumerGroupID = "oms-order-status-group"
	}
	consumerCtx, consumerCancel := context.WithCancel(context.Background())
	defer consumerCancel()

	consumer := kafka.NewConsumer(kafkaBrokers, consumerGroupID, orderRepo)
	go consumer.Start(consumerCtx)
	defer consumer.Close()

	// PricingClient: calls the real Pricing & Routing Service over HTTP when URL is provided.
	// Falls back to stub implementation when PRICING_SERVICE_URL is empty.
	pricingURL := os.Getenv("PRICING_SERVICE_URL")
	pricingClient := service.NewHTTPPricingClient(pricingURL)
	if pricingURL == "" {
		log.Printf("WARNING: PRICING_SERVICE_URL not set, using stub pricing client")
	}

	orderSvc := service.NewOrderService(orderRepo, pricingClient, producer)
	orderHandler := handler.NewOrderHandler(orderSvc)

	// --- HTTP server setup ---
	r := gin.Default()
	orderHandler.RegisterRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	log.Printf("OMS listening on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
