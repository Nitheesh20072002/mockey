package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mockey/exam-api/internal/db"
	"github.com/mockey/exam-api/internal/models"
	"github.com/mockey/exam-api/internal/server"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, relying on environment variables")
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	// initialize DB
	gdb, err := db.Init()
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	// automigrate basic models
	if err := gdb.AutoMigrate(&models.User{}, &models.Exam{}, &models.UploadJob{}); err != nil {
		log.Fatalf("auto migrate failed: %v", err)
	}

	r := gin.Default()
	server.SetupRoutes(r)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
