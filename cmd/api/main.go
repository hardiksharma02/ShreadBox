package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hardiksharma/shreadbox/config"
	"github.com/hardiksharma/shreadbox/internal/cleanup"
	"github.com/hardiksharma/shreadbox/internal/handlers"
	"github.com/hardiksharma/shreadbox/internal/storage"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if it exists
	godotenv.Load()

	// Load configuration
	cfg := config.LoadConfig()

	// Set Gin mode
	if os.Getenv("GIN_MODE") != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize storage service
	storageService, err := storage.NewStorage(cfg.StoragePath)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	// Initialize cleanup service
	cleanupService := cleanup.NewService(storageService, cfg.CleanupInterval)
	cleanupService.Start()
	defer cleanupService.Stop()

	// Initialize handlers
	handler := handlers.NewHandler(storageService)

	// Initialize router
	router := gin.Default()

	// Load templates
	router.LoadHTMLGlob("web/templates/*")

	// Serve static files
	router.Static("/static", "web/static")

	// Setup routes
	setupRoutes(router, handler)

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setupRoutes(router *gin.Engine, handler *handlers.Handler) {
	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"time":   time.Now(),
		})
	})

	// API routes
	api := router.Group("/api")
	{
		api.POST("/upload", handler.Upload)
		api.GET("/download/:token", handler.Download)
		api.GET("/status/:token", handler.Status)
	}

	// Web interface routes
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title": "ShreadBox - Secure File Sharing",
		})
	})
}
