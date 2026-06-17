package main

import (
	"auth-service/internal/domain"
	"auth-service/internal/handler"
	"auth-service/internal/repository"
	"auth-service/internal/service"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Ambil konfigurasi dari Environment Variables (disuntikkan oleh Kubernetes)
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	// Fallback nilai default jika variabel environment tidak ditemukan (untuk run lokal)
	if dbHost == "" { dbHost = "localhost" }
	if dbUser == "" { dbUser = "postgres" }
	if dbPass == "" { dbPass = "password" }
	if dbName == "" { dbName = "auth_db" }
	if dbPort == "" { dbPort = "5432" }

	// Susun DSN secara dinamis
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPass, dbName, dbPort)

	// Inisialisasi koneksi Database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal koneksi database ke %s: %v", dbHost, err)
	}

	// Konfigurasi Connection Pool untuk optimasi RPS
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetMaxOpenConns(100)

	// Jalankan Auto Migrate untuk table User
	db.AutoMigrate(&domain.User{})

	// Setup Dependency Injection
	userRepo := repository.NewUserPostgres(db)
	authServ := service.NewAuthService(userRepo)
	authHand := handler.NewAuthHandler(authServ)

	// Inisialisasi Gin Router
	r := gin.Default()

	// Definisi Routes
	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/register", authHand.Register)
		authRoutes.POST("/login", authHand.Login)
	}

	// Jalankan server pada port 8080
	log.Printf("Starting server on %s:8080", dbHost)
	r.Run(":8080")
}