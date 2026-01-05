package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"web-photobooth/backend/internal/config"
	"web-photobooth/backend/internal/handlers"
	"web-photobooth/backend/internal/models"
	"web-photobooth/backend/internal/storage"
)

func main() {
	// 1. Load Configuration
	cfg := config.LoadConfig()

	// 2. Initialize Database
	db, err := storage.InitDB(cfg)
	if err != nil {
		log.Printf("Warning: Database connection failed: %v", err)
	} else {
		log.Println("Database connected successfully")
		if err := models.Migrate(db); err != nil {
			log.Printf("Warning: Migration failed: %v", err)
		}
	}

	// 3. Initialize S3
	s3Client, err := storage.InitS3(cfg)
	if err != nil {
		log.Printf("Warning: S3 initialization failed: %v", err)
	} else {
		log.Println("S3 Client initialized")
	}

	// 4. Initialize Handler (Monolithic, no Supabase)
	h := handlers.NewHandler(db, s3Client, cfg.DOBucket, cfg.JWTSecret, cfg.GuestExpirationDays)

	// 5. Setup Router
	r := gin.Default()

	// 6. Configure CORS
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 7. API Routes (Now modularly registered)
	h.RegisterRoutes(r)

	// 8. Start Background Cleanup Worker
	go func() {
		// Run once on startup
		h.CleanupExpiredStrips()
		
		ticker := time.NewTicker(1 * time.Hour)
		for range ticker.C {
			h.CleanupExpiredStrips()
		}
	}()

	log.Printf("Starting monolithic backend on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
