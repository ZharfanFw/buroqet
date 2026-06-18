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

// buildDatabaseURL membangun connection string dari env var individual.
// Urutan prioritas:
//  1. DATABASE_URL (jika di-set langsung, dipakai apa adanya — berguna untuk local dev)
//  2. DB_HOST + DB_PORT + DB_USER + DB_PASSWORD + DB_NAME (diinjek oleh K8s ConfigMap/Secret)
func buildDatabaseURL() string {
	// Prioritas 1: DATABASE_URL langsung (local dev / docker-compose)
	if url := os.Getenv("DATABASE_URL"); url != "" {
		return url
	}

	// Prioritas 2: env var individual yang diinjek K8s
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Fallback ke nilai default untuk keperluan local dev tanpa Docker
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "5433" // port default docker-compose.test.yml
	}
	if user == "" {
		user = "user"
	}
	if password == "" {
		password = "pass"
	}
	if dbname == "" {
		dbname = "dispatch_db"
	}

	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)
}

func main() {
	// 1. Setup Database Connection (PostgreSQL + PostGIS)
	dbURL := buildDatabaseURL()

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Gagal inisialisasi database: %v", err)
	}
	defer db.Close()

	// Verifikasi koneksi saat startup (fail fast jika DB tidak reachable)
	if err := db.Ping(); err != nil {
		log.Fatalf("Gagal konek ke database: %v", err)
	}
	log.Println("✅ Database connected.")

	// 2. Inisialisasi Layer (Clean Architecture)
	// Implementasi Repository menggunakan Postgres + PostGIS
	fleetRepo := repository.NewPostgresFleetRepository(db)
	// Business Logic Service
	dispatchSvc := service.NewDispatchService(fleetRepo)
	// HTTP Handler
	dispatchHandler := handler.NewDispatchHandler(dispatchSvc)

	// 3. Routing
	mux := http.NewServeMux()

	// Endpoint bisnis
	mux.HandleFunc("/v1/dispatch/assign", dispatchHandler.Assign)

	// Health & readiness probe untuk Kubernetes liveness/readinessProbe
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// 4. Start Server
	// K8s deployment.yaml menginjek APP_PORT=8081; fallback ke 8081 jika tidak ada
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = os.Getenv("PORT") // backward compat
	}
	if port == "" {
		port = "8081"
	}

	log.Printf("Dispatch-Fleet Service running on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}