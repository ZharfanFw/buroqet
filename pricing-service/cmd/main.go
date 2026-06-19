package main

import (
	"log"
	"os"

	"pricing-service/internal/handler"
	"pricing-service/internal/repository"
	"pricing-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Load file .env untuk mengambil konfigurasi database
	// Sangat penting untuk menyembunyikan credential dari kode sumber
	if err := godotenv.Load(); err != nil {
		log.Println("Tidak ada file .env ditemukan, menggunakan sistem environment variable")
	}

	// 2. Inisialisasi Gin router
	router := gin.Default()

	// 3. Dapatkan DSN dari environment variable
	// Pastikan di server/container, key DATABASE_URL terisi dengan URL Pooler Supabase (Port 6543)
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL belum diatur di environment variable")
	}

	// 4. Inisialisasi Layered Architecture
	// Repository -> Service -> Handler
	repo := repository.NewPostgresPricingRepository(dsn)
	pricingService := service.NewPricingService(repo)
	pricingHandler := handler.NewPricingHandler(pricingService)

	// 5. Routing API
	router.POST("/pricing/calculate", pricingHandler.CalculatePricing)
	router.POST("/calculate", pricingHandler.CalculatePricing)

	// 6. Jalankan server pada port 8080
	// Port ini harus sama dengan EXPOSE di Dockerfile dan targetPort di K8s YAML
	log.Println("Pricing Service berjalan pada port 8080...")
	router.Run(":8080")
}