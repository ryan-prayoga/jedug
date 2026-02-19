package main

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"

	"jedug-backend/internal/config"
	"jedug-backend/internal/database"
	"jedug-backend/internal/handlers"
	"jedug-backend/internal/middleware"
	"jedug-backend/internal/repository"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	cfg := config.Load()

	pool := database.TryConnect(cfg.DB)
	defer database.Close()

	// Setup upload directory
	uploadDir := filepath.Join(".", "uploads")

	app := fiber.New(fiber.Config{
		AppName:   "Jedug Backend v1.0",
		BodyLimit: 10 * 1024 * 1024, // 10MB max body
	})

	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173, http://localhost:5174, http://localhost:3000",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Serve uploaded files
	app.Static("/uploads", uploadDir, fiber.Static{
		Compress: true,
		MaxAge:   86400 * 30, // 30 days cache
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Jedug Backend is running! 🚀")
	})

	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Apply rate limiter to all API routes
	v1.Use(middleware.APIRateLimiter())

	// Health check
	v1.Get("/health", func(c *fiber.Ctx) error {
		dbStatus := "disconnected"
		if pool != nil {
			if err := pool.Ping(c.Context()); err == nil {
				dbStatus = "connected"
			}
		}
		status := "ok"
		if dbStatus != "connected" {
			status = "degraded"
		}
		return c.JSON(fiber.Map{
			"status":  status,
			"message": "Service is running",
			"db":      dbStatus,
		})
	})

	if pool != nil {
		// Repositories
		districtRepo := repository.NewDistrictRepo(pool)
		reportRepo := repository.NewReportRepo(pool)
		interactionRepo := repository.NewInteractionRepo(pool)

		// Handlers
		districtHandler := handlers.NewDistrictHandler(districtRepo)
		reportHandler := handlers.NewReportHandler(reportRepo, interactionRepo)

		baseURL := "http://localhost:" + cfg.Port
		uploadHandler := handlers.NewUploadHandler(uploadDir, baseURL)

		// ===== District Routes =====
		v1.Get("/districts", districtHandler.GetDistricts)

		// ===== Report Routes =====
		reports := v1.Group("/reports")
		reports.Get("/", reportHandler.ListReports)
		reports.Get("/nearby", reportHandler.GetNearbyReports)
		reports.Get("/:id", reportHandler.GetReport)
		reports.Post("/", middleware.ReportRateLimiter(), reportHandler.CreateReport)

		// Report interactions
		reports.Post("/:id/reactions", reportHandler.AddReaction)
		reports.Get("/:id/comments", reportHandler.GetComments)
		reports.Post("/:id/comments", reportHandler.AddComment)

		// ===== Ranking Routes =====
		v1.Get("/ranking", reportHandler.GetRanking)

		// ===== Upload Routes =====
		v1.Post("/upload", middleware.UploadRateLimiter(), uploadHandler.UploadImage)

	} else {
		log.Println("⚠️  Running without database — most API routes disabled")
	}

	port := cfg.Port

	go func() {
		if err := app.Listen(":" + port); err != nil {
			log.Fatalf("❌ Server error: %v", err)
		}
	}()

	log.Printf("🚀 Jedug Backend running on port %s", port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")
	_ = app.ShutdownWithTimeout(5 * time.Second)
}
