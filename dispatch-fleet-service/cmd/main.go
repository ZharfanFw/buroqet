package main

import (
	"context"
	"database/sql"
	"dispatch-fleet/internal/handler"
	"dispatch-fleet/internal/kafka"
	"dispatch-fleet/internal/repository"
	"dispatch-fleet/internal/service"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq" // Driver PostgreSQL
)

func main() {
	// 1. Setup Database Connection (PostgreSQL)
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	if dbHost == "" { dbHost = "localhost" }
	if dbUser == "" { dbUser = "postgres" }
	if dbPass == "" { dbPass = "password" }
	if dbName == "" { dbName = "dispatch_db" }
	if dbPort == "" { dbPort = "5432" }

	dbSslMode := os.Getenv("DB_SSLMODE")
	if dbSslMode == "" {
		dbSslMode = "require"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		dbHost, dbUser, dbPass, dbName, dbPort, dbSslMode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Gagal inisialisasi database: %v", err)
	}
	defer db.Close()

	// Inisialisasi PostGIS dan Schema Tabel
	initDB(db)

	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		kafkaBrokers = "localhost:9092"
	}

	// 2. Inisialisasi Layer (Clean Architecture)
	fleetRepo := repository.NewPostgresFleetRepository(db)
	dispatchSvc := service.NewDispatchService(fleetRepo)

	// Kafka Producer
	kafkaProducer := kafka.NewProducer([]string{kafkaBrokers})
	defer kafkaProducer.Close()

	// Kafka Consumer (run in background)
	kafkaConsumer := kafka.NewConsumer([]string{kafkaBrokers}, dispatchSvc)
	go kafkaConsumer.Start(context.Background())
	defer kafkaConsumer.Close()

	dispatchHandler := handler.NewDispatchHandler(dispatchSvc, kafkaProducer)

	// 3. Routing
	mux := http.NewServeMux()

	// CORS Middleware wrap
	corsHandler := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}

	mux.HandleFunc("/v1/dispatch/assign", dispatchHandler.Assign)
	mux.HandleFunc("/v1/dispatch/start-delivery", dispatchHandler.StartDelivery)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"UP"}`))
	})

	// 4. Start Server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8083" // Sesuai project_context.md
	}

	fmt.Printf("Dispatch-Fleet Service running on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, corsHandler(mux)); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func initDB(db *sql.DB) {
	queries := []string{
		`CREATE EXTENSION IF NOT EXISTS postgis;`,
		`CREATE TABLE IF NOT EXISTS couriers (
			id VARCHAR(50) PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			current_location GEOGRAPHY(Point, 4326),
			status VARCHAR(20) NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		);`,
	}
	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			log.Printf("Warning executing init query: %v", err)
		}
	}
}