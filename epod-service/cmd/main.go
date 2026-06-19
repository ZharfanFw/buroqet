package main

import (
	"log"
	"os"


	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"epod-service/internal/handler"
	"epod-service/internal/kafka"
	"epod-service/internal/repository"
	"epod-service/internal/service"
	"epod-service/internal/storage"
	
)

func main() {
	// 1. Load file .env untuk konfigurasi lokal (DB, Kafka broker, dll)
	if err := godotenv.Load(); err != nil {
        log.Println("Tidak ada file .env ditemukan, menggunakan sistem environment variable")
    }

	dbURL := os.Getenv("DATABASE_URL")
    if dbURL == "" {
        log.Fatal("DATABASE_URL belum diatur di environment variable")
    }

	// 2. Inisialisasi Gin router
	router := gin.Default()

	router.Use(func(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Content-Type")
	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}
	c.Next()
})

	// 3. Serve folder ./uploads secara statis agar foto/tanda tangan
	// bisa diakses langsung dari frontend (mis. http://host/uploads/photos/xxx.jpg)
	router.Static("/uploads", "./uploads")

	// 4. Dapatkan DSN dari environment variable.
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL belum diatur di environment variable")
	}

	// 5. Inisialisasi Layered Architecture
	// Repository -> Service -> Handler
	repo := repository.NewPostgresEPODRepository(dsn)
	fileStorage := storage.NewLocalStorage()
	kafkaProducer := kafka.NewProducer()
	defer kafkaProducer.Close()

	epodService := service.NewEPODService(repo, fileStorage, kafkaProducer)
	epodHandler := handler.NewEPODHandler(epodService)

	// 6. Routing API
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})
	router.POST("/upload", epodHandler.Upload)
	router.GET("/epod", epodHandler.List)
	router.GET("/epod/:id", epodHandler.GetByID)
	router.GET("/epod/awb/:awb", epodHandler.GetByAWB)
	router.PATCH("/epod/:id/verify", epodHandler.Verify)

	// 7. Jalankan server pada port 8080
	log.Println("ePOD Service berjalan pada port 8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}