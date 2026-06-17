package main

import (
	"github.com/gin-gonic/gin"

	"epod-service/internal/handler"
	"epod-service/internal/kafka"
	"epod-service/internal/service"
	"epod-service/internal/storage"
)

func main() {

	r := gin.Default()

	r.Static("/uploads", "./uploads")

	storageLayer := storage.NewLocalStorage()

	kafkaProducer := kafka.NewProducer()

	epodService := service.NewEPODService(
		storageLayer,
		kafkaProducer,
	)

	epodHandler := handler.NewEPODHandler(
		epodService,
	)

	r.POST("/upload", epodHandler.Upload)

	r.Run(":8080")
}