package main

import (
	"auth-service/internal/domain"
	"auth-service/internal/handler"
	"auth-service/internal/repository"
	"auth-service/internal/service"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Ambil konfigurasi dari Environment Variables
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	if dbHost == "" { dbHost = "localhost" }
	if dbUser == "" { dbUser = "postgres" }
	if dbPass == "" { dbPass = "password" }
	if dbName == "" { dbName = "buroqet" }
	if dbPort == "" { dbPort = "5432" }

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPass, dbName, dbPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal koneksi database ke %s: %v", dbHost, err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetMaxOpenConns(100)

	// Auto Migrate untuk table User
	db.AutoMigrate(&domain.User{})

	userRepo := repository.NewUserPostgres(db)
	authServ := service.NewAuthService(userRepo)
	authHand := handler.NewAuthHandler(authServ)

	r := gin.Default()

	// CORS Middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Liveness probe
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP", "service": "auth-service"})
	})

	// Definisi Routes sesuai project_context.md
	authRoutes := r.Group("/api/v1/auth")
	{
		authRoutes.POST("/register", authHand.Register)
		authRoutes.POST("/login", authHand.Login)
		authRoutes.POST("/refresh", authHand.Refresh)
		authRoutes.GET("/me", authHand.Me)
	}

	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8080"
	}

	log.Printf("Starting auth-service server on port %s", appPort)
	r.Run(":" + appPort)
}
