package main

import (
	"log"
	"os"

	"order-management-service/internal/handler"
	"order-management-service/internal/kafka"
	"order-management-service/internal/domain"
	"order-management-service/internal/repository"
	"order-management-service/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// --- Database setup ---
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=oms_db port=5432 sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Auto-migrate will create/update the orders table based on the model struct
	if err := db.AutoMigrate(&model.Order{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	// --- Kafka producer setup ---
	kafkaBrokers := []string{os.Getenv("KAFKA_BROKER")}
	if kafkaBrokers[0] == "" {
		kafkaBrokers = []string{"localhost:9092"}
	}
	producer := kafka.NewProducer(kafkaBrokers)
	defer producer.Close()

	// --- Wire dependencies ---
	orderRepo := repository.NewOrderRepository(db)

	// PricingClient: in a real setup this would be an HTTP client pointing at Pricing Service
	// For now we use a stub so the app compiles and runs without the external service
	pricingClient := service.NewHTTPPricingClient(os.Getenv("PRICING_SERVICE_URL"))

	orderSvc := service.NewOrderService(orderRepo, pricingClient, producer)
	orderHandler := handler.NewOrderHandler(orderSvc)

	// --- HTTP server setup ---
	r := gin.Default()
	orderHandler.RegisterRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("OMS listening on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}