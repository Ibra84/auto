package main

import (
	"auto/handlers"
	"auto/models"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Подключение к базе данных
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	dsn := os.Getenv("DSN")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	// Автоматическая миграция для создания таблиц
	err = db.AutoMigrate(&models.MediaData{})
	if err != nil {
		log.Fatal("failed to migrate database")
	}

	r := gin.Default()

	handler := handlers.NewHandler(db)

	r.POST("/upload", handler.UploadData)

	r.Run(":8080")
}
