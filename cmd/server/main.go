package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mockey/internal/db"
	"github.com/mockey/internal/server"
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
	_, err := db.Init()
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	r := gin.Default()
	server.SetupRoutes(r)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
