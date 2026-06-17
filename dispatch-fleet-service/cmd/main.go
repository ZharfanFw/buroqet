package main

import (
	"database/sql"
	"dispatch-fleet/internal/handler"
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
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		// Default port 5433 sesuai setup docker-compose.test.yml
		dbURL = "postgres://user:pass@localhost:5433/dispatch_db?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Gagal inisialisasi database: %v", err)
	}
	defer db.Close()

	// 2. Inisialisasi Layer (Clean Architecture)
	// Implementasi Repository menggunakan Postgres + PostGIS
fleetRepo := repository.NewPostgresFleetRepository(db)
	// Business Logic Service
	dispatchSvc := service.NewDispatchService(fleetRepo)

	// HTTP Handler
	dispatchHandler := handler.NewDispatchHandler(dispatchSvc)

	// 3. Routing
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/dispatch/assign", dispatchHandler.Assign)

	// 4. Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081" // Port berbeda dari order-service agar tidak bentrok
	}

	fmt.Printf("Dispatch-Fleet Service running on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}