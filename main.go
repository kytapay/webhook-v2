package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kytapay/webhook-v2/config"
	"github.com/kytapay/webhook-v2/controllers"
	"github.com/kytapay/webhook-v2/routes"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize database
	db, err := config.InitDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize Gin router
	r := gin.Default()

	// Middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Initialize controllers
	webhookController := controllers.NewWebhookController(db)

	// Setup routes
	routes.SetupRoutes(r, webhookController)

	// Get port from environment or use default
	port := os.Getenv("WEBHOOK_PORT")
	if port == "" {
		port = "8081" // Default port for webhook service
	}

	log.Printf("Webhook service starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start webhook service:", err)
	}
}

